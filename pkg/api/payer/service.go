package payer

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	payer_api.UnimplementedPayerApiServer

	ctx                 context.Context
	log                 *zap.Logger
	clientManager       *ClientManager
	blockchainPublisher blockchain.IBlockchainPublisher
	payerPrivateKey     *ecdsa.PrivateKey
	nodeSelector        NodeSelectorAlgorithm
	nodeCursorTracker   *NodeCursorTracker
}

func NewPayerApiService(
	ctx context.Context,
	log *zap.Logger,
	registry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
	metadataApiClient MetadataApiClientConstructor,
) (*Service, error) {
	var metadataClient MetadataApiClientConstructor
	clientManager := NewClientManager(log, registry)
	if metadataApiClient == nil {
		metadataClient = &DefaultMetadataApiClientConstructor{clientManager: clientManager}
	} else {
		metadataClient = metadataApiClient
	}
	return &Service{
		ctx:                 ctx,
		log:                 log,
		clientManager:       clientManager,
		payerPrivateKey:     payerPrivateKey,
		blockchainPublisher: blockchainPublisher,
		nodeCursorTracker:   NewNodeCursorTracker(ctx, log, metadataClient),
		nodeSelector:        &StableHashingNodeSelectorAlgorithm{reg: registry},
	}, nil
}

func (s *Service) PublishClientEnvelopes(
	ctx context.Context,
	req *payer_api.PublishClientEnvelopesRequest,
) (*payer_api.PublishClientEnvelopesResponse, error) {
	grouped, err := s.groupEnvelopes(req.Envelopes)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error grouping envelopes: %v", err)
	}

	out := make([]*envelopesProto.OriginatorEnvelope, len(req.Envelopes))

	// For each originator found in the request, publish all matching envelopes to the node
	for originatorId, payloadsWithIndex := range grouped.forNodes {
		s.log.Info("publishing to originator", zap.Uint32("originator_id", originatorId))
		originatorEnvelopes, err := s.publishToNodes(ctx, originatorId, payloadsWithIndex)
		if err != nil {
			s.log.Error("error publishing payer envelopes", zap.Error(err))
			return nil, status.Error(codes.Internal, "error publishing payer envelopes")
		}

		// The originator envelopes come back from the API in the same order as the request
		for idx, originatorEnvelope := range originatorEnvelopes {
			out[payloadsWithIndex[idx].originalIndex] = originatorEnvelope
		}
	}

	for _, payload := range grouped.forBlockchain {
		s.log.Info(
			"publishing to blockchain",
			zap.String("topic", payload.payload.TargetTopic().String()),
		)
		var originatorEnvelope *envelopesProto.OriginatorEnvelope
		if originatorEnvelope, err = s.publishToBlockchain(ctx, payload.payload); err != nil {
			s.log.Error("error publishing payer envelopes", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "error publishing group message: %v", err)
		}
		out[payload.originalIndex] = originatorEnvelope
	}

	return &payer_api.PublishClientEnvelopesResponse{
		OriginatorEnvelopes: out,
	}, nil
}

// A struct that groups client envelopes by their intended destination
type groupedEnvelopes struct {
	// Mapping of originator ID to a list of client envelopes targeting that originator
	forNodes map[uint32][]clientEnvelopeWithIndex
	// Messages meant to be sent to the blockchain
	forBlockchain []clientEnvelopeWithIndex
}

func (s *Service) groupEnvelopes(
	rawEnvelopes []*envelopesProto.ClientEnvelope,
) (*groupedEnvelopes, error) {
	out := groupedEnvelopes{forNodes: make(map[uint32][]clientEnvelopeWithIndex)}

	for i, rawClientEnvelope := range rawEnvelopes {
		clientEnvelope, err := envelopes.NewClientEnvelope(rawClientEnvelope)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Invalid client envelope at index %d: %v",
				i,
				err,
			)
		}

		if !clientEnvelope.TopicMatchesPayload() {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Client envelope at index %d does not match topic",
				i,
			)
		}
		targetTopic := clientEnvelope.TargetTopic()
		aad := clientEnvelope.Aad()

		if shouldSendToBlockchain(targetTopic, aad) {
			out.forBlockchain = append(
				out.forBlockchain,
				newClientEnvelopeWithIndex(i, clientEnvelope),
			)
		} else {
			// backwards compatibility
			var targetNodeId uint32
			// nolint:staticcheck
			if aad.GetTargetOriginator() != 0 {
				targetNodeId = aad.GetTargetOriginator()
			} else {
				node, err := s.nodeSelector.GetNode(targetTopic)
				if err != nil {
					return nil, err
				}
				targetNodeId = node
			}
			out.forNodes[targetNodeId] = append(out.forNodes[targetNodeId], newClientEnvelopeWithIndex(i, clientEnvelope))
		}
	}

	return &out, nil
}

func (s *Service) publishToNodes(
	ctx context.Context,
	originatorID uint32,
	indexedEnvelopes []clientEnvelopeWithIndex,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	conn, err := s.clientManager.GetClient(originatorID)
	if err != nil {
		s.log.Error("error getting client", zap.Error(err))
		return nil, status.Error(codes.Internal, "error getting client")
	}
	client := message_api.NewReplicationApiClient(conn)

	payerEnvelopes, err := s.signAllClientEnvelopes(originatorID, indexedEnvelopes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error signing payer envelopes: %v", err)
	}

	resp, err := client.PublishPayerEnvelopes(ctx, &message_api.PublishPayerEnvelopesRequest{
		PayerEnvelopes: payerEnvelopes,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error publishing payer envelopes: %v", err)
	}

	return resp.OriginatorEnvelopes, nil
}

func (s *Service) publishToBlockchain(
	ctx context.Context,
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopesProto.OriginatorEnvelope, error) {
	targetTopic := clientEnvelope.TargetTopic()
	identifier := targetTopic.Identifier()
	desiredOriginatorId := uint32(1) //TODO: determine this from the chain
	desiredSequenceId := uint64(0)
	kind := targetTopic.Kind()

	// Get the group ID as [32]byte
	idBytes, err := utils.BytesToId(identifier)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error converting identifier to group ID: %v",
			err,
		)
	}

	// Serialize the clientEnvelope for publishing
	payload, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error getting client envelope bytes: %v",
			err,
		)
	}

	var unsignedOriginatorEnvelope *envelopesProto.UnsignedOriginatorEnvelope
	var hash common.Hash
	switch kind {
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		var logMessage *groupmessages.GroupMessagesMessageSent
		if logMessage, err = s.blockchainPublisher.PublishGroupMessage(ctx, idBytes, payload); err != nil {
			return nil, status.Errorf(codes.Internal, "error publishing group message: %v", err)
		}
		if logMessage == nil {
			return nil, status.Errorf(codes.Internal, "received nil logMessage")
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorId,
			logMessage.SequenceId,
			logMessage.Message,
		)
		desiredSequenceId = logMessage.SequenceId

	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		var logMessage *identityupdates.IdentityUpdatesIdentityUpdateCreated
		if logMessage, err = s.blockchainPublisher.PublishIdentityUpdate(ctx, idBytes, payload); err != nil {
			return nil, status.Errorf(codes.Internal, "error publishing identity update: %v", err)
		}
		if logMessage == nil {
			return nil, status.Errorf(codes.Internal, "received nil logMessage")
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorId,
			logMessage.SequenceId,
			logMessage.Update,
		)
		desiredSequenceId = logMessage.SequenceId

	default:
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Unknown blockchain message for topic %s",
			targetTopic.String(),
		)
	}

	unsignedBytes, err := proto.Marshal(unsignedOriginatorEnvelope)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error marshalling unsigned originator envelope: %v",
			err,
		)
	}

	// backwards compatibility
	var targetNodeId uint32
	// nolint:staticcheck
	if clientEnvelope.Aad().GetTargetOriginator() >= 100 {
		targetNodeId = clientEnvelope.Aad().GetTargetOriginator()
	} else {
		node, err := s.nodeSelector.GetNode(targetTopic)
		if err != nil {
			return nil, err
		}
		targetNodeId = node
	}

	err = s.nodeCursorTracker.BlockUntilDesiredCursorReached(
		ctx,
		targetNodeId,
		desiredOriginatorId,
		desiredSequenceId,
	)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof: &envelopesProto.OriginatorEnvelope_BlockchainProof{
			BlockchainProof: &envelopesProto.BlockchainProof{
				TransactionHash: hash.Bytes(),
			},
		},
	}, nil
}

func buildUnsignedOriginatorEnvelopeFromChain(
	targetOriginator uint32,
	sequenceID uint64,
	clientEnvelope []byte,
) *envelopesProto.UnsignedOriginatorEnvelope {
	return &envelopesProto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     targetOriginator,
		OriginatorSequenceId: sequenceID,
		OriginatorNs:         time.Now().UnixNano(), // TODO: get this data from the chain
		PayerEnvelope: &envelopesProto.PayerEnvelope{
			UnsignedClientEnvelope: clientEnvelope,
		},
	}
}

func (s *Service) signAllClientEnvelopes(originatorID uint32,
	indexedEnvelopes []clientEnvelopeWithIndex,
) ([]*envelopesProto.PayerEnvelope, error) {
	out := make([]*envelopesProto.PayerEnvelope, len(indexedEnvelopes))
	for i, indexedEnvelope := range indexedEnvelopes {
		envelope, err := s.signClientEnvelope(originatorID, indexedEnvelope.payload)
		if err != nil {
			return nil, err
		}
		out[i] = envelope
	}
	return out, nil
}

func (s *Service) signClientEnvelope(originatorID uint32,
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopesProto.PayerEnvelope, error) {
	envelopeBytes, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, err
	}

	payerSignature, err := utils.SignClientEnvelope(originatorID, envelopeBytes, s.payerPrivateKey)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: envelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator: originatorID,
	}, nil
}

func shouldSendToBlockchain(targetTopic topic.Topic, aad *envelopesProto.AuthenticatedData) bool {
	switch targetTopic.Kind() {
	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		return true
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		// nolint:staticcheck
		return aad.GetIsCommit() || aad.GetTargetOriginator() != 0 &&
			aad.GetTargetOriginator() < uint32(constants.MAX_BLOCKCHAIN_ORIGINATOR_ID)
	default:
		return false
	}
}

// Wrap a ClientEnvelope in a type that retains the original index from the request inputs
type clientEnvelopeWithIndex struct {
	originalIndex int
	payload       *envelopes.ClientEnvelope
}

func newClientEnvelopeWithIndex(
	index int,
	payload *envelopes.ClientEnvelope,
) clientEnvelopeWithIndex {
	return clientEnvelopeWithIndex{
		originalIndex: index,
		payload:       payload,
	}
}

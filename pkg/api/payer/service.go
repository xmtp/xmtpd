package payer

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
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
)

type Service struct {
	payer_api.UnimplementedPayerApiServer

	ctx                 context.Context
	log                 *zap.Logger
	clientManager       *ClientManager
	blockchainPublisher blockchain.IBlockchainPublisher
	payerPrivateKey     *ecdsa.PrivateKey
}

func NewPayerApiService(
	ctx context.Context,
	log *zap.Logger,
	registry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
) (*Service, error) {
	return &Service{
		ctx:                 ctx,
		log:                 log,
		clientManager:       NewClientManager(log, registry),
		payerPrivateKey:     payerPrivateKey,
		blockchainPublisher: blockchainPublisher,
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
		var originatorEnvelope *envelopesProto.OriginatorEnvelope
		if originatorEnvelope, err = s.publishToBlockchain(ctx, payload.payload); err != nil {
			return nil, status.Errorf(codes.Internal, "error publishing group message: %v", err)
		}
		out[payload.originalIndex] = originatorEnvelope
	}

	return nil, status.Errorf(codes.Unimplemented, "method PublishClientEnvelopes not implemented")
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
	out := groupedEnvelopes{}

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
			out.forNodes[aad.TargetOriginator] = append(out.forNodes[aad.TargetOriginator], newClientEnvelopeWithIndex(i, clientEnvelope))
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

	payerEnvelopes, err := s.signAllClientEnvelopes(indexedEnvelopes)
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

	var hash common.Hash
	switch kind {
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		hash, err = s.blockchainPublisher.PublishGroupMessage(ctx, idBytes, payload)
	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		hash, err = s.blockchainPublisher.PublishIdentityUpdate(ctx, idBytes, payload)
	default:
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Unknown blockchain message for topic %s",
			targetTopic.String(),
		)
	}
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error publishing group message: %v",
			err,
		)
	}

	return &envelopesProto.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: payload,
		Proof: &envelopesProto.OriginatorEnvelope_BlockchainProof{
			BlockchainProof: &envelopesProto.BlockchainProof{
				TransactionHash: hash.Bytes(),
			},
		},
	}, nil
}

func (s *Service) signAllClientEnvelopes(
	indexedEnvelopes []clientEnvelopeWithIndex,
) ([]*envelopesProto.PayerEnvelope, error) {
	out := make([]*envelopesProto.PayerEnvelope, len(indexedEnvelopes))
	for i, indexedEnvelope := range indexedEnvelopes {
		envelope, err := s.signClientEnvelope(indexedEnvelope.payload)
		if err != nil {
			return nil, err
		}
		out[i] = envelope
	}
	return out, nil
}

func (s *Service) signClientEnvelope(
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopesProto.PayerEnvelope, error) {
	envelopeBytes, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, err
	}

	payerSignature, err := utils.SignClientEnvelope(envelopeBytes, s.payerPrivateKey)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: envelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
	}, nil
}

func shouldSendToBlockchain(targetTopic topic.Topic, aad *envelopesProto.AuthenticatedData) bool {
	switch targetTopic.Kind() {
	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		return true
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		return aad.TargetOriginator < constants.MAX_BLOCKCHAIN_ORIGINATOR_ID
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

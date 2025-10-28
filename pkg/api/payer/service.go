// Package payer implements the Payer API service.
package payer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/deserializer"
	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/ethereum/go-ethereum/common"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	getNodesMethod               = "GetNodes"
	publishClientEnvelopesMethod = "PublishClientEnvelopes"

	requestMissingMessageError = "missing request message"
)

type Service struct {
	payer_apiconnect.UnimplementedPayerApiHandler

	ctx                 context.Context
	logger              *zap.Logger
	clientManager       *ClientManager
	blockchainPublisher blockchain.IBlockchainPublisher
	payerPrivateKey     *ecdsa.PrivateKey
	nodeSelector        NodeSelectorAlgorithm
	nodeCursorTracker   *NodeCursorTracker
	nodeRegistry        registry.NodeRegistry
}

var _ payer_apiconnect.PayerApiHandler = (*Service)(nil)

func NewPayerAPIService(
	ctx context.Context,
	logger *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
	metadataAPIClient MetadataAPIClientConstructor,
	clientMetrics *utils.ConnectClientMetrics,
) (*Service, error) {
	if clientMetrics == nil {
		clientMetrics = utils.NewConnectClientMetrics()
	}

	var metadataClient MetadataAPIClientConstructor
	clientManager := NewClientManager(logger, nodeRegistry, clientMetrics)
	if metadataAPIClient == nil {
		metadataClient = &DefaultMetadataAPIClientConstructor{clientManager: clientManager}
	} else {
		metadataClient = metadataAPIClient
	}

	return &Service{
		ctx:                 ctx,
		logger:              logger,
		clientManager:       clientManager,
		payerPrivateKey:     payerPrivateKey,
		blockchainPublisher: blockchainPublisher,
		nodeCursorTracker:   NewNodeCursorTracker(ctx, metadataClient),
		nodeSelector:        &StableHashingNodeSelectorAlgorithm{reg: nodeRegistry},
		nodeRegistry:        nodeRegistry,
	}, nil
}

// GetNodes returns the complete endpoint list of canonical nodes.
func (s *Service) GetNodes(
	ctx context.Context,
	_ *connect.Request[payer_api.GetNodesRequest],
) (*connect.Response[payer_api.GetNodesResponse], error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(getNodesMethod))
	}

	nodes, err := s.nodeRegistry.GetNodes()
	if err != nil {
		return nil, connect.NewError(
			connect.CodeUnavailable,
			fmt.Errorf("failed to fetch nodes: %w", err),
		)
	}

	metrics.EmitPayerGetNodesAvailableNodes(len(nodes))

	if len(nodes) == 0 {
		return nil, connect.NewError(
			connect.CodeUnavailable,
			errors.New("no nodes available"),
		)
	}

	response := connect.NewResponse(&payer_api.GetNodesResponse{
		Nodes: make(map[uint32]string, len(nodes)),
	})

	for _, node := range nodes {
		response.Msg.Nodes[node.NodeID] = node.HTTPAddress
	}

	return response, nil
}

func (s *Service) PublishClientEnvelopes(
	ctx context.Context,
	req *connect.Request[payer_api.PublishClientEnvelopesRequest],
) (*connect.Response[payer_api.PublishClientEnvelopesResponse], error) {
	if req.Msg == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf(requestMissingMessageError),
		)
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(publishClientEnvelopesMethod))
	}

	grouped, err := s.groupEnvelopes(req.Msg.GetEnvelopes())
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("error grouping envelopes: %w", err),
		)
	}

	response := connect.NewResponse(&payer_api.PublishClientEnvelopesResponse{
		OriginatorEnvelopes: make(
			[]*envelopesProto.OriginatorEnvelope,
			len(req.Msg.GetEnvelopes()),
		),
	})

	// For each originator found in the request, publish all matching envelopes to the node
	for originatorID, payloadsWithIndex := range grouped.forNodes {
		s.logger.Debug(
			"publishing to originator",
			utils.OriginatorIDField(originatorID),
			utils.MethodField(publishClientEnvelopesMethod),
			utils.NumEnvelopesField(len(payloadsWithIndex)),
		)

		originatorEnvelopes, err := s.publishToNodeWithRetry(ctx, originatorID, payloadsWithIndex)
		if err != nil {
			s.logger.Error("error publishing payer envelopes", zap.Error(err))
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error publishing payer envelopes: %w", err),
			)
		}

		// The originator envelopes come back from the API in the same order as the request
		for idx, originatorEnvelope := range originatorEnvelopes {
			response.Msg.OriginatorEnvelopes[payloadsWithIndex[idx].originalIndex] = originatorEnvelope
		}
	}

	for _, payload := range grouped.forBlockchain {
		s.logger.Debug(
			"publishing to blockchain",
			utils.MethodField(publishClientEnvelopesMethod),
			utils.TopicField(payload.payload.TargetTopic().String()),
		)

		var originatorEnvelope *envelopesProto.OriginatorEnvelope
		if originatorEnvelope, err = s.publishToBlockchain(ctx, payload.payload); err != nil {
			s.logger.Error("error publishing payer envelopes", zap.Error(err))
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error publishing group message: %w", err),
			)
		}

		response.Msg.OriginatorEnvelopes[payload.originalIndex] = originatorEnvelope
	}

	return response, nil
}

// A struct that groups client envelopes by their intended destination.
type groupedEnvelopes struct {
	// Mapping of originator ID to a list of client envelopes targeting that originator.
	forNodes map[uint32][]clientEnvelopeWithIndex
	// Messages meant to be sent to the blockchain.
	forBlockchain []clientEnvelopeWithIndex
}

func (s *Service) groupEnvelopes(
	rawEnvelopes []*envelopesProto.ClientEnvelope,
) (*groupedEnvelopes, error) {
	out := groupedEnvelopes{forNodes: make(map[uint32][]clientEnvelopeWithIndex)}

	for i, rawClientEnvelope := range rawEnvelopes {
		clientEnvelope, err := envelopes.NewClientEnvelope(rawClientEnvelope)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("invalid client envelope at index %d: %w", i, err),
			)
		}

		if !clientEnvelope.TopicMatchesPayload() {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("client envelope at index %d does not match topic", i),
			)
		}

		toBlockchain, err := shouldSendToBlockchain(clientEnvelope)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("client envelope at index %d can not be parsed: %w", i, err),
			)
		}

		if toBlockchain {
			out.forBlockchain = append(
				out.forBlockchain,
				newClientEnvelopeWithIndex(i, clientEnvelope),
			)
		} else {
			targetNodeID, err := s.nodeSelector.GetNode(clientEnvelope.TargetTopic())
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("error getting node for topic: %w", err),
				)
			}

			out.forNodes[targetNodeID] = append(out.forNodes[targetNodeID], newClientEnvelopeWithIndex(i, clientEnvelope))
		}
	}

	return &out, nil
}

func (s *Service) publishToNodeWithRetry(
	ctx context.Context,
	originatorID uint32,
	indexedEnvelopes []clientEnvelopeWithIndex,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	var banlist []uint32
	var result []*envelopesProto.OriginatorEnvelope
	var err error
	nodeID := originatorID

	topic := indexedEnvelopes[0].payload.TargetTopic()

	for retries := 0; retries < 5; retries++ {
		result, err = s.publishToNodes(ctx, nodeID, indexedEnvelopes)
		if err == nil {
			if retries != 0 {
				metrics.EmitPayerBanlistRetries(originatorID, retries)
			}
			return result, nil
		}

		s.logger.Error(
			"error publishing to node, will retry with the next one",
			utils.OriginatorIDField(nodeID),
			zap.Error(err),
		)

		// Add failed node to banlist and retry
		s.logger.Warn("adding failed node to banlist", utils.OriginatorIDField(nodeID))
		banlist = append(banlist, nodeID)

		nodeID, err = s.nodeSelector.GetNode(topic, banlist)
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}

func (s *Service) publishToNodes(
	ctx context.Context,
	originatorID uint32,
	indexedEnvelopes []clientEnvelopeWithIndex,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	client, err := s.clientManager.GetReplicationClient(originatorID)
	if err != nil {
		s.logger.Error("error getting client", zap.Error(err))
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error getting client: %w", err),
		)
	}

	payerEnvelopes, err := s.signAllClientEnvelopes(originatorID, indexedEnvelopes)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error signing payer envelopes: %w", err),
		)
	}

	start := time.Now()

	response, err := client.PublishPayerEnvelopes(
		ctx,
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: payerEnvelopes,
		}),
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error publishing payer envelopes: %w", err),
		)
	}

	metrics.EmitPayerNodePublishDuration(originatorID, time.Since(start).Seconds())
	metrics.EmitPayerMessageOriginated(originatorID, len(payerEnvelopes))

	return response.Msg.OriginatorEnvelopes, nil
}

func (s *Service) publishToBlockchain(
	ctx context.Context,
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopesProto.OriginatorEnvelope, error) {
	var (
		targetTopic         = clientEnvelope.TargetTopic()
		identifier          = targetTopic.Identifier()
		desiredOriginatorID uint32
		desiredSequenceID   uint64
		kind                = targetTopic.Kind()
	)

	// Serialize the clientEnvelope for publishing
	payload, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error getting client envelope bytes: %w", err),
		)
	}

	start := time.Now()

	var unsignedOriginatorEnvelope *envelopesProto.UnsignedOriginatorEnvelope
	var hash common.Hash
	switch kind {
	case topic.TopicKindGroupMessagesV1:
		desiredOriginatorID = constants.GroupMessageOriginatorID

		var logMessage *gm.GroupMessageBroadcasterMessageSent

		// Get the group ID as [16]byte
		groupID, err := utils.ParseGroupID(identifier)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error converting identifier to group ID: %w", err),
			)
		}

		if logMessage, err = metrics.MeasurePublishToBlockchainMethod("group_message", func() (*gm.GroupMessageBroadcasterMessageSent, error) {
			return s.blockchainPublisher.PublishGroupMessage(ctx, groupID, payload)
		}); err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error publishing group message: %w", err),
			)
		}
		if logMessage == nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("received nil logMessage"),
			)
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope, err = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorID,
			logMessage.SequenceId,
			logMessage.Message,
		)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error building unsigned originator envelope: %w", err),
			)
		}
		desiredSequenceID = logMessage.SequenceId

	case topic.TopicKindIdentityUpdatesV1:
		desiredOriginatorID = constants.IdentityUpdateOriginatorID

		var logMessage *iu.IdentityUpdateBroadcasterIdentityUpdateCreated

		// Get the inbox ID as [32]byte
		inboxID, err := utils.ParseInboxID(identifier)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error converting identifier to inbox ID: %w", err),
			)
		}

		if logMessage, err = metrics.MeasurePublishToBlockchainMethod("identity_update", func() (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
			return s.blockchainPublisher.PublishIdentityUpdate(ctx, inboxID, payload)
		}); err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error publishing identity update: %w", err),
			)
		}
		if logMessage == nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("received nil logMessage"),
			)
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope, err = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorID,
			logMessage.SequenceId,
			logMessage.Update,
		)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("error building unsigned originator envelope: %w", err),
			)
		}
		desiredSequenceID = logMessage.SequenceId

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("unknown blockchain message for topic %s", targetTopic.String()),
		)
	}

	metrics.EmitPayerNodePublishDuration(desiredOriginatorID, time.Since(start).Seconds())
	metrics.EmitPayerMessageOriginated(desiredOriginatorID, 1)

	s.logger.Debug(
		"published message to blockchain",
		utils.DurationMsField(time.Since(start)),
	)

	unsignedBytes, err := proto.Marshal(unsignedOriginatorEnvelope)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error marshalling unsigned originator envelope: %w", err),
		)
	}

	targetNodeID, err := s.nodeSelector.GetNode(targetTopic)
	if err != nil {
		return nil, err
	}

	s.logger.Debug(
		"waiting for message to be processed by node",
		utils.OriginatorIDField(targetNodeID),
	)

	err = s.nodeCursorTracker.BlockUntilDesiredCursorReached(
		ctx,
		targetNodeID,
		desiredOriginatorID,
		desiredSequenceID,
	)
	if err != nil {
		s.logger.Error(
			"chosen node for cursor check is unreachable",
			utils.OriginatorIDField(targetNodeID),
			zap.Error(err),
		)
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
) (*envelopesProto.UnsignedOriginatorEnvelope, error) {
	payerEnvelope := &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvelope,
	}
	payerEnvelopeBytes, err := proto.Marshal(payerEnvelope)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     targetOriginator,
		OriginatorSequenceId: sequenceID,
		OriginatorNs:         time.Now().UnixNano(), // TODO: get this data from the chain
		PayerEnvelopeBytes:   payerEnvelopeBytes,
	}, nil
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

	retentionDuration, err := determineRetentionPolicy(
		clientEnvelope,
	)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: envelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     originatorID,
		MessageRetentionDays: retentionDuration,
	}, nil
}

func determineRetentionPolicy(clientEnvelope *envelopes.ClientEnvelope) (uint32, error) {
	// TODO: mkysel determine expiration for welcomes and key packages

	switch clientEnvelope.TargetTopic().Kind() {
	case topic.TopicKindIdentityUpdatesV1:
		panic("should not be called for identity updates")
	case topic.TopicKindGroupMessagesV1:
		switch payload := clientEnvelope.Payload().(type) {
		case *envelopesProto.ClientEnvelope_GroupMessage:
			isCommit, err := deserializer.IsGroupMessageCommit(payload)
			if err != nil {
				return 0, err
			}
			if isCommit {
				panic("should not be called for group message commits")
			}
		default:
			panic("mismatched payload type")
		}
	}

	return constants.DefaultStorageDurationDays, nil
}

func shouldSendToBlockchain(clientEnvelope *envelopes.ClientEnvelope) (bool, error) {
	switch clientEnvelope.TargetTopic().Kind() {
	case topic.TopicKindIdentityUpdatesV1:
		return true, nil
	case topic.TopicKindGroupMessagesV1:
		switch payload := clientEnvelope.Payload().(type) {
		case *envelopesProto.ClientEnvelope_GroupMessage:
			isCommit, err := deserializer.IsGroupMessageCommit(payload)
			if err != nil {
				return false, err
			}
			return isCommit, nil
		default:
			panic("mismatched payload type")
		}
	default:
		return false, nil
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

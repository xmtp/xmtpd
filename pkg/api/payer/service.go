// Package payer implements the Payer API service.
package payer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/api/payer/selectors"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/deserializer"
	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/ethereum/go-ethereum/common"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const requestMissingMessageError = "missing request message"

type Service struct {
	payer_apiconnect.UnimplementedPayerApiHandler

	ctx                 context.Context
	logger              *zap.Logger
	clientManager       *ClientManager
	blockchainPublisher blockchain.IBlockchainPublisher
	payerPrivateKey     *ecdsa.PrivateKey
	nodeSelector        selectors.NodeSelectorAlgorithm
	nodeRegistry        registry.NodeRegistry
	maxPayerMessageSize uint64
}

var _ payer_apiconnect.PayerApiHandler = (*Service)(nil)

func NewPayerAPIService(
	ctx context.Context,
	logger *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
	clientMetrics *grpcprom.ClientMetrics,
	maxPayerMessageSize uint64,
) (*Service, error) {
	return NewPayerAPIServiceWithSelector(
		ctx,
		logger,
		nodeRegistry,
		payerPrivateKey,
		blockchainPublisher,
		clientMetrics,
		maxPayerMessageSize,
		nil,
	)
}

func NewPayerAPIServiceWithSelector(
	ctx context.Context,
	logger *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
	clientMetrics *grpcprom.ClientMetrics,
	maxPayerMessageSize uint64,
	nodeSelector selectors.NodeSelectorAlgorithm,
) (*Service, error) {
	if clientMetrics == nil {
		clientMetrics = grpcprom.NewClientMetrics()
	}

	if nodeSelector == nil {
		nodeSelector = selectors.NewStableHashingNodeSelectorAlgorithm(nodeRegistry)
	}

	clientManager := NewClientManager(logger, nodeRegistry, clientMetrics)

	return &Service{
		ctx:                 ctx,
		logger:              logger,
		clientManager:       clientManager,
		payerPrivateKey:     payerPrivateKey,
		blockchainPublisher: blockchainPublisher,
		nodeSelector:        nodeSelector,
		nodeRegistry:        nodeRegistry,
		maxPayerMessageSize: maxPayerMessageSize,
	}, nil
}

// GetNodes returns the complete endpoint list of canonical nodes.
func (s *Service) GetNodes(
	ctx context.Context,
	req *connect.Request[payer_api.GetNodesRequest],
) (*connect.Response[payer_api.GetNodesResponse], error) {
	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(req.Spec().Procedure))
	}

	nodes, err := s.nodeRegistry.GetNodes()
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("failed to fetch nodes: %w", err),
		)
	}

	metrics.EmitGatewayGetNodesAvailableNodes(len(nodes))

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
			errors.New(requestMissingMessageError),
		)
	}

	if s.logger.Core().Enabled(zap.DebugLevel) {
		s.logger.Debug("received request", utils.MethodField(req.Spec().Procedure))
	}

	grouped, err := s.groupEnvelopes(req.Msg.GetEnvelopes())
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("error grouping envelopes: %w", err),
		)
	}

	if s.maxPayerMessageSize != 0 {
		for _, env := range grouped.forBlockchain {
			bytes, err := env.payload.Bytes()
			if err != nil {
				return nil, err
			}
			if len(bytes) > int(s.maxPayerMessageSize) {
				return nil, status.Errorf(
					codes.InvalidArgument,
					"message at index %d too large",
					env.originalIndex,
				)
			}
		}
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
			utils.MethodField(req.Spec().Procedure),
			utils.NumEnvelopesField(len(payloadsWithIndex)),
		)

		originatorEnvelopes, err := s.publishToNodeWithRetry(ctx, originatorID, payloadsWithIndex)
		if err != nil {
			if ctx.Err() != nil {
				return nil, connect.NewError(
					connect.CodeCanceled,
					errors.New("request canceled by client"),
				)
			}

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
			utils.MethodField(req.Spec().Procedure),
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
				metrics.EmitGatewayBanlistRetries(originatorID, retries)
			}
			return result, nil
		}

		// Don't retry or ban nodes if context was cancelled.
		if ctx.Err() != nil {
			s.logger.Debug("request canceled by client", zap.Error(err))
			return nil, err
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
	conn, err := s.clientManager.GetClientConnection(originatorID)
	if err != nil {
		s.logger.Error("error getting client", zap.Error(err))
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error getting client: %w", err),
		)
	}

	client := message_api.NewReplicationApiClient(conn)

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
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: payerEnvelopes,
		},
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error publishing payer envelopes: %w", err),
		)
	}

	metrics.EmitGatewayPublishDuration(originatorID, time.Since(start).Seconds())
	metrics.EmitGatewayMessageOriginated(originatorID, len(payerEnvelopes))

	return response.GetOriginatorEnvelopes(), nil
}

func (s *Service) publishToBlockchain(
	ctx context.Context,
	clientEnvelope *envelopes.ClientEnvelope,
) (receipt *envelopesProto.OriginatorEnvelope, err error) {
	var (
		start               = time.Now()
		targetTopic         = clientEnvelope.TargetTopic()
		identifier          = targetTopic.Identifier()
		desiredOriginatorID uint32
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

	case topic.TopicKindIdentityUpdatesV1:
		desiredOriginatorID = constants.IdentityUpdateOriginatorID

		var logMessage *iu.IdentityUpdateBroadcasterIdentityUpdateCreated

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

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("unknown blockchain message for topic %s", targetTopic.String()),
		)
	}

	unsignedBytes, err := proto.Marshal(unsignedOriginatorEnvelope)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error marshalling unsigned originator envelope: %w", err),
		)
	}

	defer func() {
		if err != nil {
			s.logger.Error(
				"error publishing message to blockchain",
				utils.DurationMsField(time.Since(start)),
				utils.TopicField(targetTopic.String()),
				zap.Error(err),
			)
		} else {
			s.logger.Debug(
				"published message to blockchain",
				utils.DurationMsField(time.Since(start)),
				utils.TopicField(targetTopic.String()),
				utils.HashField(hash.String()),
			)
		}

		metrics.EmitGatewayPublishDuration(desiredOriginatorID, time.Since(start).Seconds())
		metrics.EmitGatewayMessageOriginated(desiredOriginatorID, 1)
	}()

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

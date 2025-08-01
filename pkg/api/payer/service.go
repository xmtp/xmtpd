package payer

import (
	"context"
	"crypto/ecdsa"
	"math/rand"
	"time"

	"github.com/xmtp/xmtpd/pkg/indexer/app_chain/contracts"

	"github.com/xmtp/xmtpd/pkg/deserializer"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

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
	nodeRegistry        registry.NodeRegistry
}

func NewPayerApiService(
	ctx context.Context,
	log *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	payerPrivateKey *ecdsa.PrivateKey,
	blockchainPublisher blockchain.IBlockchainPublisher,
	metadataApiClient MetadataApiClientConstructor,
	clientMetrics *grpcprom.ClientMetrics,
) (*Service, error) {
	if clientMetrics == nil {
		clientMetrics = grpcprom.NewClientMetrics()
	}

	var metadataClient MetadataApiClientConstructor
	clientManager := NewClientManager(log, nodeRegistry, clientMetrics)
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
		nodeSelector:        &StableHashingNodeSelectorAlgorithm{reg: nodeRegistry},
		nodeRegistry:        nodeRegistry,
	}, nil
}

// GetReaderNode returns a reader node and a list of backup nodes.
// For now, the reader node is chosen randomly from the list of nodes.
// In the future, different algorithms can be implemented and selected in the request.
func (s *Service) GetReaderNode(
	ctx context.Context,
	req *payer_api.GetReaderNodeRequest,
) (resp *payer_api.GetReaderNodeResponse, err error) {
	var nodes []registry.Node

	defer func() {
		metrics.EmitPayerGetReaderNodeAvailableNodes(len(nodes))
	}()

	nodes, err = s.nodeRegistry.GetNodes()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to fetch nodes: %v", err)
	}

	if len(nodes) == 0 {
		return nil, status.Errorf(codes.Unavailable, "no nodes available")
	}

	primaryUrl, backupUrls := getReaderNodeRandom(nodes)

	return &payer_api.GetReaderNodeResponse{
		ReaderNodeUrl:  primaryUrl,
		BackupNodeUrls: backupUrls,
	}, nil
}

func getReaderNodeRandom(nodes []registry.Node) (string, []string) {
	shuffled := make([]registry.Node, len(nodes))
	copy(shuffled, nodes)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	primaryUrl := shuffled[0].HttpAddress

	backupUrls := make([]string, 0, len(shuffled)-1)
	for _, node := range shuffled[1:] {
		backupUrls = append(backupUrls, node.HttpAddress)
	}

	return primaryUrl, backupUrls
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
		originatorEnvelopes, err := s.publishToNodeWithRetry(ctx, originatorId, payloadsWithIndex)
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

		toBlockchain, err := shouldSendToBlockchain(clientEnvelope)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Client envelope at index %d can not be parsed: %v",
				i, err,
			)
		}

		if toBlockchain {
			out.forBlockchain = append(
				out.forBlockchain,
				newClientEnvelopeWithIndex(i, clientEnvelope),
			)
		} else {
			targetNodeId, err := s.nodeSelector.GetNode(clientEnvelope.TargetTopic())
			if err != nil {
				return nil, err
			}

			out.forNodes[targetNodeId] = append(out.forNodes[targetNodeId], newClientEnvelopeWithIndex(i, clientEnvelope))
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

		s.log.Error(
			"error publishing to node. Retrying with the next one",
			zap.Uint32("failed_node", nodeID),
			zap.Error(err),
		)

		// Add failed node to banlist and retry
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

	start := time.Now()
	resp, err := client.PublishPayerEnvelopes(ctx, &message_api.PublishPayerEnvelopesRequest{
		PayerEnvelopes: payerEnvelopes,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error publishing payer envelopes: %v", err)
	}

	metrics.EmitPayerNodePublishDuration(originatorID, time.Since(start).Seconds())
	metrics.EmitPayerMessageOriginated(originatorID, len(payerEnvelopes))
	return resp.OriginatorEnvelopes, nil
}

func (s *Service) publishToBlockchain(
	ctx context.Context,
	clientEnvelope *envelopes.ClientEnvelope,
) (*envelopesProto.OriginatorEnvelope, error) {
	var (
		targetTopic         = clientEnvelope.TargetTopic()
		identifier          = targetTopic.Identifier()
		desiredOriginatorId uint32
		desiredSequenceId   uint64
		kind                = targetTopic.Kind()
	)

	// Serialize the clientEnvelope for publishing
	payload, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error getting client envelope bytes: %v",
			err,
		)
	}

	start := time.Now()

	var unsignedOriginatorEnvelope *envelopesProto.UnsignedOriginatorEnvelope
	var hash common.Hash
	switch kind {
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		desiredOriginatorId = contracts.GROUP_MESSAGE_ORIGINATOR_ID

		var logMessage *gm.GroupMessageBroadcasterMessageSent

		// Get the group ID as [16]byte
		groupID, err := utils.ParseGroupID(identifier)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"error converting identifier to group ID: %v",
				err,
			)
		}

		if logMessage, err = metrics.MeasurePublishToBlockchainMethod("group_message", func() (*gm.GroupMessageBroadcasterMessageSent, error) {
			return s.blockchainPublisher.PublishGroupMessage(ctx, groupID, payload)
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "error publishing group message: %v", err)
		}
		if logMessage == nil {
			return nil, status.Errorf(codes.Internal, "received nil logMessage")
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope, err = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorId,
			logMessage.SequenceId,
			logMessage.Message,
		)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"error building unsigned originator envelope: %v",
				err,
			)
		}
		desiredSequenceId = logMessage.SequenceId

	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		desiredOriginatorId = contracts.IDENTITY_UPDATE_ORIGINATOR_ID

		var logMessage *iu.IdentityUpdateBroadcasterIdentityUpdateCreated

		// Get the inbox ID as [32]byte
		inboxID, err := utils.ParseInboxID(identifier)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"error converting identifier to inbox ID: %v",
				err,
			)
		}

		if logMessage, err = metrics.MeasurePublishToBlockchainMethod("identity_update", func() (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
			return s.blockchainPublisher.PublishIdentityUpdate(ctx, inboxID, payload)
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "error publishing identity update: %v", err)
		}
		if logMessage == nil {
			return nil, status.Errorf(codes.Internal, "received nil logMessage")
		}

		hash = logMessage.Raw.TxHash
		unsignedOriginatorEnvelope, err = buildUnsignedOriginatorEnvelopeFromChain(
			desiredOriginatorId,
			logMessage.SequenceId,
			logMessage.Update,
		)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"error building unsigned originator envelope: %v",
				err,
			)
		}
		desiredSequenceId = logMessage.SequenceId

	default:
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Unknown blockchain message for topic %s",
			targetTopic.String(),
		)
	}

	metrics.EmitPayerNodePublishDuration(desiredOriginatorId, time.Since(start).Seconds())
	metrics.EmitPayerMessageOriginated(desiredOriginatorId, 1)

	s.log.Debug(
		"published message to blockchain",
		zap.Float64("seconds", time.Since(start).Seconds()),
	)

	unsignedBytes, err := proto.Marshal(unsignedOriginatorEnvelope)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error marshalling unsigned originator envelope: %v",
			err,
		)
	}

	targetNodeId, err := s.nodeSelector.GetNode(targetTopic)
	if err != nil {
		return nil, err
	}

	s.log.Debug(
		"Waiting for message to be processed by node",
		zap.Uint32("target_node_id", targetNodeId),
	)

	err = s.nodeCursorTracker.BlockUntilDesiredCursorReached(
		ctx,
		targetNodeId,
		desiredOriginatorId,
		desiredSequenceId,
	)
	if err != nil {
		s.log.Error(
			"Chosen node for cursor check is unreachable",
			zap.Uint32("targetNodeId", targetNodeId),
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
	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		panic("should not be called for identity updates")
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
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

	return constants.DEFAULT_STORAGE_DURATION_DAYS, nil
}

func shouldSendToBlockchain(clientEnvelope *envelopes.ClientEnvelope) (bool, error) {
	switch clientEnvelope.TargetTopic().Kind() {
	case topic.TOPIC_KIND_IDENTITY_UPDATES_V1:
		return true, nil
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
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

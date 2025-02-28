package indexer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func startIndexing(t *testing.T) (*queries.Queries, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := testutils.NewLog(t)
	db, _, cleanup := testutils.NewDB(t, ctx)
	appChainOptions, _ := testutils.GetContractsOptions(t)
	validationService := mlsvalidate.NewMockMLSValidationService(t)

	indx := NewIndexer(ctx, logger)
	err := indx.StartIndexer(db, appChainOptions, validationService)
	require.NoError(t, err)

	return queries.New(db), ctx, func() {
		cleanup()
		cancel()
	}
}

func messagePublisher(t *testing.T, ctx context.Context) *blockchain.BlockchainPublisher {
	payerCfg := testutils.GetPayerOptions(t)
	appChainOptions, _ := testutils.GetContractsOptions(t)
	var signer blockchain.TransactionSigner
	signer, err := blockchain.NewPrivateKeySigner(payerCfg.PrivateKey, appChainOptions.ChainID)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, appChainOptions.RpcUrl)
	require.NoError(t, err)

	publisher, err := blockchain.NewBlockchainPublisher(
		ctx,
		testutils.NewLog(t),
		client,
		signer,
		appChainOptions,
	)
	require.NoError(t, err)

	return publisher
}

func TestStoreMessages(t *testing.T) {
	querier, ctx, cleanup := startIndexing(t)
	publisher := messagePublisher(t, ctx)
	defer cleanup()

	message := testutils.RandomBytes(78)
	groupID := testutils.RandomGroupID()
	msgTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupID[:]).Bytes()

	clientEnvelope := envelopesTestUtils.CreateGroupMessageClientEnvelope(groupID, message)
	clientEnvelopeBytes, err := proto.Marshal(clientEnvelope)
	require.NoError(t, err)

	// Publish the message onto the blockchain
	_, err = publisher.PublishGroupMessage(ctx, groupID, clientEnvelopeBytes)
	require.NoError(t, err)

	// Poll the DB until the stored message shows up
	require.Eventually(t, func() bool {
		results, err := querier.SelectGatewayEnvelopes(
			context.Background(),
			queries.SelectGatewayEnvelopesParams{
				Topics: [][]byte{msgTopic},
			},
		)
		require.NoError(t, err)

		if len(results) == 0 {
			return false
		}

		firstEnvelope := results[0]
		_, err = envelopes.NewOriginatorEnvelopeFromBytes(
			firstEnvelope.OriginatorEnvelope,
		)
		require.NoError(t, err)
		require.Equal(t, firstEnvelope.Topic, msgTopic)

		return true
	}, 5*time.Second, 100*time.Millisecond, "Failed to find indexed envelope")
}

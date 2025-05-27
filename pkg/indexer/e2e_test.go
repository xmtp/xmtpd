package indexer_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/indexer"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/mocks/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	envelopesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func startIndexing(
	t *testing.T,
) (*sql.DB, *queries.Queries, config.ContractsOptions, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := testutils.NewLog(t)
	db, _ := testutils.NewDB(t, ctx)

	rpcUrl := anvil.StartAnvil(t, false)
	cfg := testutils.NewContractsOptions(t, rpcUrl)

	validationService := mlsvalidate.NewMockMLSValidationService(t)

	indx, err := indexer.NewIndexer(ctx, logger, db, cfg, validationService)
	require.NoError(t, err)
	indx.StartIndexer()

	return db, queries.New(db), cfg, ctx
}

func messagePublisher(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	contractsCfg config.ContractsOptions,
) *blockchain.BlockchainPublisher {
	payerCfg := testutils.GetPayerOptions(t)
	var signer blockchain.TransactionSigner
	signer, err := blockchain.NewPrivateKeySigner(
		payerCfg.PrivateKey,
		contractsCfg.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsCfg.AppChain.RpcURL)
	require.NoError(t, err)

	nonceManager := blockchain.NewSQLBackedNonceManager(db, testutils.NewLog(t))

	publisher, err := blockchain.NewBlockchainPublisher(
		ctx,
		testutils.NewLog(t),
		client,
		signer,
		contractsCfg,
		nonceManager,
	)
	require.NoError(t, err)

	return publisher
}

func TestStoreMessages(t *testing.T) {
	db, querier, cfg, ctx := startIndexing(t)
	publisher := messagePublisher(t, ctx, db, cfg)

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

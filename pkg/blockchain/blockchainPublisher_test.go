package blockchain_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"strings"
	"sync"
	"testing"
)

func buildPublisher(t *testing.T) (*blockchain.BlockchainPublisher, string) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := testutils.NewLog(t)
	rpcUrl, cleanup := anvil.StartAnvil(t, true)
	t.Cleanup(cleanup)
	go func() {
		err := testutils.SubscribePendingTxs(t, rpcUrl)
		require.NoError(t, err)
	}()
	contractsOptions := testutils.NewContractsOptions(rpcUrl)
	// Set the nodes contract address to the newly deployed contract
	contractsOptions.NodesContractAddress = testutils.DeployNodesContract(t, rpcUrl)
	contractsOptions.MessagesContractAddress = testutils.DeployGroupMessagesContract(t, rpcUrl)
	contractsOptions.IdentityUpdatesContractAddress = testutils.DeployIdentityUpdatesContract(
		t,
		rpcUrl,
	)

	go func() {
		err := testutils.SubscribeToLogs(t, rpcUrl, contractsOptions.NodesContractAddress)
		require.NoError(t, err)
	}()
	go func() {
		err := testutils.SubscribeToLogs(t, rpcUrl, contractsOptions.MessagesContractAddress)
		require.NoError(t, err)
	}()
	go func() {
		err := testutils.SubscribeToLogs(t, rpcUrl, contractsOptions.IdentityUpdatesContractAddress)
		require.NoError(t, err)
	}()

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(ctx, contractsOptions.RpcUrl)
	require.NoError(t, err)
	t.Cleanup(
		func() {
			client.Close()
		})

	nonceManager := NewTestNonceManager(logger)

	publisher, err := blockchain.NewBlockchainPublisher(
		ctx,
		logger,
		client,
		signer,
		contractsOptions,
		nonceManager,
	)
	require.NoError(t, err)

	return publisher, rpcUrl
}

func TestPublishIdentityUpdate(t *testing.T) {
	publisher, _ := buildPublisher(t)
	tests := []struct {
		name           string
		inboxId        [32]byte
		identityUpdate []byte
		ctx            context.Context
		wantErr        bool
	}{
		{
			name:           "cancelled context",
			inboxId:        testutils.RandomGroupID(),
			identityUpdate: testutils.RandomBytes(100),
			ctx:            testutils.CancelledContext(),
			wantErr:        true,
		},
		{
			name:           "empty update",
			inboxId:        testutils.RandomGroupID(),
			identityUpdate: []byte{},
			ctx:            context.Background(),
			wantErr:        true,
		},
		{
			name:           "happy path",
			inboxId:        testutils.RandomGroupID(),
			identityUpdate: testutils.RandomBytes(104),
			ctx:            context.Background(),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessage, err := publisher.PublishIdentityUpdate(
				tt.ctx,
				tt.inboxId,
				tt.identityUpdate,
			)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessage)
			require.Equal(t, tt.inboxId, logMessage.InboxId)
			require.Equal(t, tt.identityUpdate, logMessage.Update)
			require.Greater(t, logMessage.SequenceId, uint64(0))
			require.NotNil(t, logMessage.Raw.TxHash)
		})
	}
}

func TestPublishGroupMessage(t *testing.T) {
	publisher, _ := buildPublisher(t)

	tests := []struct {
		name    string
		groupID [32]byte
		message []byte
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "happy path",
			groupID: testutils.RandomGroupID(),
			message: testutils.RandomBytes(100),
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "empty message",
			groupID: testutils.RandomGroupID(),
			message: []byte{},
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:    "cancelled context",
			groupID: testutils.RandomGroupID(),
			message: testutils.RandomBytes(100),
			ctx:     testutils.CancelledContext(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessage, err := publisher.PublishGroupMessage(tt.ctx, tt.groupID, tt.message)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessage)
			require.Equal(t, tt.groupID, logMessage.GroupId)
			require.Equal(t, tt.message, logMessage.Message)
			require.Greater(t, logMessage.SequenceId, uint64(0))
			require.NotNil(t, logMessage.Raw.TxHash)
		})
	}
}

func TestPublishGroupMessageConcurrent(t *testing.T) {
	publisher, rpcUrl := buildPublisher(t)
	defer testutils.LogTransactionsIfFailed(t, rpcUrl)

	const parallelRuns = 30
	var wg sync.WaitGroup
	errSet := sync.Map{}

	for i := 0; i < parallelRuns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := publisher.PublishGroupMessage(
				context.Background(),
				testutils.RandomGroupID(),
				testutils.RandomBytes(100),
			)
			if err != nil {
				errSet.Store(err.Error(), struct{}{})
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Collect and print unique errors
	var uniqueErrors []string
	errSet.Range(func(key, value interface{}) bool {
		uniqueErrors = append(uniqueErrors, key.(string))
		return true
	})

	if len(uniqueErrors) > 0 {
		t.Errorf("Errors encountered:\n%s", strings.Join(uniqueErrors, "\n"))
	}
}

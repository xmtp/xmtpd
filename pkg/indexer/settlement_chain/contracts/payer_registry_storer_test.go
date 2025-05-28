package contracts

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type payerRegistryStorerTester struct {
	abi    *abi.ABI
	storer *PayerRegistryStorer
}

func TestStorePayerRegistryErrorNoTopics(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	err := tester.storer.StoreLog(t.Context(), types.Log{})

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerRegistryLog,
		errors.New("no topics"),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerRegistryErrorUnknownEvent(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := types.Log{
		Topics: []common.Hash{common.HexToHash("UnknownEvent")},
	}

	err := tester.storer.StoreLog(t.Context(), log)

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerRegistryLog,
		fmt.Errorf("no event with id: %#x", log.Topics[0].Hex()),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerRegistryErrorUnhandledEvent(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := tester.newLog(t, "WithdrawalFinalized")

	err := tester.storer.StoreLog(t.Context(), log)

	expectedErr := re.NewNonRecoverableError(
		ErrPayerRegistryUnhandledEvent,
		errors.New("WithdrawalFinalized"),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerRegistryUsageSettled(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := tester.newLog(t, "UsageSettled")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func TestStorePayerRegistryDeposit(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := tester.newLog(t, "Deposit")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func TestStorePayerRegistryWithdrawalRequested(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := tester.newLog(t, "WithdrawalRequested")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func TestStorePayerRegistryWithdrawalCancelled(t *testing.T) {
	tester := buildPayerRegistryStorerTester(t)

	log := tester.newLog(t, "WithdrawalCancelled")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func buildPayerRegistryStorerTester(t *testing.T) *payerRegistryStorerTester {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Dependencies.
	db, _ := testutils.NewDB(t, ctx)
	queryImpl := queries.New(db)
	rpcUrl := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(rpcUrl)

	// Chain client.
	client, err := blockchain.NewClient(ctx, config.AppChain.RpcURL)
	require.NoError(t, err)

	// Contract.
	contract, err := pr.NewPayerRegistry(
		common.HexToAddress(config.SettlementChain.PayerRegistryAddress),
		client,
	)
	require.NoError(t, err)

	// Storer and ABI.
	storer, err := NewPayerRegistryStorer(queryImpl, testutils.NewLog(t), contract)
	require.NoError(t, err)

	abi, err := pr.PayerRegistryMetaData.GetAbi()
	require.NoError(t, err)

	return &payerRegistryStorerTester{
		abi:    abi,
		storer: storer,
	}
}

// TODO: Placeholder. This will be replaced with newDepositLog, newWithdrawalRequestedLog, etc.
func (st *payerRegistryStorerTester) newLog(t *testing.T, event string) types.Log {
	topic, err := utils.GetEventTopic(st.abi, event)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
	}
}

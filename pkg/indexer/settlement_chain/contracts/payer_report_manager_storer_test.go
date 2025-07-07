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
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type payerReportManagerStorerTester struct {
	abi    *abi.ABI
	storer *PayerReportManagerStorer
}

func TestStorePayerReportManagerErrorNoTopics(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	err := tester.storer.StoreLog(t.Context(), types.Log{})

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerReportManagerLog,
		errors.New("no topics"),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerReportManagerErrorUnknownEvent(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	log := types.Log{
		Topics: []common.Hash{common.HexToHash("UnknownEvent")},
	}

	err := tester.storer.StoreLog(t.Context(), log)

	expectedErr := re.NewNonRecoverableError(
		ErrParsePayerReportManagerLog,
		fmt.Errorf("no event with id: %#x", log.Topics[0].Hex()),
	)

	require.Error(t, err)
	require.ErrorAs(t, err, &expectedErr)
}

func TestStorePayerReportManagerPayerReportSubmitted(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	log := tester.newLog(t, "PayerReportSubmitted")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func TestStorePayerReportManagerPayerReportSubsetSettled(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	log := tester.newLog(t, "PayerReportSubsetSettled")

	err := tester.storer.StoreLog(t.Context(), log)

	require.NoError(t, err)
}

func buildPayerReportManagerStorerTester(t *testing.T) *payerReportManagerStorerTester {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Dependencies.
	db, _ := testutils.NewDB(t, ctx)
	queryImpl := queries.New(db)
	wsUrl := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(t, wsUrl)

	// Chain client.
	client, err := blockchain.NewClient(
		ctx,
		blockchain.WithWebSocketURL(config.AppChain.WssURL),
	)
	require.NoError(t, err)

	// Contract.
	contract, err := p.NewPayerReportManager(
		common.HexToAddress(config.SettlementChain.PayerReportManagerAddress),
		client,
	)
	require.NoError(t, err)

	// Storer and ABI.
	storer, err := NewPayerReportManagerStorer(queryImpl, testutils.NewLog(t), contract)
	require.NoError(t, err)

	abi, err := p.PayerReportManagerMetaData.GetAbi()
	require.NoError(t, err)

	return &payerReportManagerStorerTester{
		abi:    abi,
		storer: storer,
	}
}

func (st *payerReportManagerStorerTester) newLog(t *testing.T, event string) types.Log {
	topic, err := utils.GetEventTopic(st.abi, event)
	require.NoError(t, err)

	return types.Log{
		Topics: []common.Hash{topic},
	}
}

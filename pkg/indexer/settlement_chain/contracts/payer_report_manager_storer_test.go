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
	payerreport "github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type payerReportManagerStorerTester struct {
	abi     *abi.ABI
	storer  *PayerReportManagerStorer
	queries *queries.Queries
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

	originatorNodeID := uint32(1)

	log := tester.newPayerReportSubmittedLog(t, &payerreport.PayerReport{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       100,
		EndMinuteSinceEpoch: 200,
		PayersMerkleRoot:    testutils.RandomInboxIDBytes(),
		ActiveNodeIDs:       []uint32{1, 2, 3},
	}, 0)

	err := tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	res, queryErr := tester.queries.FetchPayerReports(t.Context(), queries.FetchPayerReportsParams{
		OriginatorNodeID: utils.NewNullInt32(&originatorNodeID),
	})
	require.Nil(t, queryErr)
	require.Len(t, res, 1)

	require.Equal(t, int32(200), res[0].EndMinuteSinceEpoch)
	require.Equal(t, int64(0), res[0].StartSequenceID)
	require.Equal(t, int64(100), res[0].EndSequenceID)
	require.Equal(t, []int32{1, 2, 3}, res[0].ActiveNodeIds)
}

func TestStorePayerReportManagerPayerReportSubmittedIdempotency(t *testing.T) {
	tester := buildPayerReportManagerStorerTester(t)

	originatorNodeID := uint32(1)

	log := tester.newPayerReportSubmittedLog(t, &payerreport.PayerReport{
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  0,
		EndSequenceID:    100,
		PayersMerkleRoot: testutils.RandomInboxIDBytes(),
		ActiveNodeIDs:    []uint32{1, 2, 3},
	}, 0)

	err := tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	err = tester.storer.StoreLog(t.Context(), log)
	require.NoError(t, err)

	res, queryErr := tester.queries.FetchPayerReports(t.Context(), queries.FetchPayerReportsParams{
		OriginatorNodeID: utils.NewNullInt32(&originatorNodeID),
	})
	require.Nil(t, queryErr)
	require.Len(t, res, 1)
}

func buildPayerReportManagerStorerTester(t *testing.T) *payerReportManagerStorerTester {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Dependencies.
	db, _ := testutils.NewDB(t, ctx)
	wsURL, rpcURL := anvil.StartAnvil(t, false)
	config := testutils.NewContractsOptions(t, rpcURL, wsURL)

	// Chain client.
	client, err := blockchain.NewRPCClient(
		ctx,
		config.AppChain.RPCURL,
	)
	require.NoError(t, err)

	// Contract.
	contract, err := p.NewPayerReportManager(
		common.HexToAddress(config.SettlementChain.PayerReportManagerAddress),
		client,
	)
	require.NoError(t, err)

	// Storer and ABI.
	storer, err := NewPayerReportManagerStorer(db, testutils.NewLog(t), contract)
	require.NoError(t, err)

	abi, err := p.PayerReportManagerMetaData.GetAbi()
	require.NoError(t, err)

	return &payerReportManagerStorerTester{
		abi:     abi,
		storer:  storer,
		queries: queries.New(db),
	}
}

func (st *payerReportManagerStorerTester) newPayerReportSubmittedLog(
	t *testing.T,
	report *payerreport.PayerReport,
	payerReportIndex uint64,
) types.Log {
	return testutils.BuildPayerReportSubmittedEvent(
		t,
		report.OriginatorNodeID,
		payerReportIndex,
		report.StartSequenceID,
		report.EndSequenceID,
		uint64(report.EndMinuteSinceEpoch),
		report.PayersMerkleRoot,
		report.ActiveNodeIDs,
	)
}

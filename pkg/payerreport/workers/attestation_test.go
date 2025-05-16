package workers

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	payerreportMocks "github.com/xmtp/xmtpd/pkg/mocks/payerreport"
	registrantMocks "github.com/xmtp/xmtpd/pkg/mocks/registrant"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

var domainSeparator = common.BytesToHash(testutils.RandomBytes(32))

func testAttestationWorker(
	t *testing.T,
	pollInterval time.Duration,
) (*AttestationWorker, *payerreport.Store, *registrantMocks.MockIRegistrant, *payerreportMocks.MockIPayerReportVerifier) {
	log := testutils.NewLog(t)
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	store := payerreport.NewStore(db, log)
	mockRegistrant := registrantMocks.NewMockIRegistrant(t)
	mockRegistrant.EXPECT().
		SignPayerReportAttestation(mock.Anything, mock.Anything).
		Return(&payerreport.NodeSignature{
			Signature: []byte("signature"),
			NodeID:    1,
		}, nil).
		Maybe()
	mockRegistrant.EXPECT().
		SignClientEnvelopeToSelf(mock.Anything).
		Return([]byte("signature"), nil).
		Maybe()
	mockRegistrant.EXPECT().NodeID().Return(uint32(1)).Maybe()

	verifier := payerreportMocks.NewMockIPayerReportVerifier(t)
	worker := NewAttestationWorker(ctx, log, mockRegistrant, store, pollInterval, domainSeparator)
	worker.verifier = verifier

	return worker, store, mockRegistrant, verifier
}

func storeReport(
	t *testing.T,
	store *payerreport.Store,
	report *payerreport.PayerReport,
) *payerreport.PayerReportWithStatus {
	id, err := store.StoreReport(t.Context(), report)
	require.NoError(t, err)
	require.NotNil(t, id)
	reportWithStatus, err := store.FetchReport(t.Context(), *id)
	require.NoError(t, err)

	return reportWithStatus
}

func setReportAttestationStatus(
	t *testing.T,
	store payerreport.IPayerReportStore,
	id payerreport.ReportID,
	attestationStatus payerreport.AttestationStatus,
) {
	require.NoError(
		t,
		store.SetReportAttestationStatus(
			t.Context(),
			id,
			[]payerreport.AttestationStatus{payerreport.AttestationPending},
			attestationStatus,
		),
	)
}

func TestFindReport(t *testing.T) {
	worker, store, _, _ := testAttestationWorker(t, time.Second)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: 1,
		StartSequenceID:  1,
		EndSequenceID:    10,
		DomainSeparator:  domainSeparator,
		NodeIDs:          []uint32{1},
	})
	require.NoError(t, err)
	storedReport := storeReport(t, store, &report.PayerReport)

	reports, err := worker.findReportsNeedingAttestation()
	require.NoError(t, err)
	require.Len(t, reports, 1)
	require.Equal(t, storedReport.ID, reports[0].ID)

	setReportAttestationStatus(t, store, storedReport.ID, payerreport.AttestationApproved)

	reports, err = worker.findReportsNeedingAttestation()
	require.NoError(t, err)
	require.Len(t, reports, 0)
}

func TestAttestFirstReport(t *testing.T) {
	worker, store, _, mockVerifier := testAttestationWorker(t, time.Second)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: 1,
		StartSequenceID:  0,
		EndSequenceID:    10,
		NodeIDs:          []uint32{1},
		DomainSeparator:  domainSeparator,
	})
	require.NoError(t, err)
	storedReport := storeReport(t, store, &report.PayerReport)
	require.NoError(t, err)

	mockVerifier.EXPECT().
		IsValidReport(mock.Anything, (*payerreport.PayerReport)(nil), &report.PayerReport).
		Return(true, nil)

	err = worker.attestReport(storedReport)
	require.NoError(t, err)

	fromDB, err := store.FetchReport(t.Context(), storedReport.ID)
	require.NoError(t, err)
	require.Equal(
		t,
		payerreport.AttestationStatus(payerreport.AttestationApproved),
		fromDB.AttestationStatus,
	)
}

package workers

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
	payerreportMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/payerreport"
	registrantMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/registrant"
)

const originatorNodeID = 100

var domainSeparator = common.BytesToHash(testutils.RandomBytes(32))

func testAttestationWorker(
	t *testing.T,
	pollInterval time.Duration,
) (*AttestationWorker, *payerreport.Store, *registrantMocks.MockIRegistrant, *payerreportMocks.MockIPayerReportVerifier) {
	var (
		log            = testutils.NewLog(t)
		ctx            = t.Context()
		db, _          = testutils.NewDB(t, ctx)
		store          = payerreport.NewStore(log, db)
		mockRegistrant = registrantMocks.NewMockIRegistrant(t)
		verifier       = payerreportMocks.NewMockIPayerReportVerifier(t)
		worker         = NewAttestationWorker(
			ctx,
			log,
			mockRegistrant,
			store,
			pollInterval,
			domainSeparator,
		)
	)

	worker.verifier = verifier

	mockRegistrant.EXPECT().
		SignPayerReportAttestation(mock.Anything).
		Return(&payerreport.NodeSignature{
			Signature: []byte("signature"),
			NodeID:    originatorNodeID,
		}, nil).
		Maybe()

	mockRegistrant.EXPECT().
		SignClientEnvelopeToSelf(mock.Anything).
		Return([]byte("signature"), nil).
		Maybe()

	mockRegistrant.EXPECT().NodeID().Return(uint32(1)).Maybe()

	// Create real sequence IDs, from 1 to 10.
	for i := uint64(1); i <= 10; i++ {
		testutils.InsertGatewayEnvelopes(t, db.DB(), []queries.InsertGatewayEnvelopeParams{
			{
				OriginatorNodeID:     originatorNodeID,
				OriginatorSequenceID: int64(i),
				OriginatorEnvelope:   testutils.RandomBytes(32),
				Topic:                testutils.RandomBytes(32),
			},
		})
	}

	return worker, store, mockRegistrant, verifier
}

func storeReport(
	t *testing.T,
	store *payerreport.Store,
	report *payerreport.PayerReport,
) *payerreport.PayerReportWithStatus {
	numRows, err := store.StoreReport(t.Context(), report)
	require.NoError(t, err)
	require.Equal(t, int64(1), numRows)
	reportWithStatus, err := store.FetchReport(t.Context(), report.ID)
	require.NoError(t, err)

	return reportWithStatus
}

func TestFindReport(t *testing.T) {
	worker, store, _, _ := testAttestationWorker(t, time.Second)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  1,
		EndSequenceID:    10,
		DomainSeparator:  domainSeparator,
		NodeIDs:          []uint32{originatorNodeID},
	})
	require.NoError(t, err)
	storedReport := storeReport(t, store, &report.PayerReport)

	reports, err := worker.findReportsNeedingAttestation()
	require.NoError(t, err)
	require.Len(t, reports, 1)
	require.Equal(t, storedReport.ID, reports[0].ID)

	require.NoError(
		t,
		store.SetReportAttestationApproved(
			t.Context(),
			storedReport.ID,
		),
	)

	reports, err = worker.findReportsNeedingAttestation()
	require.NoError(t, err)
	require.Empty(t, reports)
}

func TestDontAttestReportsIfSeqIDNotProcessed(t *testing.T) {
	worker, store, _, _ := testAttestationWorker(t, time.Second)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  1,
		EndSequenceID:    20,
		DomainSeparator:  domainSeparator,
		NodeIDs:          []uint32{originatorNodeID},
	})
	require.NoError(t, err)
	_ = storeReport(t, store, &report.PayerReport)

	reports, err := worker.findReportsNeedingAttestation()
	require.NoError(t, err)
	require.Empty(t, reports)
}

// TestDontAttestReportsInNonPendingStates covers the following states:
// (0,1), (0,2)        - Already attested (approved/rejected)
// (1,0), (1,1), (1,2) - Already submitted
// (2,0), (2,1), (2,2) - Already settled
// (3,1)               - Submission rejected
func TestDontAttestReportsInNonPendingStates(t *testing.T) {
	type reportState int
	const (
		stateAttestationApproved reportState = iota
		stateAttestationRejected
		stateSubmissionSubmitted
		stateSubmissionRejected
		stateSubmissionSettled
	)

	prepareReport := func(
		t *testing.T,
		state reportState,
		store payerreport.IPayerReportStore,
		r *payerreport.PayerReportWithStatus,
	) {
		t.Helper()

		switch state {
		case stateAttestationApproved:
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		case stateAttestationRejected:
			require.NoError(t, store.SetReportAttestationRejected(t.Context(), r.ID))

		case stateSubmissionSubmitted:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID, 0))

		case stateSubmissionRejected:
			require.NoError(t, store.SetReportSubmissionRejected(t.Context(), r.ID))

		case stateSubmissionSettled:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID, 0))
			require.NoError(t, store.SetReportSettled(t.Context(), r.ID))

		default:
			t.Fatalf("unknown target state: %v", state)
		}
	}

	testCases := []struct {
		name                      string
		state                     reportState
		expectedAttestationStatus payerreport.AttestationStatus
		expectedSubmissionStatus  payerreport.SubmissionStatus
	}{
		{
			"don't attest already attested report",
			stateAttestationApproved,
			payerreport.AttestationApproved,
			payerreport.SubmissionPending,
		},
		{
			"don't attest already attested report",
			stateAttestationRejected,
			payerreport.AttestationRejected,
			payerreport.SubmissionPending,
		},
		{
			"don't attest already submitted report",
			stateSubmissionSubmitted,
			payerreport.AttestationPending,
			payerreport.SubmissionSubmitted,
		},
		{
			"don't attest rejected report",
			stateSubmissionRejected,
			payerreport.AttestationPending,
			payerreport.SubmissionRejected,
		},
		{
			"don't attest already settled report",
			stateSubmissionSettled,
			payerreport.AttestationPending,
			payerreport.SubmissionSettled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker, store, _, _ := testAttestationWorker(t, time.Second)

			report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
				OriginatorNodeID: originatorNodeID,
				StartSequenceID:  0,
				EndSequenceID:    10,
				DomainSeparator:  domainSeparator,
				NodeIDs:          []uint32{originatorNodeID},
			})
			require.NoError(t, err)

			stored := storeReport(t, store, &report.PayerReport)

			prepareReport(t, tc.state, store, stored)

			reports, err := worker.findReportsNeedingAttestation()
			require.NoError(t, err)

			// Every single case should not find any reports needing attestation.
			require.Empty(t, reports)

			got, err := store.FetchReport(t.Context(), stored.ID)
			require.NoError(t, err)
			require.Equal(t, tc.expectedAttestationStatus, got.AttestationStatus)
			require.Equal(t, tc.expectedSubmissionStatus, got.SubmissionStatus)
		})
	}
}

func TestAttestFirstReport(t *testing.T) {
	worker, store, _, mockVerifier := testAttestationWorker(t, time.Second)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID: originatorNodeID,
		StartSequenceID:  0,
		EndSequenceID:    10,
		NodeIDs:          []uint32{originatorNodeID},
		DomainSeparator:  domainSeparator,
	})
	require.NoError(t, err)
	storedReport := storeReport(t, store, &report.PayerReport)
	require.NoError(t, err)

	mockVerifier.EXPECT().
		VerifyReport(mock.Anything, (*payerreport.PayerReport)(nil), &report.PayerReport).
		Return(payerreport.VerifyReportResult{
			IsValid: true,
			Reason:  "valid report",
		}, nil)

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

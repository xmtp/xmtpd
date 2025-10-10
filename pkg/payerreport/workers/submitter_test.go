package workers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	registrantMocks "github.com/xmtp/xmtpd/pkg/mocks/registrant"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/registry"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

func testSubmitterWorker(
	t *testing.T,
) (*SubmitterWorker, *payerreport.Store, *blockchainMocks.MockPayerReportsManager) {
	var (
		log            = testutils.NewLog(t)
		ctx            = t.Context()
		db, _          = testutils.NewDB(t, ctx)
		store          = payerreport.NewStore(db, log)
		mockRegistrant = registrantMocks.NewMockIRegistrant(t)
		registry       = mocks.NewMockNodeRegistry(t)
		reportsManager = blockchainMocks.NewMockPayerReportsManager(t)
		myNodeID       = uint32(1)
	)

	mockRegistrant.EXPECT().
		SignPayerReportAttestation(mock.Anything).
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

	worker := NewSubmitterWorker(
		ctx,
		log,
		store,
		registry,
		reportsManager,
		myNodeID,
	)

	return worker, store, reportsManager
}

// TestSubmitReports covers all possible valid and invalid states for the submitter worker.
//
// Valid transitions:
//
//	(0,1) → (1,1) - SubmitterWorker succeeds
//	(0,1) → (3,1) - SubmitterWorker hits race condition
//
// Invalid transitions:
//
//	(0,0) - Not attested
//	(0,2) - Rejected
//	(1,X) - Already submitted
//	(2,X) - Already settled
//	(3,X) - Submission rejected
func TestSubmitterStatesAndTransitions(t *testing.T) {
	type reportState int
	const (
		stateSubmissionPendingAttestationPending reportState = iota
		stateSubmissionPendingAttestationApprovedSuccess
		stateSubmissionPendingAttestationApprovedRaceCondition
		stateSubmissionPendingAttestationRejected
		stateSubmissionSubmittedAttestationPending
		stateSubmissionSubmittedAttestationApproved
		stateSubmissionSubmittedAttestationRejected
		stateSubmissionSettledAttestationPending
		stateSubmissionSettledAttestationApproved
		stateSubmissionSettledAttestationRejected
		stateSubmissionRejectedAttestationPending
		stateSubmissionRejectedAttestationApproved
		stateSubmissionRejectedAttestationRejected
	)

	prepareReport := func(
		t *testing.T,
		state reportState,
		store payerreport.IPayerReportStore,
		r *payerreport.PayerReportWithStatus,
	) {
		t.Helper()

		switch state {
		//	(0,0) - Not attested
		case stateSubmissionPendingAttestationPending:
			// Nothing to do!

		//	(0,1) → (1,1) - SubmitterWorker succeeds
		case stateSubmissionPendingAttestationApprovedSuccess:
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		//	(0,1) → (3,1) - SubmitterWorker hits race condition
		case stateSubmissionPendingAttestationApprovedRaceCondition:
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		//	(0,2) - Rejected
		case stateSubmissionPendingAttestationRejected:
			require.NoError(t, store.SetReportAttestationRejected(t.Context(), r.ID))

		//	(1,0) - Already submitted, attestation pending
		case stateSubmissionSubmittedAttestationPending:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))

		//	(1,1) - Already submitted, attestation approved
		case stateSubmissionSubmittedAttestationApproved:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		//	(1,2) - Already submitted, attestation rejected
		case stateSubmissionSubmittedAttestationRejected:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationRejected(t.Context(), r.ID))

		//	(2,0) - Already settled
		case stateSubmissionSettledAttestationPending:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))
			require.NoError(t, store.SetReportSettled(t.Context(), r.ID))

		//	(2,1) - Already settled, attestation approved
		case stateSubmissionSettledAttestationApproved:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))
			require.NoError(t, store.SetReportSettled(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		//	(2,2) - Already settled, attestation rejected
		case stateSubmissionSettledAttestationRejected:
			require.NoError(t, store.SetReportSubmitted(t.Context(), r.ID))
			require.NoError(t, store.SetReportSettled(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationRejected(t.Context(), r.ID))

		//	(3,0) - Submission rejected, attestation pending
		case stateSubmissionRejectedAttestationPending:
			require.NoError(t, store.SetReportSubmissionRejected(t.Context(), r.ID))

		//	(3,1) - Submission rejected, attestation approved
		case stateSubmissionRejectedAttestationApproved:
			require.NoError(t, store.SetReportSubmissionRejected(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationApproved(t.Context(), r.ID))

		//	(3,2) - Submission rejected, attestation rejected
		case stateSubmissionRejectedAttestationRejected:
			require.NoError(t, store.SetReportSubmissionRejected(t.Context(), r.ID))
			require.NoError(t, store.SetReportAttestationRejected(t.Context(), r.ID))

		default:
			t.Fatalf("unknown target state: %v", state)
		}
	}

	testCases := []struct {
		name                      string
		state                     reportState
		expectedAttestationStatus payerreport.AttestationStatus
		expectedSubmissionStatus  payerreport.SubmissionStatus
		wantSubmitRejected        bool
	}{
		//	(0,0) - Not attested
		{
			"don't submit not attested report",
			stateSubmissionPendingAttestationPending,
			payerreport.AttestationPending,
			payerreport.SubmissionPending,
			false,
		},
		//	(0,1) → (1,1) - SubmitterWorker succeeds
		{
			"submit a pending, approved report",
			stateSubmissionPendingAttestationApprovedSuccess,
			payerreport.AttestationApproved,
			payerreport.SubmissionPending,
			false,
		},
		//	(0,1) → (3,1) - SubmitterWorker hits race condition:
		// - Some other node submitted a valid report that goes before the one submitted.
		// - The contract throws InvalidSequenceIDs or InvalidStartSequenceID.
		{
			"submit a pending, approved report, reject on chain",
			stateSubmissionPendingAttestationApprovedRaceCondition,
			payerreport.AttestationApproved,
			payerreport.SubmissionRejected,
			true,
		},
		//	(0,2) - Rejected
		{
			"don't submit rejected report",
			stateSubmissionPendingAttestationRejected,
			payerreport.AttestationRejected,
			payerreport.SubmissionPending,
			false,
		},
		//	(1,0) - Already submitted, attestation pending
		{
			"don't submit already submitted report",
			stateSubmissionSubmittedAttestationPending,
			payerreport.AttestationPending,
			payerreport.SubmissionSubmitted,
			false,
		},
		//	(1,1) - Already submitted, attestation approved
		{
			"don't submit already submitted report",
			stateSubmissionSubmittedAttestationApproved,
			payerreport.AttestationApproved,
			payerreport.SubmissionSubmitted,
			false,
		},
		//	(1,2) - Already submitted, attestation rejected
		{
			"don't submit already submitted report",
			stateSubmissionSubmittedAttestationRejected,
			payerreport.AttestationRejected,
			payerreport.SubmissionSubmitted,
			false,
		},
		//	(2,0) - Already settled, attestation pending
		{
			"don't submit already settled report",
			stateSubmissionSettledAttestationPending,
			payerreport.AttestationPending,
			payerreport.SubmissionSettled,
			false,
		},
		//	(2,1) - Already settled, attestation approved
		{
			"don't submit already settled report",
			stateSubmissionSettledAttestationApproved,
			payerreport.AttestationApproved,
			payerreport.SubmissionSettled,
			false,
		},
		//	(2,2) - Already settled, attestation rejected
		{
			"don't submit already settled report",
			stateSubmissionSettledAttestationRejected,
			payerreport.AttestationRejected,
			payerreport.SubmissionSettled,
			false,
		},
		//	(3,0) - Submission rejected, attestation pending
		{
			"don't submit rejected report",
			stateSubmissionRejectedAttestationPending,
			payerreport.AttestationPending,
			payerreport.SubmissionRejected,
			false,
		},
		//	(3,1) - Submission rejected, attestation approved
		{
			"don't submit rejected report",
			stateSubmissionRejectedAttestationApproved,
			payerreport.AttestationApproved,
			payerreport.SubmissionRejected,
			false,
		},
		//	(3,2) - Submission rejected, attestation rejected
		{
			"don't submit rejected report",
			stateSubmissionRejectedAttestationRejected,
			payerreport.AttestationRejected,
			payerreport.SubmissionRejected,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker, store, reportsManager := testSubmitterWorker(t)

			report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
				OriginatorNodeID: 1,
				StartSequenceID:  0,
				EndSequenceID:    10,
				DomainSeparator:  domainSeparator,
				NodeIDs:          []uint32{1},
			})
			require.NoError(t, err)

			stored := storeReport(t, store, &report.PayerReport)

			prepareReport(t, tc.state, store, stored)

			if tc.wantSubmitRejected {
				require.NoError(
					t,
					store.StoreAttestation(t.Context(), &payerreport.PayerReportAttestation{
						Report: &report.PayerReport,
						NodeSignature: payerreport.NodeSignature{
							Signature: []byte("0xSignature"),
							NodeID:    1,
						},
					}),
				)

				reportsManager.EXPECT().
					SubmitPayerReport(mock.Anything, mock.Anything).
					Return(blockchain.NewBlockchainError(fmt.Errorf("execution reverted: 0x84e23433")))
			}

			err = worker.SubmitReports(t.Context())
			require.NoError(t, err)

			got, err := store.FetchReport(t.Context(), stored.ID)
			require.NoError(t, err)
			require.Equal(t, tc.expectedAttestationStatus, got.AttestationStatus)
			require.Equal(t, tc.expectedSubmissionStatus, got.SubmissionStatus)
		})
	}
}

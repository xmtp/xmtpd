package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	registrantMocks "github.com/xmtp/xmtpd/pkg/mocks/registrant"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

var ErrTestReportGenerationFailed = errors.New("test: report generation failed")

type TrackingPayerReportGenerator struct {
	t                *testing.T
	reportsGenerated int
}

var _ payerreport.IPayerReportGenerator = &TrackingPayerReportGenerator{}

func (t *TrackingPayerReportGenerator) GenerateReport(
	ctx context.Context,
	params payerreport.PayerReportGenerationParams,
) (*payerreport.PayerReportWithInputs, error) {
	t.t.Log("Generating mock report and erroring out")
	t.reportsGenerated++
	return nil, ErrTestReportGenerationFailed
}

func (t *TrackingPayerReportGenerator) GetReportsGenerated() int {
	return t.reportsGenerated
}

func newTestGenerator(
	t *testing.T,
) (*GeneratorWorker, *payerreport.Store, *registrantMocks.MockIRegistrant, *TrackingPayerReportGenerator) {
	log := testutils.NewLog(t)
	ctx := t.Context()
	db, _ := testutils.NewDB(t, ctx)
	store := payerreport.NewStore(log, db)

	mockRegistrant := registrantMocks.NewMockIRegistrant(t)

	mockRegistrant.EXPECT().NodeID().Return(originatorNodeID).Maybe()

	generator := TrackingPayerReportGenerator{
		t: t,
	}

	worker := &GeneratorWorker{
		ctx:                  ctx,
		logger:               log.With(zap.String("test", "generator")),
		store:                store,
		registrant:           mockRegistrant,
		generator:            &generator,
		generateSelfPeriod:   0, // always allow generation
		generateOthersPeriod: 0,
		expirySelfPeriod:     24 * time.Hour, // no expiration unless we force it
		expiryOthersPeriod:   24 * time.Hour,
	}

	return worker, store, mockRegistrant, &generator
}

func TestGenerator_NoDuplicateWhenReportAlreadyExists(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(time.Now().Unix() / 60)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), report.ID, 0))

	report, err = payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     10,
		EndSequenceID:       20,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)

	err = worker.maybeGenerateReport(originatorNodeID)
	require.NoError(t, err)
	require.Equal(t, 0, generator.GetReportsGenerated())
}

func TestGenerator_DoGenerate(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), report.ID, 0))

	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())
}

func TestGenerator_ExpireReport(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)

	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())

	reportWithStatus, err := store.FetchReport(t.Context(), report.ID)
	require.NoError(t, err)
	require.EqualValues(t, payerreport.SubmissionRejected, reportWithStatus.SubmissionStatus)

	// a new reports get generated to replace the expired one
	require.Equal(t, 1, generator.GetReportsGenerated())
}

func TestGenerator_ExpiredWithPreExisting(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), report.ID, 0))

	report, err = payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     10,
		EndSequenceID:       20,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)

	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())
}

func TestGenerator_ExpirationDoesNotTouchSubmitted(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), report.ID, 0))

	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())

	reportWithStatus, err := store.FetchReport(t.Context(), report.ID)
	require.NoError(t, err)
	require.EqualValues(t, payerreport.SubmissionSubmitted, reportWithStatus.SubmissionStatus)
}

func TestGenerator_ExpirationDoesNotTouchSettled(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSettled(t.Context(), report.ID))

	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())

	reportWithStatus, err := store.FetchReport(t.Context(), report.ID)
	require.NoError(t, err)
	require.EqualValues(t, payerreport.SubmissionSettled, reportWithStatus.SubmissionStatus)
}

func TestGenerator_FutureMinuteGetsSkipped(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	futureMinute := uint32(time.Now().Add(5*time.Minute).Unix() / 60)

	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: futureMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), report.ID, 0))

	report, err = payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     10,
		EndSequenceID:       20,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: futureMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &report.PayerReport)

	// nothing should get generated, the future minute is within the current generation
	err = worker.maybeGenerateReport(originatorNodeID)
	require.NoError(t, err)
	require.Equal(t, 0, generator.GetReportsGenerated())
}

func TestGenerator_FirstReportFromScratch(t *testing.T) {
	worker, _, _, generator := newTestGenerator(t)

	// No reports stored at all.

	err := worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())
}

func TestGenerator_IgnoresRejectedReportsAndGeneratesNew(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(time.Now().Unix() / 60)

	// Store a report that will be marked as rejected.
	report, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	stored := storeReport(t, store, &report.PayerReport)

	require.NoError(t, store.SetReportSubmissionRejected(t.Context(), stored.ID))

	// Snapshot shouldn't see this rejected report, so generator should act
	// as if there are no reports and generate from 0.
	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())
}

func TestGenerator_OverlappingNonBoundaryReportDoesNotBlockGeneration(t *testing.T) {
	worker, store, _, generator := newTestGenerator(t)
	currentMinute := uint32(10)

	// Submitted report: [0, 10]
	r1, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     0,
		EndSequenceID:       10,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &r1.PayerReport)
	require.NoError(t, store.SetReportSubmitted(t.Context(), r1.ID, 0))

	// Overlapping pending report: [5, 15] – does NOT start at 10.
	r2, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    originatorNodeID,
		StartSequenceID:     5,
		EndSequenceID:       15,
		DomainSeparator:     domainSeparator,
		NodeIDs:             []uint32{originatorNodeID},
		EndMinuteSinceEpoch: currentMinute,
	})
	require.NoError(t, err)
	storeReport(t, store, &r2.PayerReport)

	// Since there is no report starting exactly at 10, generator should
	// still try to generate [10, ...].
	err = worker.maybeGenerateReport(originatorNodeID)
	require.ErrorIs(t, err, ErrTestReportGenerationFailed)
	require.Equal(t, 1, generator.GetReportsGenerated())
}

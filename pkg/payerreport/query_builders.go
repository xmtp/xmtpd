package payerreport

import (
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type FetchReportsQuery struct {
	SubmissionStatusIn  []SubmissionStatus
	AttestationStatusIn []AttestationStatus
	StartSequenceID     *uint64
	EndSequenceID       *uint64
	CreatedAfter        time.Time
	OriginatorNodeID    *uint32
	PayerReportID       []byte
	MinAttestations     *int32
}

func (f *FetchReportsQuery) toParams() queries.FetchPayerReportsParams {
	return queries.FetchPayerReportsParams{
		CreatedAfter:        utils.NewNullTime(f.CreatedAfter),
		SubmissionStatusIn:  utils.NewNullInt16Slice(f.SubmissionStatusIn),
		AttestationStatusIn: utils.NewNullInt16Slice(f.AttestationStatusIn),
		StartSequenceID:     utils.NewNullInt64(f.StartSequenceID),
		EndSequenceID:       utils.NewNullInt64(f.EndSequenceID),
		OriginatorNodeID:    utils.NewNullInt32(f.OriginatorNodeID),
		PayerReportID:       utils.NewNullBytes(f.PayerReportID),
		MinAttestations:     utils.NewNullInt32(f.MinAttestations),
	}
}

func NewFetchReportsQuery() *FetchReportsQuery {
	return &FetchReportsQuery{}
}

func (f *FetchReportsQuery) WithSubmissionStatus(statuses ...SubmissionStatus) *FetchReportsQuery {
	f.SubmissionStatusIn = append(f.SubmissionStatusIn, statuses...)
	return f
}

func (f *FetchReportsQuery) WithAttestationStatus(
	statuses ...AttestationStatus,
) *FetchReportsQuery {
	f.AttestationStatusIn = append(f.AttestationStatusIn, statuses...)
	return f
}

func (f *FetchReportsQuery) WithCreatedAfter(createdAfter time.Time) *FetchReportsQuery {
	f.CreatedAfter = createdAfter.UTC()
	return f
}

func (f *FetchReportsQuery) WithReportID(payerReportID ReportID) *FetchReportsQuery {
	f.PayerReportID = payerReportID[:]
	return f
}

func (f *FetchReportsQuery) WithStartSequenceID(startSequenceID uint64) *FetchReportsQuery {
	f.StartSequenceID = &startSequenceID
	return f
}

func (f *FetchReportsQuery) WithEndSequenceID(endSequenceID uint64) *FetchReportsQuery {
	f.EndSequenceID = &endSequenceID
	return f
}

func (f *FetchReportsQuery) WithOriginatorNodeID(originatorNodeID uint32) *FetchReportsQuery {
	f.OriginatorNodeID = &originatorNodeID
	return f
}

func (f *FetchReportsQuery) WithMinAttestations(minAttestations int32) *FetchReportsQuery {
	f.MinAttestations = &minAttestations
	return f
}

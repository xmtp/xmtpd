package metadata

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
)

type PayerInfoGroupBy string

const (
	PayerInfoGroupByHour PayerInfoGroupBy = "hour"
	PayerInfoGroupByDay  PayerInfoGroupBy = "day"
)

type IPayerInfoFetcher interface {
	GetPayerByAddress(ctx context.Context, address string) (int32, error)
	GetPayerInfo(
		ctx context.Context,
		payerID int32,
		groupBy PayerInfoGroupBy,
	) (*metadata_api.GetPayerInfoResponse_PayerInfo, error)
}

type PayerInfoFetcher struct {
	queries *queries.Queries
}

func NewPayerInfoFetcher(db *sql.DB) *PayerInfoFetcher {
	return &PayerInfoFetcher{
		queries: queries.New(db),
	}
}

// Gets the total spend and message count for a payer between the two timestamps, grouped by the appropriate granularity.
func (f *PayerInfoFetcher) GetPayerInfo(
	ctx context.Context,
	payerID int32,
	groupBy PayerInfoGroupBy,
) (*metadata_api.GetPayerInfoResponse_PayerInfo, error) {
	result, err := f.queries.GetPayerInfoReport(ctx, queries.GetPayerInfoReportParams{
		GroupBy: string(groupBy),
		PayerID: payerID,
	})
	if err != nil {
		return nil, err
	}

	payerInfo := &metadata_api.GetPayerInfoResponse_PayerInfo{
		PeriodSummaries: make([]*metadata_api.GetPayerInfoResponse_PeriodSummary, 0),
	}

	for _, row := range result {
		periodSummary := &metadata_api.GetPayerInfoResponse_PeriodSummary{
			AmountSpentPicodollars: uint64(row.TotalSpendPicodollars),
			NumMessages:            uint64(row.TotalMessageCount),
			PeriodStartUnixSeconds: uint64(row.TimePeriod),
		}

		payerInfo.PeriodSummaries = append(payerInfo.PeriodSummaries, periodSummary)
	}

	return payerInfo, nil
}

// GetPayerByAddress looks up a payer ID by address
func (f *PayerInfoFetcher) GetPayerByAddress(ctx context.Context, address string) (int32, error) {
	return f.queries.GetPayerByAddress(ctx, address)
}

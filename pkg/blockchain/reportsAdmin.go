package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"go.uber.org/zap"
)

type ReportsAdmin struct {
	client *ethclient.Client
	signer TransactionSigner
	log    *zap.Logger
}

func NewReportsAdmin(
	log *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
) *ReportsAdmin {
	return &ReportsAdmin{
		log:    log,
		client: client,
		signer: signer,
	}
}

func (r *ReportsAdmin) SubmitPayerReport(
	ctx context.Context,
	report *payerreport.PayerReportWithStatus,
) error {
	return nil
}

func (r *ReportsAdmin) GetDomainSeparator(ctx context.Context) (common.Hash, error) {
	return common.Hash{}, nil
}

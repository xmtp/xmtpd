package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
)

// PayerReportManagerContract is an interface for the PayerReportManager contract
type PayerReportManagerContract interface {
	// Parse events
	ParsePayerReportSubmitted(log types.Log) (*p.PayerReportManagerPayerReportSubmitted, error)
	ParsePayerReportSubsetSettled(
		log types.Log,
	) (*p.PayerReportManagerPayerReportSubsetSettled, error)

	// Contract calls
	GetPayerReport(
		opts *bind.CallOpts,
		originatorNodeID uint32,
		payerReportIndex *big.Int,
	) (p.IPayerReportManagerPayerReport, error)
	DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error)
}

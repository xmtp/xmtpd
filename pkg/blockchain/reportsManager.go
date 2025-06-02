package blockchain

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	reportManager "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"go.uber.org/zap"
)

type ReportsManager struct {
	client                *ethclient.Client
	signer                TransactionSigner
	log                   *zap.Logger
	reportManagerContract *reportManager.PayerReportManager
	domainSeparator       *common.Hash
	domainSeparatorLock   sync.Mutex
}

func NewReportsManager(
	log *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	options config.SettlementChainOptions,
) (*ReportsManager, error) {
	reportManagerContract, err := reportManager.NewPayerReportManager(
		common.HexToAddress(options.PayerReportManagerAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &ReportsManager{
		log:                   log,
		client:                client,
		signer:                signer,
		reportManagerContract: reportManagerContract,
		domainSeparator:       nil,
	}, nil
}

func (r *ReportsManager) SubmitPayerReport(
	ctx context.Context,
	report *payerreport.PayerReportWithStatus,
) error {
	err := ExecuteTransaction(
		ctx,
		r.signer,
		r.log,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			signatures := make(
				[]reportManager.IPayerReportManagerPayerReportSignature,
				len(report.AttestationSignatures),
			)

			for i, signature := range report.AttestationSignatures {
				// The smart contract expects the signature to use the legacy EIP-155 recovery ID
				// which is the last byte of the signature. We need to add 27 to the last byte
				// to get the correct recovery ID.
				sig := make([]byte, len(signature.Signature))
				copy(sig, signature.Signature)
				if sig[64] < 27 {
					sig[64] += 27
				}

				signatures[i] = reportManager.IPayerReportManagerPayerReportSignature{
					NodeId:    signature.NodeID,
					Signature: sig,
				}
			}

			r.log.Info(
				"submitting report",
				zap.Any("report", report),
				zap.Any("signatures", signatures),
			)

			return r.reportManagerContract.Submit(
				opts,
				report.OriginatorNodeID,
				report.StartSequenceID,
				report.EndSequenceID,
				report.PayersMerkleRoot,
				report.ActiveNodeIDs,
				signatures,
			)
		},
		func(log *types.Log) (any, error) {
			return r.reportManagerContract.ParsePayerReportSubmitted(*log)
		},
		func(event any) {
			if parsedEvent, ok := event.(*reportManager.PayerReportManagerPayerReportSubmitted); ok {
				r.log.Info(
					"payer report submitted",
					zap.Any("event", parsedEvent),
				)
			} else {
				r.log.Warn("unknown event type", zap.Any("event", event))
			}
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReportsManager) GetReport(
	ctx context.Context,
	originatorNodeID uint32,
	index uint64,
) (*payerreport.PayerReport, error) {
	report, err := r.reportManagerContract.GetPayerReport(&bind.CallOpts{
		Context: ctx,
	}, originatorNodeID, big.NewInt(int64(index)))
	if err != nil {
		return nil, err
	}

	domainSeparator, err := r.GetDomainSeparator(ctx)
	if err != nil {
		return nil, err
	}

	return transformOnChainReport(&report, originatorNodeID, domainSeparator)
}

func (r *ReportsManager) GetReportID(
	ctx context.Context,
	payerReport *payerreport.PayerReportWithStatus,
) (payerreport.ReportID, error) {
	digest, err := r.reportManagerContract.GetPayerReportDigest(&bind.CallOpts{
		Context: ctx,
	}, payerReport.OriginatorNodeID, payerReport.StartSequenceID, payerReport.EndSequenceID, payerReport.PayersMerkleRoot, payerReport.ActiveNodeIDs)
	if err != nil {
		return payerreport.ReportID{}, err
	}
	return payerreport.ReportID(digest), nil
}

func (r *ReportsManager) GetDomainSeparator(ctx context.Context) (common.Hash, error) {
	r.domainSeparatorLock.Lock()
	defer r.domainSeparatorLock.Unlock()

	if r.domainSeparator == nil {
		domainSeparator, err := r.reportManagerContract.DOMAINSEPARATOR(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return common.Hash{}, err
		}
		asHash := common.Hash(domainSeparator)
		r.domainSeparator = &asHash
	}
	return *r.domainSeparator, nil
}

func transformOnChainReport(
	report *reportManager.IPayerReportManagerPayerReport,
	nodeID uint32,
	domainSeparator common.Hash,
) (*payerreport.PayerReport, error) {
	id, err := payerreport.BuildPayerReportID(
		nodeID,
		report.StartSequenceId,
		report.EndSequenceId,
		report.PayersMerkleRoot,
		report.NodeIds,
		domainSeparator,
	)
	if err != nil {
		return nil, err
	}

	return &payerreport.PayerReport{
		ID:               *id,
		OriginatorNodeID: nodeID,
		StartSequenceID:  report.StartSequenceId,
		EndSequenceID:    report.EndSequenceId,
		PayersMerkleRoot: report.PayersMerkleRoot,
		ActiveNodeIDs:    report.NodeIds,
	}, nil
}

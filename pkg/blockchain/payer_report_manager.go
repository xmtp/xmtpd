package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	payerRegistry "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	reportManager "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/merkle"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type ReportsManager struct {
	client                *ethclient.Client
	signer                TransactionSigner
	logger                *zap.Logger
	reportManagerContract *reportManager.PayerReportManager
	payerRegistryContract *payerRegistry.PayerRegistry
	domainSeparator       *common.Hash
	domainSeparatorLock   sync.Mutex
}

type SettlementSummary struct {
	Offset    uint32
	IsSettled bool
}

func NewReportsManager(
	logger *zap.Logger,
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

	payerRegistryContract, err := payerRegistry.NewPayerRegistry(
		common.HexToAddress(options.PayerRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &ReportsManager{
		logger:                logger.Named(utils.PayerReportManagerAdminLoggerName),
		client:                client,
		signer:                signer,
		reportManagerContract: reportManagerContract,
		payerRegistryContract: payerRegistryContract,
		domainSeparator:       nil,
	}, nil
}

func (r *ReportsManager) SubmitPayerReport(
	ctx context.Context,
	report *payerreport.PayerReportWithStatus,
) (int32, ProtocolError) {
	var reportIndex int32
	var eventErr ProtocolError
	var foundEvent bool
	err := ExecuteTransaction(
		ctx,
		r.signer,
		r.logger,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			signatures := prepareSignatures(report.AttestationSignatures)

			if r.logger.Core().Enabled(zap.DebugLevel) {
				r.logger.Debug(
					"submitting report",
					utils.PayerReportIDField(report.ID.String()),
					utils.BodyField(signatures),
				)
			}

			return r.reportManagerContract.Submit(
				opts,
				report.OriginatorNodeID,
				report.StartSequenceID,
				report.EndSequenceID,
				report.EndMinuteSinceEpoch,
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
				foundEvent = true
				var err error
				reportIndex, err = payerreport.ValidateReportIndex(parsedEvent.PayerReportIndex)
				if err != nil {
					eventErr = NewBlockchainError(err)
					r.logger.Error(
						"payer report index validation failed",
						zap.Error(err),
					)
					return
				}

				r.logger.Info(
					"payer report submitted",
					utils.PayerReportIDField(report.ID.String()),
				)
			} else {
				r.logger.Warn("unknown event type")
			}
		},
	)
	if err != nil {
		return 0, err
	}
	if eventErr != nil {
		return 0, eventErr
	}
	if !foundEvent {
		return 0, NewBlockchainError(errors.New("no PayerReportSubmitted event found"))
	}
	return reportIndex, nil
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

	return transformOnChainReport(&report, originatorNodeID, domainSeparator, index)
}

func (r *ReportsManager) GetReportID(
	ctx context.Context,
	payerReport *payerreport.PayerReportWithStatus,
) (payerreport.ReportID, error) {
	digest, err := r.reportManagerContract.GetPayerReportDigest(
		&bind.CallOpts{
			Context: ctx,
		},
		payerReport.OriginatorNodeID,
		payerReport.StartSequenceID,
		payerReport.EndSequenceID,
		payerReport.EndMinuteSinceEpoch,
		payerReport.PayersMerkleRoot,
		payerReport.ActiveNodeIDs,
	)
	if err != nil {
		return payerreport.ReportID{}, err
	}
	return digest, nil
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

func (r *ReportsManager) SettlementSummary(
	ctx context.Context,
	originatorNodeID uint32,
	index uint64,
) (*SettlementSummary, error) {
	report, err := r.reportManagerContract.GetPayerReport(&bind.CallOpts{
		Context: ctx,
	}, originatorNodeID, new(big.Int).SetUint64(index))
	if err != nil {
		return nil, err
	}

	return &SettlementSummary{
		Offset:    report.Offset,
		IsSettled: report.IsSettled,
	}, nil
}

func (r *ReportsManager) SettleReport(
	ctx context.Context,
	originatorNodeID uint32,
	index uint64,
	proof *merkle.MultiProof,
) error {
	leaves, proofElements, err := prepareProof(proof)
	if err != nil {
		return err
	}

	return ExecuteTransaction(
		ctx,
		r.signer,
		r.logger,
		r.client,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return r.reportManagerContract.Settle(
				opts,
				originatorNodeID,
				big.NewInt(int64(index)),
				leaves,
				proofElements,
			)
		},
		func(log *types.Log) (any, error) {
			if usageSettled, err := r.payerRegistryContract.ParseUsageSettled(*log); err == nil {
				r.logger.Info(
					"usage settled for",
					zap.String("payer", usageSettled.Payer.String()),
					zap.Any("amount", usageSettled.Amount),
				)
			}

			return r.reportManagerContract.ParsePayerReportSubsetSettled(*log)
		},
		func(event any) {
			if parsedEvent, ok := event.(*reportManager.PayerReportManagerPayerReportSubsetSettled); ok {
				r.logger.Info(
					"settled report",
					utils.OriginatorIDField(originatorNodeID),
					zap.Uint64("index", index),
					zap.Uint32("remaining", parsedEvent.Remaining),
				)
			} else {
				r.logger.Warn("unknown event type")
			}
		},
	)
}

func transformOnChainReport(
	report *reportManager.IPayerReportManagerPayerReport,
	nodeID uint32,
	domainSeparator common.Hash,
	index uint64,
) (*payerreport.PayerReport, error) {
	id, err := payerreport.BuildPayerReportID(
		nodeID,
		report.StartSequenceId,
		report.EndSequenceId,
		report.EndMinuteSinceEpoch,
		report.PayersMerkleRoot,
		report.NodeIds,
		domainSeparator,
	)
	if err != nil {
		return nil, err
	}

	if index > uint64(math.MaxUint32) {
		return nil, fmt.Errorf("report index %d exceeds max uint32", index)
	}

	index32 := uint32(index)

	return &payerreport.PayerReport{
		ID:                   *id,
		OriginatorNodeID:     nodeID,
		StartSequenceID:      report.StartSequenceId,
		EndSequenceID:        report.EndSequenceId,
		EndMinuteSinceEpoch:  report.EndMinuteSinceEpoch,
		PayersMerkleRoot:     report.PayersMerkleRoot,
		ActiveNodeIDs:        report.NodeIds,
		SubmittedReportIndex: &index32,
	}, nil
}

func prepareSignatures(
	signatures []payerreport.NodeSignature,
) []reportManager.IPayerReportManagerPayerReportSignature {
	// Copy and sort signatures by NodeID without mutating the input
	sortedSigs := make([]payerreport.NodeSignature, len(signatures))
	copy(sortedSigs, signatures)
	sort.Slice(sortedSigs, func(i, j int) bool {
		return sortedSigs[i].NodeID < sortedSigs[j].NodeID
	})

	out := make([]reportManager.IPayerReportManagerPayerReportSignature, len(sortedSigs))
	for i, sig := range sortedSigs {
		// Convert signature to use legacy EIP-155 recovery ID if needed
		sigBytes := sig.Signature
		if len(sigBytes) > 64 && sigBytes[64] < 27 {
			sigBytes = make([]byte, len(sig.Signature))
			copy(sigBytes, sig.Signature)
			sigBytes[64] += 27
		}

		out[i] = reportManager.IPayerReportManagerPayerReportSignature{
			NodeId:    sig.NodeID,
			Signature: sigBytes,
		}
	}

	return out
}

func prepareProof(proof *merkle.MultiProof) ([][]byte, [][32]byte, error) {
	// Coerce the types from merkle.Leaf to []byte
	rawLeaves := proof.GetLeaves()
	leaves := make([][]byte, len(rawLeaves))
	for idx, leaf := range rawLeaves {
		leaves[idx] = leaf
	}

	rawProofElements := proof.GetProofElements()
	proofElements := make([][32]byte, len(rawProofElements))
	for idx, element := range rawProofElements {
		elem32, err := utils.SliceToArray32(element)
		if err != nil {
			return nil, nil, err
		}
		proofElements[idx] = elem32
	}

	return leaves, proofElements, nil
}

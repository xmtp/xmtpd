package types

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	dm "github.com/xmtp/xmtpd/pkg/abi/distributionmanager"
	ft "github.com/xmtp/xmtpd/pkg/abi/feetoken"
	mft "github.com/xmtp/xmtpd/pkg/abi/mockunderlyingfeetoken"
	"github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/e2e/keys"
	"github.com/xmtp/xmtpd/pkg/fees"
	"go.uber.org/zap"
)

// chainClients holds lazily initialized blockchain clients for direct contract calls.
type chainClients struct {
	ethClient              *ethclient.Client
	adminSigner            blockchain.TransactionSigner
	contractsOpts          *config.ContractsOptions
	settlementAdmin        blockchain.ISettlementChainAdmin
	ratesAdmin             blockchain.IRatesAdmin
	payerRegistry          *payerregistry.PayerRegistry
	distributionManager    *dm.DistributionManager
	feeToken               *ft.FeeToken
	mockUnderlyingFeeToken *mft.MockUnderlyingFeeToken
}

// initChainClients lazily initializes the blockchain clients.
// Uses the external RPC URL (host-accessible) since these calls run from the test process.
// Unlike sync.Once, this retries on failure so transient errors don't permanently break
// all chain operations.
func (e *Environment) initChainClients(ctx context.Context) error {
	e.contractsMu.Lock()
	defer e.contractsMu.Unlock()

	if e.contracts != nil {
		return nil // already initialized
	}
	return e.doInitChainClients(ctx)
}

func (e *Environment) doInitChainClients(ctx context.Context) error {
	rpcURL := e.Chain.RPCURL()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to chain at %s: %w", rpcURL, err)
	}

	contractsOpts, err := config.LoadContractsConfig(config.ContractsSource{
		Environment: "anvil",
	})
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to load anvil contracts config: %w", err)
	}

	signer, err := blockchain.NewPrivateKeySigner(
		keys.AdminKey(),
		contractsOpts.SettlementChain.ChainID,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create admin signer: %w", err)
	}

	paramAdmin, err := blockchain.NewSettlementParameterAdmin(
		e.Logger.Named("param-admin"),
		client,
		signer,
		contractsOpts,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create parameter admin: %w", err)
	}

	settlementAdmin, err := blockchain.NewSettlementChainAdmin(
		e.Logger.Named("settlement-admin"),
		client,
		signer,
		contractsOpts,
		paramAdmin,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create settlement chain admin: %w", err)
	}

	ratesAdmin, err := blockchain.NewRatesAdmin(
		e.Logger.Named("rates-admin"),
		client,
		signer,
		paramAdmin,
		contractsOpts,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create rates admin: %w", err)
	}

	pr, err := payerregistry.NewPayerRegistry(
		common.HexToAddress(contractsOpts.SettlementChain.PayerRegistryAddress),
		client,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create payer registry binding: %w", err)
	}

	distMgr, err := dm.NewDistributionManager(
		common.HexToAddress(contractsOpts.SettlementChain.DistributionManagerAddress),
		client,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create distribution manager binding: %w", err)
	}

	feeToken, err := ft.NewFeeToken(
		common.HexToAddress(contractsOpts.SettlementChain.FeeToken),
		client,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create fee token binding: %w", err)
	}

	mockUnderlying, err := mft.NewMockUnderlyingFeeToken(
		common.HexToAddress(contractsOpts.SettlementChain.UnderlyingFeeToken),
		client,
	)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create mock underlying fee token binding: %w", err)
	}

	e.contracts = &chainClients{
		ethClient:              client,
		adminSigner:            signer,
		contractsOpts:          contractsOpts,
		settlementAdmin:        settlementAdmin,
		ratesAdmin:             ratesAdmin,
		payerRegistry:          pr,
		distributionManager:    distMgr,
		feeToken:               feeToken,
		mockUnderlyingFeeToken: mockUnderlying,
	}

	return nil
}

// signerForKey creates a TransactionSigner for the given private key hex string.
func (e *Environment) signerForKey(privateKey string) (blockchain.TransactionSigner, error) {
	return blockchain.NewPrivateKeySigner(
		privateKey,
		e.contracts.contractsOpts.SettlementChain.ChainID,
	)
}

// --- Rate operations ---

// RateOptions holds the parameters for updating rates.
type RateOptions struct {
	MessageFee    int64
	StorageFee    int64
	CongestionFee int64
	TargetRate    uint64
	StartTime     uint64 // 0 = default (2h from now)
}

// UpdateRates adds new rates to the RateRegistry.
func (e *Environment) UpdateRates(ctx context.Context, opts RateOptions) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	startTime := opts.StartTime
	if startTime == 0 {
		startTime = uint64(time.Now().Add(2 * time.Hour).Unix())
	}

	rates := fees.Rates{
		MessageFee:          currency.PicoDollar(opts.MessageFee),
		StorageFee:          currency.PicoDollar(opts.StorageFee),
		CongestionFee:       currency.PicoDollar(opts.CongestionFee),
		TargetRatePerMinute: opts.TargetRate,
		StartTime:           startTime,
	}

	if err := e.contracts.ratesAdmin.AddRates(ctx, rates); err != nil {
		return fmt.Errorf("add rates failed: %w", err)
	}

	e.Logger.Info("rates updated", zap.Any("rates", rates))
	return nil
}

// --- Settlement operations ---

// SendExcessToFeeDistributor moves excess funds from PayerRegistry to the
// DistributionManager.
func (e *Environment) SendExcessToFeeDistributor(ctx context.Context) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	if err := e.contracts.settlementAdmin.SendExcessToFeeDistributor(ctx); err != nil {
		return fmt.Errorf("send excess to fee distributor failed: %w", err)
	}

	e.Logger.Info("excess sent to fee distributor")
	return nil
}

// GetPayerRegistryExcess returns the current excess balance in the PayerRegistry.
func (e *Environment) GetPayerRegistryExcess(ctx context.Context) (*big.Int, error) {
	if err := e.initChainClients(ctx); err != nil {
		return nil, err
	}
	return e.contracts.settlementAdmin.GetPayerRegistryExcess(ctx)
}

// ClaimFromDistributionManager claims earned fees for a node from the
// DistributionManager. Must be called with the node owner's private key.
func (e *Environment) ClaimFromDistributionManager(
	ctx context.Context,
	nodeOwnerKey string,
	nodeID uint32,
	originatorNodeIDs []uint32,
	payerReportIndices []*big.Int,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	signer, err := e.signerForKey(nodeOwnerKey)
	if err != nil {
		return fmt.Errorf("failed to create node owner signer: %w", err)
	}

	txErr := blockchain.ExecuteTransaction(
		ctx,
		signer,
		e.Logger.Named("dm-claim"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.distributionManager.Claim(
				opts, nodeID, originatorNodeIDs, payerReportIndices,
			)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.distributionManager.ParseClaim(*log)
		},
		func(event any) {
			ev, ok := event.(*dm.DistributionManagerClaim)
			if !ok {
				return
			}
			e.Logger.Info("claimed from distribution manager",
				zap.Uint32("node_id", ev.NodeId),
				zap.String("amount", ev.Amount.String()))
		},
	)
	if txErr != nil {
		return fmt.Errorf("dm claim for node %d failed: %w", nodeID, txErr)
	}

	return nil
}

// WithdrawFromDistributionManager withdraws claimed fees for a node from the
// DistributionManager. Must be called with the node owner's private key.
// The recipient defaults to the node owner's address.
func (e *Environment) WithdrawFromDistributionManager(
	ctx context.Context,
	nodeOwnerKey string,
	nodeID uint32,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	signer, err := e.signerForKey(nodeOwnerKey)
	if err != nil {
		return fmt.Errorf("failed to create node owner signer: %w", err)
	}

	recipient := signer.FromAddress()

	txErr := blockchain.ExecuteTransaction(
		ctx,
		signer,
		e.Logger.Named("dm-withdraw"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.distributionManager.Withdraw(opts, nodeID, recipient)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.distributionManager.ParseWithdrawal(*log)
		},
		func(event any) {
			ev, ok := event.(*dm.DistributionManagerWithdrawal)
			if !ok {
				return
			}
			e.Logger.Info("withdrawn from distribution manager",
				zap.Uint32("node_id", ev.NodeId),
				zap.String("amount", ev.Amount.String()))
		},
	)
	if txErr != nil {
		return fmt.Errorf("dm withdraw for node %d failed: %w", nodeID, txErr)
	}

	return nil
}

// GetDistributionManagerOwedFees returns the owed fees for a given node.
func (e *Environment) GetDistributionManagerOwedFees(
	ctx context.Context,
	nodeID uint32,
) (*big.Int, error) {
	if err := e.initChainClients(ctx); err != nil {
		return nil, err
	}
	return e.contracts.settlementAdmin.GetDistributionManagerOwedFees(ctx, nodeID)
}

// --- Balance operations ---

// GetGasBalance returns the native ETH balance for the given address.
func (e *Environment) GetGasBalance(
	ctx context.Context,
	addr common.Address,
) (*big.Int, error) {
	if err := e.initChainClients(ctx); err != nil {
		return nil, err
	}
	return e.contracts.ethClient.BalanceAt(ctx, addr, nil)
}

// --- Payer operations ---

// DepositPayer deposits funds into the PayerRegistry for a given payer address.
// The deposit is made by the admin signer.
func (e *Environment) DepositPayer(
	ctx context.Context,
	payer common.Address,
	amount *big.Int,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	err := blockchain.ExecuteTransaction(
		ctx,
		e.contracts.adminSigner,
		e.Logger.Named("payer-deposit"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.payerRegistry.Deposit(opts, payer, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.payerRegistry.ParseDeposit(*log)
		},
		func(event any) {
			ev, ok := event.(*payerregistry.PayerRegistryDeposit)
			if !ok {
				return
			}
			e.Logger.Info("payer deposit",
				zap.String("payer", ev.Payer.Hex()),
				zap.String("amount", ev.Amount.String()))
		},
	)
	if err != nil {
		return fmt.Errorf("deposit for payer %s failed: %w", payer.Hex(), err)
	}

	return nil
}

// MintFeeToken mints mock underlying tokens to the admin, wraps them into fee
// tokens (xUSD), and returns the amount minted. This is only available on anvil.
func (e *Environment) MintFeeToken(ctx context.Context, amount *big.Int) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	admin := e.contracts.adminSigner.FromAddress()
	prAddr := common.HexToAddress(e.contracts.contractsOpts.SettlementChain.PayerRegistryAddress)
	feeTokenAddr := common.HexToAddress(e.contracts.contractsOpts.SettlementChain.FeeToken)

	// 1. Mint mock underlying to admin
	err := blockchain.ExecuteTransaction(
		ctx,
		e.contracts.adminSigner,
		e.Logger.Named("mint-underlying"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.mockUnderlyingFeeToken.Mint(opts, admin, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.mockUnderlyingFeeToken.ParseTransfer(*log)
		},
		func(event any) {},
	)
	if err != nil {
		return fmt.Errorf("mint underlying failed: %w", err)
	}

	// 2. Approve FeeToken to spend underlying
	err = blockchain.ExecuteTransaction(
		ctx,
		e.contracts.adminSigner,
		e.Logger.Named("approve-fee-token"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.mockUnderlyingFeeToken.Approve(opts, feeTokenAddr, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.mockUnderlyingFeeToken.ParseApproval(*log)
		},
		func(event any) {},
	)
	if err != nil {
		return fmt.Errorf("approve fee token failed: %w", err)
	}

	// 3. Wrap underlying into fee token (xUSD)
	err = blockchain.ExecuteTransaction(
		ctx,
		e.contracts.adminSigner,
		e.Logger.Named("wrap-fee-token"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.feeToken.Deposit(opts, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.feeToken.ParseTransfer(*log)
		},
		func(event any) {},
	)
	if err != nil {
		return fmt.Errorf("wrap fee token failed: %w", err)
	}

	// 4. Approve PayerRegistry to spend fee token
	err = blockchain.ExecuteTransaction(
		ctx,
		e.contracts.adminSigner,
		e.Logger.Named("approve-payer-registry"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.feeToken.Approve(opts, prAddr, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.feeToken.ParseApproval(*log)
		},
		func(event any) {},
	)
	if err != nil {
		return fmt.Errorf("approve payer registry failed: %w", err)
	}

	e.Logger.Info("fee tokens minted and approved",
		zap.String("amount", amount.String()),
		zap.String("admin", admin.Hex()))

	return nil
}

// FundPayer mints fee tokens and deposits them into the PayerRegistry for the
// given payer address. Handles the full flow: mint underlying → wrap → approve → deposit.
func (e *Environment) FundPayer(
	ctx context.Context,
	payer common.Address,
	amount *big.Int,
) error {
	if err := e.MintFeeToken(ctx, amount); err != nil {
		return fmt.Errorf("failed to mint fee tokens: %w", err)
	}
	if err := e.DepositPayer(ctx, payer, amount); err != nil {
		return fmt.Errorf("failed to deposit for payer: %w", err)
	}
	return nil
}

// GetPayerBalance returns the payer's balance in the PayerRegistry.
func (e *Environment) GetPayerBalance(
	ctx context.Context,
	payer common.Address,
) (*big.Int, error) {
	if err := e.initChainClients(ctx); err != nil {
		return nil, err
	}
	return e.contracts.payerRegistry.GetBalance(
		&bind.CallOpts{Context: ctx}, payer,
	)
}

// GetFeeTokenBalance returns the fee token (xUSD) balance for the given address.
func (e *Environment) GetFeeTokenBalance(
	ctx context.Context,
	addr common.Address,
) (*big.Int, error) {
	if err := e.initChainClients(ctx); err != nil {
		return nil, err
	}
	return e.contracts.feeToken.BalanceOf(
		&bind.CallOpts{Context: ctx}, addr,
	)
}

// RequestPayerWithdrawal requests a withdrawal from the PayerRegistry.
// Must be called with the payer's private key since only the payer can request withdrawals.
func (e *Environment) RequestPayerWithdrawal(
	ctx context.Context,
	payerPrivateKey string,
	amount *big.Int,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	signer, err := e.signerForKey(payerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create payer signer: %w", err)
	}

	txErr := blockchain.ExecuteTransaction(
		ctx,
		signer,
		e.Logger.Named("payer-request-withdrawal"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.payerRegistry.RequestWithdrawal(opts, amount)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.payerRegistry.ParseWithdrawalRequested(*log)
		},
		func(event any) {
			ev, ok := event.(*payerregistry.PayerRegistryWithdrawalRequested)
			if !ok {
				return
			}
			e.Logger.Info("withdrawal requested",
				zap.String("payer", ev.Payer.Hex()),
				zap.String("amount", ev.Amount.String()),
				zap.Uint32("withdrawable_timestamp", ev.WithdrawableTimestamp))
		},
	)
	if txErr != nil {
		return fmt.Errorf("request withdrawal failed: %w", txErr)
	}

	return nil
}

// CancelPayerWithdrawal cancels a pending withdrawal from the PayerRegistry.
// Must be called with the payer's private key.
func (e *Environment) CancelPayerWithdrawal(
	ctx context.Context,
	payerPrivateKey string,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	signer, err := e.signerForKey(payerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create payer signer: %w", err)
	}

	txErr := blockchain.ExecuteTransaction(
		ctx,
		signer,
		e.Logger.Named("payer-cancel-withdrawal"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.payerRegistry.CancelWithdrawal(opts)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.payerRegistry.ParseWithdrawalCancelled(*log)
		},
		func(event any) {
			ev, ok := event.(*payerregistry.PayerRegistryWithdrawalCancelled)
			if !ok {
				return
			}
			e.Logger.Info("withdrawal cancelled",
				zap.String("payer", ev.Payer.Hex()))
		},
	)
	if txErr != nil {
		return fmt.Errorf("cancel withdrawal failed: %w", txErr)
	}

	return nil
}

// FinalizePayerWithdrawal finalizes a pending withdrawal from the PayerRegistry,
// transferring funds to the given recipient address.
// Must be called with the payer's private key.
func (e *Environment) FinalizePayerWithdrawal(
	ctx context.Context,
	payerPrivateKey string,
	recipient common.Address,
) error {
	if err := e.initChainClients(ctx); err != nil {
		return err
	}

	signer, err := e.signerForKey(payerPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create payer signer: %w", err)
	}

	txErr := blockchain.ExecuteTransaction(
		ctx,
		signer,
		e.Logger.Named("payer-finalize-withdrawal"),
		e.contracts.ethClient,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return e.contracts.payerRegistry.FinalizeWithdrawal(opts, recipient)
		},
		func(log *types.Log) (any, error) {
			return e.contracts.payerRegistry.ParseWithdrawalFinalized(*log)
		},
		func(event any) {
			ev, ok := event.(*payerregistry.PayerRegistryWithdrawalFinalized)
			if !ok {
				return
			}
			e.Logger.Info("withdrawal finalized",
				zap.String("payer", ev.Payer.Hex()),
				zap.String("recipient", recipient.Hex()))
		},
	)
	if txErr != nil {
		return fmt.Errorf("finalize withdrawal failed: %w", txErr)
	}

	return nil
}

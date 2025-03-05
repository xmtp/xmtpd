package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/contracts/pkg/standardmerkletree"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/payerreports/merkle"
	"go.uber.org/zap"
)

type MerkleCaller struct {
	logger           *zap.Logger
	contract         *standardmerkletree.StandardMerkleTree
	contractsOptions config.ContractsOptions
}

func NewMerkleCaller(
	logger *zap.Logger,
	client *ethclient.Client,
	contractsOptions config.ContractsOptions,
) (*MerkleCaller, error) {
	contract, err := standardmerkletree.NewStandardMerkleTree(
		common.HexToAddress(contractsOptions.MerkleContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}
	return &MerkleCaller{
		logger:           logger,
		contract:         contract,
		contractsOptions: contractsOptions,
	}, nil
}

func (c *MerkleCaller) VerifyMultiProof(
	ctx context.Context,
	proof *merkle.MultiProof,
) (bool, error) {
	return c.contract.MultiProofVerify(
		&bind.CallOpts{Context: ctx},
		proof.Proof(),
		proof.ProofFlags(),
		proof.Offset(),
		proof.Root(),
		proof.PayerAddresses(),
		proof.Amounts(),
		proof.Indices(),
	)
}

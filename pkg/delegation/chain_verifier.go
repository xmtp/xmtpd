package delegation

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// PayerRegistryCaller is an interface matching the methods needed from the generated
// payerregistry.PayerRegistryCaller. This allows for easy testing and decouples
// the delegation package from the generated ABI bindings.
//
// NOTE: After deploying the updated PayerRegistry contract with delegation support,
// regenerate the ABI bindings using `dev/gen/abi` and this interface will be
// satisfied by the generated PayerRegistryCaller type.
type PayerRegistryCaller interface {
	IsAuthorized(opts *bind.CallOpts, payer common.Address, delegate common.Address) (bool, error)
	GetDelegation(opts *bind.CallOpts, payer common.Address, delegate common.Address) (struct {
		IsActive  bool
		Expiry    uint64
		CreatedAt uint64
	}, error)
}

// PayerRegistryChainVerifier implements ChainVerifier using the PayerRegistry contract.
type PayerRegistryChainVerifier struct {
	contract PayerRegistryCaller
}

// NewPayerRegistryChainVerifier creates a new chain verifier using the PayerRegistry contract.
func NewPayerRegistryChainVerifier(contract PayerRegistryCaller) *PayerRegistryChainVerifier {
	return &PayerRegistryChainVerifier{
		contract: contract,
	}
}

// IsAuthorized checks if a delegate is authorized to sign on behalf of a payer.
func (v *PayerRegistryChainVerifier) IsAuthorized(
	ctx context.Context,
	payer, delegate common.Address,
) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	authorized, err := v.contract.IsAuthorized(opts, payer, delegate)
	if err != nil {
		return false, fmt.Errorf("failed to check authorization: %w", err)
	}
	return authorized, nil
}

// GetDelegation returns the delegation info for a payer/delegate pair.
func (v *PayerRegistryChainVerifier) GetDelegation(
	ctx context.Context,
	payer, delegate common.Address,
) (*DelegationInfo, error) {
	opts := &bind.CallOpts{Context: ctx}
	delegation, err := v.contract.GetDelegation(opts, payer, delegate)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegation: %w", err)
	}
	return &DelegationInfo{
		IsActive:  delegation.IsActive,
		Expiry:    delegation.Expiry,
		CreatedAt: delegation.CreatedAt,
	}, nil
}

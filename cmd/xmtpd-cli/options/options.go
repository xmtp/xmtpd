// Package options implements the options for the CLI.
package options

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

/* ---------- Target ---------- */

type Target string

const (
	TargetIdentity               Target = "identity"
	TargetGroup                  Target = "group"
	TargetAppChainGateway        Target = "app-chain-gateway"
	TargetDistributionManager    Target = "distribution-manager"
	TargetPayerRegistry          Target = "payer-registry"
	TargetSettlementChainGateway Target = "settlement-chain-gateway"
)

func (t Target) String() string {
	if t == "" {
		return ""
	}
	return string(t)
}

func (t *Target) Set(s string) error {
	s = strings.TrimSpace(s)
	switch Target(s) {
	case TargetIdentity,
		TargetGroup,
		TargetAppChainGateway,
		TargetDistributionManager,
		TargetPayerRegistry,
		TargetSettlementChainGateway:
		*t = Target(s)
		return nil
	default:
		return fmt.Errorf(
			"invalid target %q (allowed: identity|group|app-chain-gateway|distribution-manager|payer-registry|settlement-chain-gateway)",
			s,
		)
	}
}

func (t Target) Type() string { return "target" }

/* ---------- PayloadBound ---------- */

type PayloadBound string

const (
	PayloadMin PayloadBound = "min"
	PayloadMax PayloadBound = "max"
)

func (b PayloadBound) String() string {
	if b == "" {
		return ""
	}
	return string(b)
}

func (b *PayloadBound) Set(s string) error {
	switch strings.TrimSpace(s) {
	case string(PayloadMin), string(PayloadMax):
		*b = PayloadBound(s)
		return nil
	default:
		return fmt.Errorf("invalid bound %q (allowed: min|max)", s)
	}
}

func (b PayloadBound) Type() string { return "payload-bound" }

/* ---------- AddressFlag ---------- */

type AddressFlag struct {
	Address common.Address
}

func (a AddressFlag) String() string {
	if (a.Address == common.Address{}) {
		return ""
	}
	return a.Address.Hex()
}

func (a *AddressFlag) Set(s string) error {
	s = strings.TrimSpace(s)
	if !common.IsHexAddress(s) {
		return fmt.Errorf("invalid Ethereum address: %q", s)
	}
	a.Address = common.HexToAddress(s)
	return nil
}

func (a AddressFlag) Type() string { return "address" }

// Common returns the go-ethereum common.Address.
func (a AddressFlag) Common() common.Address {
	return a.Address
}

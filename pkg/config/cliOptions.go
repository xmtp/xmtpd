package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/xmtp/xmtpd/pkg/currency"
)

type GlobalOptions struct {
	Contracts ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Log       LogOptions       `group:"Log Options"       namespace:"log"`
}

// NodeRegistryAdminOptions is the options for the node registry admin.
// It is intended to be used as a namespace inside a command option struct.
type NodeRegistryAdminOptions struct {
	AdminPrivateKey string `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
}

// RateRegistryAdminOptions is the options for the rate registry admin.
// It is intended to be used as a namespace inside a command option struct.
type RateRegistryAdminOptions struct {
	AdminPrivateKey string `long:"private-key" description:"Private key of the admin to administer the rates" required:"true"`
}

// NodeRegistryManagerOptions is the options for the node registry manager.
// It is intended to be used as a namespace inside a command option struct.
type NodeRegistryManagerOptions struct {
	NodePrivateKey string `long:"manager-private-key" description:"Private key of the node manager"`
	NodeId         int64  `long:"node-id"             description:"NodeId of the node to administer"`
}

/*
*
Command options
*
*/

type GenerateKeyOptions struct{}

type GetAllNodesOptions struct {
	OutFile string `long:"out-file" description:"File to write the nodes to"`
}

type GetNodeOptions struct {
	NodeId uint32 `long:"node-id" description:"NodeId of the node to get" required:"true"`
}

type MigrateNodesOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	InFile       string                   `                                        long:"in-file" description:"File to read the nodes from"`
}

type GetMaxCanonicalOptions struct{}

type SetMaxCanonicalOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Limit        uint8                    `                                        long:"limit" description:"Limit of max canonical nodes" required:"true"`
}

type GetBootstrapperAddressOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetBootstrapperAddressOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Address      AddressFlag              `                                        long:"address" description:"New bootstrapper address" required:"true"`
}

type GetPubKeyOptions struct {
	PrivateKey string `long:"private-key" description:"Private key you want the public key for" required:"true"`
}

type AddRatesOptions struct {
	AdminOptions  RateRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	MessageFee    currency.PicoDollar      `                                        long:"message-fee"    description:"Message fee"            required:"true"`
	StorageFee    currency.PicoDollar      `                                        long:"storage-fee"    description:"Storage fee"            required:"true"`
	CongestionFee currency.PicoDollar      `                                        long:"congestion-fee" description:"Congestion fee"         required:"true"`
	TargetRate    uint64                   `                                        long:"target-rate"    description:"Target rate per minute" required:"true"`
}

type GetRatesOptions struct{}

type RegisterNodeOptions struct {
	AdminOptions  NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	HttpAddress   string                   `                                        long:"http-address"         description:"HTTP address to register for the node"                            required:"true"`
	OwnerAddress  AddressFlag              `                                        long:"node-owner-address"   description:"Blockchain address of the intended owner of the registration NFT" required:"true"`
	SigningKeyPub string                   `                                        long:"node-signing-key-pub" description:"Signing key of the node to register"                              required:"true"`
	Force         bool                     `                                        long:"force"                description:"Register even if pubkey already exists"`
}

type NetworkAdminOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	NodeId       uint32                   `                                        long:"node-id" description:"NodeId to add to the network"`
}

type SetHttpAddressOptions struct {
	NodeManagerOptions NodeRegistryManagerOptions `group:"Node Manager Options" namespace:"node-manager"`
	Address            string                     `                                                      long:"address" description:"New HTTP address"`
	NodeId             uint32                     `                                                      long:"node-id" description:"NodeId to add to the network"`
}

type IdentityUpdatesStressOptions struct {
	PrivateKey string      `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
	Contract   AddressFlag `long:"contract"    description:"Contract address"                                required:"true"`
	Rpc        string      `long:"rpc"         description:"RPC URL"`
	Count      int         `long:"count"       description:"Number of transactions to send"`
	Async      bool        `long:"async"       description:"Send transactions asynchronously"`
}

type WatcherOptions struct {
	Contract AddressFlag `long:"contract" description:"Contract address" required:"true"`
	Wss      string      `long:"wss"      description:"WSS URL"          required:"true"`
}

type Target string

const (
	TargetIdentity               Target = "identity"
	TargetGroup                  Target = "group"
	TargetAppChainGateway        Target = "app-chain-gateway"
	TargetDistributionManager    Target = "distribution-manager"
	TargetPayerRegistry          Target = "payer-registry"
	TargetSettlementChainGateway Target = "settlement-chain-gateway"
)

func (t *Target) UnmarshalFlag(v string) error {
	switch v {
	case string(TargetIdentity),
		string(TargetGroup),
		string(TargetAppChainGateway),
		string(TargetDistributionManager),
		string(TargetPayerRegistry),
		string(TargetSettlementChainGateway):
		*t = Target(v)
		return nil
	default:
		return fmt.Errorf(
			"invalid target %q (allowed: identity|group|app-chain-gateway|distribution-manager|payer-registry|settlement-chain-gateway)",
			v,
		)
	}
}

type PausedFlag bool

const (
	PausedTrue  PausedFlag = true
	PausedFalse PausedFlag = false
)

func (p *PausedFlag) UnmarshalFlag(v string) error {
	switch v {
	case "true", "1", "yes", "on":
		*p = PausedTrue
		return nil
	case "false", "0", "no", "off":
		*p = PausedFalse
		return nil
	default:
		return fmt.Errorf("invalid paused value %q (allowed: true|false)", v)
	}
}

func (p *PausedFlag) Bool() bool {
	return bool(*p)
}

type SetPauseOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`

	// Which contract to pause/unpause.
	Target Target `long:"target" required:"true" description:"Target to pause: identity|group|app-chain-gateway|settlement-chain-gateway|payer-registry|distribution-manager"`

	// The desired paused state.
	Paused PausedFlag `long:"paused" required:"true" description:"Pause state to set (true|false)"`
}

type GetPauseOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`

	// Which contract’s pause status to read.
	Target Target `long:"target" required:"true" description:"Target to read: identity|group|app-chain-gateway|settlement-chain-gateway|payer-registry|distribution-manager"`
}

// For payload size keys (min / max) on the broadcasters.
type PayloadBound string

const (
	PayloadMin PayloadBound = "min"
	PayloadMax PayloadBound = "max"
)

func (b *PayloadBound) UnmarshalFlag(v string) error {
	switch v {
	case string(PayloadMin), string(PayloadMax):
		*b = PayloadBound(v)
		return nil
	default:
		return fmt.Errorf("invalid bound %q (allowed: min|max)", v)
	}
}

type GetPayloadSizeOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Target       Target                   `                                        long:"target" required:"true" description:"Which broadcaster: identity|group"`
	Bound        PayloadBound             `                                        long:"bound"  required:"true" description:"Which bound: min|max"`
}

type SetPayloadSizeOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Target       Target                   `                                        long:"target" required:"true" description:"Which broadcaster: identity|group"`
	Bound        PayloadBound             `                                        long:"bound"  required:"true" description:"Which bound: min|max"`
	Size         uint64                   `                                        long:"size"   required:"true" description:"Payload size in bytes"`
}

type GetNodeRegistryAdminOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type AddressFlag struct {
	Address common.Address
}

func (a *AddressFlag) UnmarshalFlag(value string) error {
	if !common.IsHexAddress(value) {
		return fmt.Errorf("invalid Ethereum address: %q", value)
	}
	a.Address = common.HexToAddress(value)
	return nil
}

func (a AddressFlag) String() string {
	return a.Address.Hex()
}

func (a AddressFlag) Common() common.Address {
	return a.Address
}

type SetNodeRegistryAdminOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Address      AddressFlag              `                                        long:"address" description:"New node registry admin address" required:"true"`
}

type GetDistributionManagerProtocolFeesRecipientOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetDistributionManagerProtocolFeesRecipientOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Address      AddressFlag              `                                        long:"address" description:"New DistributionManager protocolFeesRecipient address" required:"true"`
}

type GetPayerMinimumDepositOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetPayerMinimumDepositOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Amount       currency.MicroDollar     `                                        long:"amount" required:"true" description:"Minimum deposit in microdollars (1e-6 USD)"`
}

type GetPayerWithdrawLockPeriodOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetPayerWithdrawLockPeriodOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Seconds      uint32                   `                                        long:"seconds" required:"true" description:"Withdraw lock period in seconds"`
}

type GetPayerReportProtocolFeeRateOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetPayerReportProtocolFeeRateOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	FeeRateBps   uint16                   `                                        long:"bps" required:"true" description:"Protocol fee rate in basis points (uint16)"`
}

type GetRateRegistryMigratorOptions struct {
	AdminOptions RateRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
}

type SetRateRegistryMigratorOptions struct {
	AdminOptions RateRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	Address      AddressFlag              `                                        long:"address" description:"New RateRegistry migrator address" required:"true"`
}

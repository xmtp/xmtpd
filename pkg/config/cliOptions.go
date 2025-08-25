package config

import "github.com/xmtp/xmtpd/pkg/currency"

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
	Address      string                   `                                        long:"address" description:"New bootstrapper address"`
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
	OwnerAddress  string                   `                                        long:"node-owner-address"   description:"Blockchain address of the intended owner of the registration NFT" required:"true"`
	SigningKeyPub string                   `                                        long:"node-signing-key-pub" description:"Signing key of the node to register"                              required:"true"`
	Force         bool                     `                                        long:"force"                description:"Register even if pubkey already exists"                           required:"false"`
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
	PrivateKey string `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
	Contract   string `long:"contract"    description:"Contract address"`
	Rpc        string `long:"rpc"         description:"RPC URL"`
	Count      int    `long:"count"       description:"Number of transactions to send"`
	Async      bool   `long:"async"       description:"Send transactions asynchronously"`
}

type WatcherOptions struct {
	Contract string `long:"contract" description:"Contract address" required:"true"`
	Wss      string `long:"wss"      description:"WSS URL"          required:"true"`
}

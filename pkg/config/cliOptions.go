package config

type GlobalOptions struct {
	Contracts ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Log       LogOptions       `group:"Log Options"       namespace:"log"`
}

// NodeRegistryAdminOptions is the options for the node registry admin.
// It is intended to be used as a namespace inside a command option struct.
type NodeRegistryAdminOptions struct {
	AdminPrivateKey string `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
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
	NodeId int64 `long:"node-id" description:"NodeId of the node to get" required:"true"`
}

type MigrateNodesOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	InFile       string                   `                                        long:"in-file" description:"File to read the nodes from"`
}

type GetPubKeyOptions struct {
	PrivateKey string `long:"private-key" description:"Private key you want the public key for" required:"true"`
}

type AddRatesOptions struct {
	AdminPrivateKey string `long:"admin-private-key" description:"Private key of the admin to administer the node"`
	MessageFee      uint64 `long:"message-fee"       description:"Message fee"`
	StorageFee      uint64 `long:"storage-fee"       description:"Storage fee"`
	CongestionFee   uint64 `long:"congestion-fee"    description:"Congestion fee"`
	TargetRate      uint64 `long:"target-rate"       description:"Target rate per minute"`
	DelayDays       uint   `long:"delay-days"        description:"Delay the rates going into effect for N days"    default:"0"`
}

type RegisterNodeOptions struct {
	AdminOptions              NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	HttpAddress               string                   `                                        long:"http-address"                  description:"HTTP address to register for the node"                            required:"true"`
	OwnerAddress              string                   `                                        long:"node-owner-address"            description:"Blockchain address of the intended owner of the registration NFT" required:"true"`
	SigningKeyPub             string                   `                                        long:"node-signing-key-pub"          description:"Signing key of the node to register"                              required:"true"`
	MinMonthlyFeeMicroDollars int64                    `                                        long:"min-monthly-fee-micro-dollars" description:"Minimum monthly fee to register the node"                         required:"false"`
	Force                     bool                     `                                        long:"force"                         description:"Register even if pubkey already exists"                           required:"false"`
}

type NetworkAdminOptions struct {
	AdminOptions NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	NodeId       int64                    `                                        long:"node-id" description:"NodeId to add to the network"`
}

type SetHttpAddressOptions struct {
	NodeManagerOptions NodeRegistryManagerOptions `group:"Node Manager Options" namespace:"node-manager"`
	Address            string                     `                                                      long:"address" description:"New HTTP address"`
	NodeId             int64                      `                                                      long:"node-id" description:"NodeId to add to the network"`
}

type SetMinMonthlyFeeOptions struct {
	NodeManagerOptions        NodeRegistryManagerOptions `group:"Node Manager Options" namespace:"node-manager"`
	MinMonthlyFeeMicroDollars int64                      `                                                      long:"min-monthly-fee-micro-dollars" description:"Minimum monthly fee to register the node"`
	NodeId                    int64                      `                                                      long:"node-id"                       description:"NodeId to add to the network"`
}

type SetMaxActiveNodesOptions struct {
	AdminOptions   NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	MaxActiveNodes uint8                    `                                        long:"max-active-nodes" description:"Maximum number of active nodes"`
}

type SetNodeOperatorCommissionPercentOptions struct {
	AdminOptions      NodeRegistryAdminOptions `group:"Admin Options" namespace:"admin"`
	CommissionPercent int64                    `                                        long:"commission-percent" description:"Commission percent to set for the node operator"`
}

type IdentityUpdatesStressOptions struct {
	PrivateKey string `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
	Contract   string `long:"contract"    description:"Contract address"`
	Rpc        string `long:"rpc"         description:"RPC URL"`
	Count      int    `long:"count"       description:"Number of transactions to send"`
}

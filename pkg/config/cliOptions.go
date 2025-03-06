package config

type GlobalOptions struct {
	Contracts ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Log       LogOptions       `group:"Log Options"       namespace:"log"`
}

type AdminOptions struct {
	AdminPrivateKey string `long:"private-key" description:"Private key of the admin to administer the node" required:"true"`
	NodeId          int64  `long:"node-id"     description:"NodeId of the node to administer"`
}

type NodeManagerOptions struct {
	NodePrivateKey string `long:"manager-private-key" description:"Private key of the node manager"`
	NodeId         int64  `long:"node-id"             description:"NodeId of the node to administer"`
}

type NodeOperatorOptions struct {
	NodePrivateKey string `long:"node-private-key" description:"Private key of the node to administer"`
	NodeId         int64  `long:"node-id"          description:"NodeId of the node to administer"`
	Enable         bool   `long:"enable"           description:"Enable the node"`
}

type GenerateKeyOptions struct{}

type GetAllNodesOptions struct {
	OutFile string `long:"out-file" description:"File to write the nodes to"`
}

type GetOptions struct{}

type MigrateNodesOptions struct {
	AdminOptions AdminOptions `group:"Admin Options" namespace:"admin"`
	InFile       string       `                                        long:"in-file" description:"File to read the nodes from"`
}

type GetPubKeyOptions struct {
	PrivateKey string `long:"private-key" description:"Private key you want the public key for" required:"true"`
}

type RegisterNodeOptions struct {
	AdminOptions  AdminOptions `group:"Admin Options" namespace:"admin"`
	HttpAddress   string       `                                        long:"http-address"         description:"HTTP address to register for the node"                            required:"true"`
	OwnerAddress  string       `                                        long:"node-owner-address"   description:"Blockchain address of the intended owner of the registration NFT" required:"true"`
	SigningKeyPub string       `                                        long:"node-signing-key-pub" description:"Signing key of the node to register"                              required:"true"`
	MinMonthlyFee int64        `                                        long:"min-monthly-fee"      description:"Minimum monthly fee to register the node"                         required:"false"`
}

type SetHttpAddressOptions struct {
	NodeManagerOptions NodeManagerOptions `group:"Node Manager Options" namespace:"node-manager"`
	Address            string             `                                                      long:"address" description:"New HTTP address"`
}

type SetMinMonthlyFeeOptions struct {
	NodeManagerOptions NodeManagerOptions `group:"Node Manager Options" namespace:"node-manager"`
	MinMonthlyFee      int64              `                                                      long:"min-monthly-fee" description:"Minimum monthly fee to register the node"`
}

type SetMaxActiveNodesOptions struct {
	AdminOptions   AdminOptions `group:"Admin Options" namespace:"admin"`
	MaxActiveNodes uint8        `                                        long:"max-active-nodes" description:"Maximum number of active nodes"`
}

type SetNodeOperatorCommissionPercentOptions struct {
	AdminOptions      AdminOptions `group:"Admin Options" namespace:"admin"`
	CommissionPercent int64        `                                        long:"commission-percent" description:"Commission percent to set for the node operator"`
}

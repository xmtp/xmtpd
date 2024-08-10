package indexer

type ContractsOptions struct {
	RpcUrl                  string `log:"rpc-url" description:"Blockchain RPC URL"`
	NodesContractAddress    string `long:"nodes-address" description:"Node contract address"`
	MessagesContractAddress string `long:"messages-address" description:"Message contract address"`
}

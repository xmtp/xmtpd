package config

type ChainConfig struct {
	AppChainDeploymentBlock          int    `json:"appChainDeploymentBlock"`
	AppChainID                       int    `json:"appChainId"`
	DistributionManager              string `json:"distributionManager"`
	GroupMessageBroadcaster          string `json:"groupMessageBroadcaster"`
	IdentityUpdateBroadcaster        string `json:"identityUpdateBroadcaster"`
	NodeRegistry                     string `json:"nodeRegistry"`
	PayerRegistry                    string `json:"payerRegistry"`
	PayerReportManager               string `json:"payerReportManager"`
	RateRegistry                     string `json:"rateRegistry"`
	SettlementChainDeploymentBlock   int    `json:"settlementChainDeploymentBlock"`
	SettlementChainID                int    `json:"settlementChainId"`
	SettlementChainParameterRegistry string `json:"settlementChainParameterRegistry"`
}

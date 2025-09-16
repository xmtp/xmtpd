package config

type ChainConfig struct {
	AppChainDeploymentBlock          int    `json:"appChainDeploymentBlock"`
	AppChainID                       int    `json:"appChainId"`
	AppChainGateway                  string `json:"appChainGateway"`
	DistributionManager              string `json:"distributionManager"`
	GroupMessageBroadcaster          string `json:"groupMessageBroadcaster"`
	IdentityUpdateBroadcaster        string `json:"identityUpdateBroadcaster"`
	NodeRegistry                     string `json:"nodeRegistry"`
	PayerRegistry                    string `json:"payerRegistry"`
	PayerReportManager               string `json:"payerReportManager"`
	RateRegistry                     string `json:"rateRegistry"`
	SettlementChainDeploymentBlock   int    `json:"settlementChainDeploymentBlock"`
	SettlementChainGateway           string `json:"settlementChainGateway"`
	SettlementChainID                int    `json:"settlementChainId"`
	SettlementChainParameterRegistry string `json:"settlementChainParameterRegistry"`
	UnderlyingFeeToken               string `json:"underlyingFeeToken"`
}

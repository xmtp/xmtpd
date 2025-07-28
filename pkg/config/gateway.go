package config

type GatewayConfig struct {
	API       ApiOptions       `group:"API Options"            namespace:"api"`
	Contracts ContractsOptions `group:"Contracts Options"      namespace:"contracts"`
	Log       LogOptions       `group:"Log Options"            namespace:"log"`
	Metrics   MetricsOptions   `group:"Metrics Options"        namespace:"metrics"`
	Payer     PayerOptions     `group:"Payer Options"          namespace:"payer"`
	Redis     RedisOptions     `group:"Redis Options"          namespace:"redis"`
	Tracing   TracingOptions   `group:"DD APM Tracing Options" namespace:"tracing"`
}

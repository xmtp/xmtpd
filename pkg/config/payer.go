package config

type PayerConfig struct {
	API       ApiOptions       `group:"API Options"            namespace:"api"`
	DB        DbOptions        `group:"Database Options"       namespace:"db"`
	Contracts ContractsOptions `group:"Contracts Options"      namespace:"contracts"`
	Log       LogOptions       `group:"Log Options"            namespace:"log"`
	Metrics   MetricsOptions   `group:"Metrics Options"        namespace:"metrics"`
	Payer     PayerOptions     `group:"Payer Options"          namespace:"payer"`
	Tracing   TracingOptions   `group:"DD APM Tracing Options" namespace:"tracing"`
}

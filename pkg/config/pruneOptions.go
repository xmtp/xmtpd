package config

type PruneConfig struct {
	MaxCycles int   `long:"max-prune-cycles" env:"XMTPD_PRUNE_MAX_CYCLES" description:"Maximum pruning cycles" default:"10"`
	BatchSize int32 `long:"batch-size"       env:"XMTPD_PRUNE_BATCH_SIZE" description:"Batch size"             default:"10000"`
	DryRun    bool  `long:"dry-run"          env:"XMTPD_PRUNE_DRY_RUN"    description:"Dry run mode"`
}
type PruneOptions struct {
	DB          DbOptions        `group:"Database Options"  namespace:"db"`
	Log         LogOptions       `group:"Log Options"       namespace:"log"`
	Signer      SignerOptions    `group:"Signer Options"    namespace:"signer"`
	Contracts   ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	PruneConfig PruneConfig      `group:"Prune Options"     namespace:"prune"`
}

package config

type PruneOptions struct {
	DB  DbOptions  `group:"Database Options" namespace:"db"`
	Log LogOptions `group:"Log Options"      namespace:"log"`
}

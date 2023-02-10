package bolt

type Options struct {
	DataPath string `long:"data-path" env:"BOLT_DATA_PATH" description:"Path to the bolt DB file" default:""`
}

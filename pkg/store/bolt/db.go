package bolt

import bolt "go.etcd.io/bbolt"

func NewDB(opts *Options) (*bolt.DB, error) {
	return bolt.Open(opts.DataPath, 0600, nil)
}

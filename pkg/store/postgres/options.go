package postgresstore

type Options struct {
	DSN string `long:"dsn" env:"POSTGRES_DSN" description:"Postgres connection string" default:""`
}

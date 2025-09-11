package config

type RedisOptions struct {
	RedisURL  string `long:"redis-url"  env:"XMTPD_REDIS_URL"        description:"Redis URL"`
	KeyPrefix string `long:"key-prefix" env:"XMTPD_REDIS_KEY_PREFIX" description:"Redis key prefix" default:"xmtpd:"`
}

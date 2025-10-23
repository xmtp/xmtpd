package config

import "time"

type RedisOptions struct {
	RedisURL       string        `long:"redis-url"       env:"XMTPD_REDIS_URL"             description:"Redis URL"`
	KeyPrefix      string        `long:"key-prefix"      env:"XMTPD_REDIS_KEY_PREFIX"      description:"Redis key prefix"         default:"xmtpd:"`
	ConnectTimeout time.Duration `long:"connect-timeout" env:"XMTPD_REDIS_CONNECT_TIMEOUT" description:"Redis connection timeout" default:"10s"`
}

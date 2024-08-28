package config

import (
	"time"
)

type Config struct {
	HTTPServerHost              string
	HTTPServerPort              string
	HTTPServerReadTimeout       time.Duration
	HTTPServerReadHeaderTimeout time.Duration
	HTTPServerWriteTimeout      time.Duration
	HTTPServerIdleTimeout       time.Duration
	DBURI                       string
	Logging                     bool
}

package config

import (
	"time"

	"github.com/spf13/viper"
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

func New(path string) (Config, error) {
	viper.AutomaticEnv()
	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return Config{}, err
		}
	}
	return Config{
		HTTPServerHost:              viper.GetString("HTTP_SERVER_HOST"),
		HTTPServerPort:              viper.GetString("HTTP_SERVER_PORT"),
		HTTPServerReadTimeout:       viper.GetDuration("HTTP_SERVER_READ_TIMEOUT"),
		HTTPServerReadHeaderTimeout: viper.GetDuration("HTTP_SERVER_READ_HEADER_TIMEOUT"),
		HTTPServerWriteTimeout:      viper.GetDuration("HTTP_SERVER_WRITE_TIMEOUT"),
		HTTPServerIdleTimeout:       viper.GetDuration("HTTP_SERVER_IDLE_TIMEOUT"),
		DBURI:                       viper.GetString("DB_URI"),
		Logging:                     viper.GetBool("LOGGING"),
	}, nil
}

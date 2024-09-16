package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/bamboo-firewall/be/config"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

func main() {
	var pathConfig string
	flag.StringVar(&pathConfig, "config-file", "", "path to env config file")
	flag.Parse()

	cfg, err := loadConfig(pathConfig)
	if err != nil {
		slog.Warn("read config from file fail", "error", err)
	}
	newApp, err := NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err = newApp.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	interruptHandle(newApp)
}

func loadConfig(path string) (config.Config, error) {
	viper.AutomaticEnv()
	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return config.Config{}, err
		}
	}
	return config.Config{
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

func interruptHandle(app App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	slog.Debug("Listening Signal...")
	s := <-c
	slog.Info("Shutting down Server ...", "Received signal.", s)

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}

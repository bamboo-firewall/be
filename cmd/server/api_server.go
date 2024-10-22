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

	"github.com/bamboo-firewall/be/config"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

func main() {
	var pathConfig string
	flag.StringVar(&pathConfig, "config-file", "", "path to env config file")
	flag.Parse()

	cfg, err := config.New(pathConfig)
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

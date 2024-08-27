package main

import (
	"context"
	"github.com/bamboo-firewall/be/v1"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

func main() {
	app, err := v1.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		app.Start()
	}()
	interruptHandle(app)
}

func interruptHandle(app v1.App) {
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

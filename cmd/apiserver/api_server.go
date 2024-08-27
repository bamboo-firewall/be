package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bamboo-firewall/be/cmd/apiserver/rest"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

func main() {
	app, err := rest.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		app.Start()
	}()
	interruptHandle(app)
}

func interruptHandle(app rest.App) {
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

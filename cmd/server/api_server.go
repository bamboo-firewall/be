package main

import (
	"context"
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
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		app.Start()
	}()
	interruptHandle(app)
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

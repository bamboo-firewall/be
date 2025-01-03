package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/bamboo-firewall/be/buildinfo"
	"github.com/bamboo-firewall/be/config"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

func main() {
	var (
		pathConfig  string
		versionFlag bool
	)
	flag.StringVar(&pathConfig, "config-file", "", "path to env config file")
	flag.BoolVar(&versionFlag, "version", false, "show version information")
	flag.Parse()

	if versionFlag {
		version := fmt.Sprintf("Version: \t%s\nBranch: \t%s\nBuild: \t%s\nOrganzition: \t%s", buildinfo.Version, buildinfo.GitBranch, buildinfo.BuildDate, buildinfo.Organization)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, version)
		w.Flush()
		os.Exit(0)
	}

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

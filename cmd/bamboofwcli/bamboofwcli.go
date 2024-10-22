package main

import (
	"log/slog"
	"os"

	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	command.Execute()
}

// Package main is the entry point for the demo application.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ComixHe/go-template/internal/metadata"
)

func main() {
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(metadata.String())

		return
	}

	if err := run(context.Background()); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

// run starts the application and blocks until the context is canceled by
// SIGINT or SIGTERM. Cleanup registered with defer runs before returning.
func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting",
		"app", metadata.AppName,
		"version", metadata.Version,
		"commit", metadata.Commit,
	)

	<-ctx.Done()

	slog.Info("shutting down", "reason", context.Cause(ctx))

	return nil
}

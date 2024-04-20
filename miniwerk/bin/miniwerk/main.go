package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RobinThrift/miniwerk/miniwerk/app"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("fatal error running miniwerk:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := app.NewConfigFromEnv()
	if err != nil {
		return fmt.Errorf("error parsing config from env: %w", err)
	}

	app, err := app.New(ctx, config)
	if err != nil {
		return fmt.Errorf("error creating miniwerk app: %w", err)
	}

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		sig := <-shutdown

		stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCtxCancel()

		slog.InfoContext(ctx, fmt.Sprintf("received signal %v: triggering shutdown", sig))

		err := app.Stop(stopCtx) //nolint: contextcheck // false positive
		if err != nil {
			slog.ErrorContext(ctx, "could not stop gracefully", "error", err)
		}
	}()

	return app.Start(ctx)
}

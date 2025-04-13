package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lorenzogood/x/quotestack/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("failed to start quotestack", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	app.Run(ctx)
	<-ctx.Done()

	return nil
}

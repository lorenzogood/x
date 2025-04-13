package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Serve(ctx context.Context, addr string, handler http.Handler) {
	logger := slog.With("address", addr)

	s := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second, //nolint:mnd // fine
		ReadHeaderTimeout: 2 * time.Second,  //nolint:mnd // fine
		Handler:           handler,
		Addr:              addr,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Shutdown(sCtx); err != nil {
			logger.Warn("failed while stopping http server", "error", err)
		}

		logger.Debug("stopped http server")
	}()

	logger.Debug("started http server thread")
}

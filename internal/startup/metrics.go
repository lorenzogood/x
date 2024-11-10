package startup

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var metricsPort = flag.Int("metrics-port", 2112, "where to expose prometheus metrics.")

// Expose all registered prometheus metrics over metricsPort.
func Metrics(ctx context.Context) {
	addr := fmt.Sprintf(":%d", *metricsPort)
	srv := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second, //nolint:mnd // fine
		ReadHeaderTimeout: 2 * time.Second,  //nolint:mnd // fine
		Handler:           promhttp.Handler(),
		Addr:              addr,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("failed to start metrics http server", zap.Error(err))
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutDownCtx); err != nil {
			zap.L().Error("failed to stop metrics http server", zap.Error(err))
		}
	}()

	zap.L().Debug("started metrics http server", zap.Int("port", *metricsPort))
}

package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lorenzogood/x/quotestack/app/website"
	"github.com/lorenzogood/x/quotestack/internal/settings"
	"github.com/lorenzogood/x/quotestack/internal/web"
	"github.com/lorenzogood/x/quotestack/internal/web/mid"
)

func Run(ctx context.Context) error {
	cfg := settings.Get()

	webApp, err := website.New()
	if err != nil {
		return fmt.Errorf("error starting webapp: %w", err)
	}

	r := chi.NewMux()
	r.Use(mid.LogRecover)
	r.Use(middleware.Compress(5))
	r.Method(http.MethodGet, "/", web.Handler(webApp.Index))

	addr := fmt.Sprintf(fmt.Sprintf("0.0.0.0:%d", cfg.Port))

	web.Serve(ctx, addr, r)
	slog.Info("server started", "address", addr)

	return nil
}

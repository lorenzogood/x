package mid

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LogRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic in request handler",
					"status", ww.Status(),
					"duration", time.Since(start),
					"path", r.URL.Path,
					"panic", err,
					"ua", r.UserAgent())
			} else {
				mills := time.Since(start) / time.Millisecond
				logFn := slog.Debug

				if mills >= 500 {
					logFn = slog.Info
				}
				if ww.Status() >= 500 {
					logFn = slog.Error
				}

				logFn(
					"http request",
					"status", ww.Status(),
					"duration", time.Since(start),
					"path", r.URL.Path,
					"ua", r.UserAgent(),
				)
			}
		}()

		next.ServeHTTP(ww, r)
	})
}

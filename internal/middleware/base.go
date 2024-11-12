package middleware

import (
	"time"

	"github.com/lorenzogood/x/internal/router"
	"go.uber.org/zap"
)

func LogRecover(h router.Handler) router.Handler {
	return func(c *router.Ctx) error {
		start := time.Now().UTC()

		defer func() {
			if err := recover(); err != nil {
				zap.L().Error(
					"panic in http request handler",
					zap.Int("status", int(c.StatusCode())),
					zap.Duration("duration", time.Since(start)*time.Millisecond),
					zap.String("path", c.Request().Pattern),
					zap.Any("panic", err),
					zap.String("user_agent", c.Request().Header.Get("User-Agent")),
				)
			} else {
				mills := time.Since(start) / time.Millisecond
				logFn := zap.L().Debug

				if mills >= 500 {
					logFn = zap.L().Info
				}
				if c.StatusCode() >= 500 {
					logFn = zap.L().Error
				}

				logFn(
					"http request",
					zap.Int("status", int(c.StatusCode())),
					zap.Duration("duration", time.Since(start)*time.Millisecond),
					zap.String("method", string(c.Method())),
					zap.String("path", c.Request().Pattern),
					zap.String("user_agent", c.Request().Header.Get("User-Agent")),
				)
			}
		}()

		return h(c)
	}
}

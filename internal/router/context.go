package router

import (
	"net/http"

	"go.uber.org/zap"
)

type Ctx struct {
	r          *http.Request
	w          http.ResponseWriter
	statusCode HttpStatusCode
}

func (c *Ctx) Request() *http.Request {
	return c.r
}

func (c *Ctx) Response() http.ResponseWriter {
	return c.w
}

func (c *Ctx) Log() *zap.Logger {
	logger := zap.L().Named("web")

	return logger
}

func (c *Ctx) SetStatus(s HttpStatusCode) {
	c.statusCode = s

	c.Response().WriteHeader(int(s))
}

func (c *Ctx) StatusCode() HttpStatusCode {
	return c.statusCode
}

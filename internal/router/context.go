package router

import (
	"encoding/json"
	"net/http"

	"github.com/lorenzogood/x/internal/templates"
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

func (c *Ctx) Header() http.Header {
	return c.w.Header()
}

func (c *Ctx) StatusCode() HttpStatusCode {
	return c.statusCode
}

func (c *Ctx) Respond(status HttpStatusCode, b []byte) error {
	c.SetStatus(status)
	_, err := c.Response().Write(b)
	return err
}

func (c *Ctx) RespondString(status HttpStatusCode, b string) error {
	return c.Respond(status, []byte(b))
}

func (c *Ctx) RespondTemplate(t *templates.TemplateRenderer, status HttpStatusCode, name string, data any) error {
	c.SetStatus(status)
	c.Header().Set("Content-Type", "text/html")

	return t.Render(name, data, c.w)
}

func (c *Ctx) RespondJson(status HttpStatusCode, data any) error {
	c.SetStatus(status)
	c.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(c.w).Encode(data); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) Method() RouteMethod {
	return RouteMethod(c.r.Method)
}

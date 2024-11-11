package router

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type Handler func(c *Ctx) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Ctx{
		r: r,
		w: w,
	}

	if err := h(ctx); err != nil {
		ctx.Log().Error("handler error", zap.Error(err))
	}
}

type Middleware func(Handler) Handler

// Route, either a group or a leaf node.
// Routes get added onto a serve mux, at the provied path.
type Route interface {
	Register(m *http.ServeMux)
}

type Router struct {
	// Base route the router is mounted at, all routes are at {route}/{subroute}.
	route  string
	routes []Route

	// Only here to be copied onto leaf nodes.
	middleware []Middleware
}

func New() *Router {
	return new(Router)
}

func (r *Router) Handle(m RouteMethod, route string, h Handler, middleware ...Middleware) {
	mid := make([]Middleware, 0)
	mid = append(mid, r.middleware...)
	mid = append(mid, middleware...)

	route = strings.TrimPrefix(route, "/")

	l := &Leaf{
		route:      fmt.Sprintf("%s/%s", r.route, route),
		handler:    h,
		method:     m,
		middleware: mid,
	}

	r.routes = append(r.routes, l)
}

func (r *Router) Group(route string, fn func(r *Router)) {
	mid := make([]Middleware, 0)
	mid = append(mid, r.middleware...)

	ro := &Router{
		route:      fmt.Sprintf("%s/%s", r.route, route),
		middleware: mid,
	}

	fn(ro)

	r.routes = append(r.routes, ro)
}

func (r *Router) Use(middleware ...Middleware) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) Register(m *http.ServeMux) {
	for _, route := range r.routes {
		route.Register(m)
	}
}

type Leaf struct {
	method     RouteMethod
	route      string
	handler    Handler
	middleware []Middleware
}

func (l *Leaf) Register(m *http.ServeMux) {
	h := l.handler

	for _, middleware := range l.middleware {
		h = middleware(h)
	}

	m.Handle(fmt.Sprintf("%s %s", l.method, l.route), h)
}

package website

import (
	"fmt"

	"github.com/lorenzogood/x/quotestack"
	"github.com/lorenzogood/x/quotestack/internal/templates"
	"github.com/lorenzogood/x/quotestack/internal/web"
)

type App struct {
	templates *templates.Templates
}

func New() (*App, error) {
	templ, err := templates.New(quotestack.Templates, "templates")
	if err != nil {
		return nil, fmt.Errorf("error building templates: %w", err)
	}

	return &App{templates: templ}, nil
}

func (a *App) Render(ctx *web.Ctx, status web.HttpStatusCode, template string, data any) error {
	render_data := struct {
		Data any
	}{
		Data: data,
	}

	ctx.Header().Set("Content-Type", "text/html")
	ctx.SetStatus(status)

	return a.templates.Run(ctx.Response(), template, render_data)
}

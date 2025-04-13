package website

import "github.com/lorenzogood/x/quotestack/internal/web"

func (a *App) Index(ctx *web.Ctx) error {
	return a.Render(ctx, web.OK, "base.tmpl.html", nil)
}

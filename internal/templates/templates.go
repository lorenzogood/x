package templates

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var templateRenderTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "x_template_render_time",
	Help: "Template render times.",
}, []string{"template", "is_error"})

type TemplateRenderer struct {
	templ *template.Template
}

func New(base string, funcMap template.FuncMap) (*TemplateRenderer, error) {
	base = filepath.Clean(base)

	logger := zap.L().Named("template build").With(zap.String("dir", base))

	t := template.New("")
	_ = t.Funcs(funcMap)

	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk error: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		name := strings.TrimPrefix(path, base+"/")

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		if t, err = template.New(name).Parse(string(content)); err != nil {
			return fmt.Errorf("template render error for %s: %w", path, err)
		}

		logger.Debug("added template", zap.String("name", name))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		templ: t,
	}, nil
}

func (t *TemplateRenderer) Render(name string, data any, w io.Writer) error {
	start := time.Now()

	err := t.templ.ExecuteTemplate(w, name, data)
	if err != nil {
		templateRenderTime.With(prometheus.Labels{"template": name, "is_error": "true"}).Observe(float64(time.Since(start)))
		return err
	}

	templateRenderTime.With(prometheus.Labels{"template": name, "is_error": "false"}).Observe(float64(time.Since(start)))

	return nil
}

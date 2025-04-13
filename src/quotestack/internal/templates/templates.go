package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"
)

type Templates struct {
	t *template.Template
}

func New(f fs.FS, root string) (*Templates, error) {
	base := template.New("")
	err := fs.WalkDir(f, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(f, path)
		if err != nil {
			return fmt.Errorf("error reading template file: %w", err)
		}

		fname := strings.TrimPrefix(path, root+"/")

		if _, err := base.New(fname).Parse(string(content)); err != nil {
			return fmt.Errorf("error parsing tempalte file: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Templates{t: base}, nil
}

func (t *Templates) List() []string {
	templates := t.t.Templates()
	names := make([]string, len(templates))
	for i, v := range templates {
		names[i] = v.Name()
	}

	return names
}

func (t *Templates) Run(w io.Writer, name string, data any) error {
	return t.t.ExecuteTemplate(w, name, data)
}

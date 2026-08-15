package web

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"path/filepath"
)

type Templates struct {
	pages map[string]*template.Template
}

func NewTemplates() (*Templates, error) {
	t := &Templates{
		pages: make(map[string]*template.Template),
	}

	pages, err := filepath.Glob("web/templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"web/templates/layouts/base.html",
			page,
		}

		components, err := filepath.Glob(
			"web/templates/components/*.html",
		)
		if err != nil {
			return nil, err
		}

		files = append(files, components...)

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		t.pages[name] = tmpl
	}

	return t, nil
}

func (t *Templates) RenderBase(
	w io.Writer,
	page string,
	data any,
) error {
	tmpl, exists := t.pages[page]
	if !exists {
		slog.Error("template not found", "page", page)
		return fmt.Errorf("template not found: %s", page)
	}

	return tmpl.ExecuteTemplate(
		w,
		"base",
		data,
	)
}

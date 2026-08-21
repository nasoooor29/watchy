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

func newTemplates() *Templates {
	t := &Templates{
		pages: make(map[string]*template.Template),
	}

	pages, err := filepath.Glob("web/templates/pages/*.html")
	if err != nil {
		panic(err)
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
			panic(err)
		}

		files = append(files, components...)

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			panic(err)
		}

		t.pages[name] = tmpl
	}

	return t
}

func RenderBase(
	w io.Writer,
	page string,
	data any,
) error {
	// NOTE: change it in prod to create the template once and reuse it
	t := newTemplates()
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

func RenderComponent(
	w io.Writer,
	page string,
	component string,
	data any,
) error {
	// NOTE: change it in prod to create the template once and reuse it
	t := newTemplates()
	tmpl, exists := t.pages[page]
	if !exists {
		slog.Error("template not found", "page", page)
		return fmt.Errorf("template not found: %s", page)
	}

	return tmpl.ExecuteTemplate(w, component, data)
}

func RenderContent(
	w io.Writer,
	page string,
	data any,
) error {
	// NOTE: change it in prod to create the template once and reuse it
	t := newTemplates()
	tmpl, exists := t.pages[page]
	if !exists {
		slog.Error("template not found", "page", page)
		return fmt.Errorf("template not found: %s", page)
	}

	return tmpl.ExecuteTemplate(w, "content", data)
}

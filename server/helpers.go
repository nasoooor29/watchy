package server

import (
	"backend/web"
	"log/slog"
	"net/http"
)

func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, page string, data ...any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		if len(data) > 0 {
			if err := web.RenderContent(w, page, data[0]); err != nil {
				slog.Error("rendering page fragment failed", "page", page, "err", err)
			}
		} else {
			if err := web.RenderContent(w, page, nil); err != nil {
				slog.Error("rendering page fragment failed", "page", page, "err", err)
			}
		}
		return
	}
	w.WriteHeader(status)
	var pageData any
	if len(data) > 0 {
		pageData = data[0]
	}
	if err := web.RenderBase(w, page, pageData); err != nil {
		slog.Error("rendering page failed", "page", page, "err", err)
	}
}
func (s *Server) renderComponent(
	w http.ResponseWriter,
	template string,
	component string,
	data any,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := web.RenderComponent(w, template, component, data); err != nil {
		slog.Error("failed to render component",
			"template", template,
			"component", component,
			"err", err,
		)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
func (s *Server) internalError(
	w http.ResponseWriter,
	msg string,
	err error,
	args ...any,
) {
	args = append(args, "err", err)
	slog.Error(msg, args...)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
func isHtmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request, location string) {
	if isHtmx(r) {
		w.Header().Set("HX-Redirect", location)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func (s *Server) RenderPageHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, http.StatusOK, page)
	}
}

package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"backend/web"
)

func (s *Server) GetShowPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	detail, err := web.BuildShowDetail(s.DB, id)
	if err != nil {
		slog.Error("failed to build show detail", "id", id, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.render(w, r, http.StatusOK, "show.html", detail)
}

func (s *Server) GetPoster(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	show, err := s.DB.GetShow(id)
	if err != nil {
		slog.Error("failed to get show", "id", id, "err", err)
		http.NotFound(w, r)
		return
	}

	if len(show.Poster) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(show.Poster)
}

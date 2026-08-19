package server

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) ShowsPage(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(s.env.TvDir)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	slog.Info("shows page", "dirs", len(dirs))
	shows := []db.Show{}
	withMD := 0
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		show, err := s.Indexer.IndexShow(filepath.Join(s.env.TvDir, dir.Name()))
		if err != nil {
			slog.Error("Failed to get metadata for show", "show", dir.Name(), "err", err)
			continue
		}
		shows = append(shows, *show)
	}
	// print the number of the shows that have metadata.json file
	slog.Info("shows page", "shows with metadata", withMD)
	s.render(w, r, http.StatusOK, "home.html", map[string]any{
		"TotalSeries": 5,
		"Shows":       shows,
	})
}

func (s *Server) GetShowCover(w http.ResponseWriter, r *http.Request) {
	coverPath := filepath.Join(s.env.TvDir, r.PathValue("ShowName"), "cover.jpg")

	if _, err := os.Stat(coverPath); err != nil {
		slog.Error("Failed to get cover for show", "show", r.PathValue("ShowName"), "err", err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// set the content type to image/jpeg
	w.Header().Set("Content-Type", "image/jpeg")

	// serve the cover.jpg file
	http.ServeFile(w, r, coverPath)
}

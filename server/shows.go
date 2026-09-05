package server

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"backend/db"
	"backend/models"
	"backend/utils"
)

func buildLibraryPage(tvDir string) ([]models.Show, error) {
	entries, err := os.ReadDir(tvDir)
	if err != nil {
		return nil, err
	}

	libShows := []models.Show{}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "" || entry.Name()[0] == '.' {
			continue
		}
		showPath := filepath.Join(tvDir, entry.Name())
		// if the path have .hide-from-list file, if so, skip it
		if utils.IsHiddenShow(showPath) {
			slog.Debug("skipping hidden show", "show", entry.Name(), "path", showPath)
			continue
		}
		libShows = append(libShows, models.Show{
			Path:     showPath,
			Metadata: models.Metadata{Title: entry.Name(), Image: "/public/fallback.svg"},
		})
	}

	return libShows, nil
}

func buildShowDetail(database *db.SQL, showPath string) (models.Show, error) {
	if utils.IsHiddenShow(showPath) {
		return models.Show{}, models.ErrHidden
	}

	detail := models.Show{
		Path:     showPath,
		Metadata: models.Metadata{Title: filepath.Base(showPath), Image: "/public/fallback.svg"},
	}
	if show, err := database.GetShowByPath(showPath); err == nil {
		detail.Metadata.Title = show.Title
		detail.Metadata.Image = "/api/poster/" + showPath
	}

	fseasons, err := utils.GetShowSeasons(showPath)
	if err != nil {
		return models.Show{}, err
	}
	if len(fseasons) == 0 {
		return models.Show{}, models.ErrFailedToBuildShowDetails
	}

	detail.Seasons = make(map[string][]models.Episode)
	slices.Reverse(fseasons)
	for i, name := range fseasons {
		eps, err := utils.GetShowEpisodes(fseasons[i], showPath)
		if err != nil {
			return models.Show{}, err
		}
		slices.Reverse(eps)
		detail.Seasons[name] = eps
	}

	return detail, nil
}

func (s *Server) GetShow(w http.ResponseWriter, r *http.Request) {
	id, err := url.PathUnescape(r.PathValue("id"))
	if err != nil || id == "" {
		http.NotFound(w, r)
		return
	}

	detail, err := buildShowDetail(s.DB, id)
	if err != nil {
		s.internalError(w, "failed to build show detail", err, "id", id)
		return
	}

	if isHtmx(r) {
		s.renderComponent(w, "home.html", "show-panel", detail)
		return
	}

	page, err := buildLibraryPage(s.env.TvDir)
	if err != nil {
		s.internalError(w, "failed to build library page", err)
		return
	}

	s.render(w, r, http.StatusOK, "home.html", map[string]any{
		"Shows": page,
	})
}

func (s *Server) GetPoster(w http.ResponseWriter, r *http.Request) {
	showPath, err := url.PathUnescape(r.PathValue("id"))
	if err != nil || showPath == "" {
		http.NotFound(w, r)
		return
	}

	show, err := s.DB.GetShowByPath(showPath)
	if err != nil {
		slog.Error("failed to get show", "path", showPath, "err", err)
		http.NotFound(w, r)
		return
	}

	if len(show.Poster) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	// cache the poster for a year and immutable
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(show.Poster)
}

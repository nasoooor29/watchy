package server

import (
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"backend/db"
	"backend/models"
	"backend/utils"
)

func buildLibraryPage(database *db.SQL, tvDir string) ([]models.Show, error) {
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
		relativeShowPath, err := filepath.Rel(tvDir, showPath)
		if err != nil {
			return nil, err
		}
		relativeShowPath = filepath.ToSlash(relativeShowPath)
		// if the path have .hide-from-list file, if so, skip it
		if utils.IsHiddenShow(showPath) {
			slog.Debug("skipping hidden show", "show", entry.Name(), "path", showPath)
			continue
		}
		show := models.Show{
			Path:     relativeShowPath,
			Metadata: models.Metadata{Title: entry.Name(), Image: "/public/fallback.svg"},
		}
		if record, err := database.GetShowByPath(relativeShowPath); err == nil {
			show.Metadata = record.GetMetadata()
			if len(record.Poster) > 0 {
				show.Metadata.Image = "/api/poster/" + relativeShowPath
			}
		}
		// show.Metadata.LatestEp = utils.GetLatestFileTimestamp(showPath)
		show.Metadata.LatestEp = rand.Int63n(1000000)

		libShows = append(libShows, show)
	}

	return libShows, nil
}

func buildShowDetail(database *db.SQL, tvDir, showPath string) (models.Show, error) {
	showPath = filepath.FromSlash(showPath)
	if !filepath.IsAbs(showPath) {
		showPath = filepath.Join(tvDir, showPath)
	}
	showPath, _ = filepath.Abs(showPath)
	if utils.IsHiddenShow(showPath) {
		return models.Show{}, models.ErrHidden
	}

	relativeShowPath, _ := filepath.Rel(tvDir, showPath)
	relativeShowPath = filepath.ToSlash(relativeShowPath)
	detail := models.Show{
		Path:     relativeShowPath,
		Metadata: models.Metadata{Title: filepath.Base(showPath), Image: "/public/fallback.svg"},
	}
	if show, err := database.GetShowByPath(filepath.ToSlash(relativeShowPath)); err == nil {
		detail.Metadata = show.GetMetadata()
		if len(show.Poster) > 0 {
			detail.Metadata.Image = "/api/poster/" + filepath.ToSlash(relativeShowPath)
		}
	}

	fseasons, err := utils.GetShowSeasons(showPath)
	if err != nil {
		return models.Show{}, err
	}
	if len(fseasons) == 0 {
		return models.Show{}, models.ErrFailedToBuildShowDetails
	}

	slices.Reverse(fseasons)
	detail.Seasons = make([]models.Season, 0, len(fseasons))
	for i, name := range fseasons {
		eps, err := utils.GetShowEpisodes(fseasons[i], showPath)
		if err != nil {
			return models.Show{}, err
		}
		slices.Reverse(eps)
		for episodeIndex := range eps {
			for pathIndex, episodePath := range eps[episodeIndex].Paths {
				relativePath, err := filepath.Rel(tvDir, episodePath)
				if err != nil {
					return models.Show{}, err
				}
				eps[episodeIndex].Paths[pathIndex] = filepath.ToSlash(relativePath)
			}
		}
		detail.Seasons = append(detail.Seasons, models.Season{Name: name, Episodes: eps})
	}

	return detail, nil
}

func (s *Server) GetShow(w http.ResponseWriter, r *http.Request) {
	id, err := url.PathUnescape(r.PathValue("id"))
	if err != nil || id == "" {
		http.NotFound(w, r)
		return
	}

	detail, err := buildShowDetail(s.DB, s.env.TvDir, id)
	if err != nil {
		s.internalError(w, "failed to build show detail", err, "id", id)
		return
	}

	if isHtmx(r) {
		s.renderComponent(w, "home.html", "show-panel", detail)
		return
	}

	page, err := buildLibraryPage(s.DB, s.env.TvDir)
	if err != nil {
		s.internalError(w, "failed to build library page", err)
		return
	}

	s.render(w, r, http.StatusOK, "home.html", map[string]any{
		"Shows": page,
		"Show":  detail,
	})
}

func (s *Server) GetPoster(w http.ResponseWriter, r *http.Request) {
	showPath, err := url.PathUnescape(r.PathValue("id"))
	if err != nil || showPath == "" {
		http.NotFound(w, r)
		return
	}
	showPath = filepath.ToSlash(showPath)

	relativeShowPath := filepath.Clean(showPath)
	show, err := s.DB.GetShowByPath(filepath.ToSlash(relativeShowPath))
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

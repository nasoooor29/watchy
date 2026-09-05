package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"backend/db"
	"backend/models"
	"backend/utils"
)

func buildLibraryPage(database *db.SQL) ([]models.Show, error) {
	dbShows, err := database.GetAllShows()
	if err != nil {
		return nil, err
	}

	libShows := []models.Show{}
	for _, s := range dbShows {
		// if the path have .hide-from-list file, if so, skip it
		if utils.IsHiddenShow(s.Path) {
			slog.Debug("skipping hidden show", "show", s.Title, "path", s.Path)
			continue
		}
		libShows = append(libShows, models.Show{
			Path:     s.Path,
			Metadata: models.Metadata{Title: s.Title, Image: "/api/poster/" + fmt.Sprint(s.ID)},
		})
	}

	return libShows, nil
}

func buildShowDetail(database *db.SQL, id int64) (models.Show, error) {
	show, err := database.GetShow(id)
	if err != nil {
		return models.Show{}, err
	}
	if utils.IsHiddenShow(show.Path) {
		return models.Show{}, models.ErrHidden
	}

	detail := models.Show{
		Path:     show.Path,
		Metadata: models.Metadata{Title: show.Title, Image: "/api/poster/" + fmt.Sprint(show.ID)},
	}

	fseasons, err := utils.GetShowSeasons(show.Path)
	if err != nil {
		return models.Show{}, err
	}
	if len(fseasons) == 0 {
		return models.Show{}, models.ErrFailedToBuildShowDetails
	}

	detail.Seasons = make(map[string][]models.Episode)
	slices.Reverse(fseasons)
	for i, name := range fseasons {
		eps, err := utils.GetShowEpisodes(fseasons[i], show.Path)
		if err != nil {
			return models.Show{}, err
		}
		slices.Reverse(eps)
		detail.Seasons[name] = eps
	}

	return detail, nil
}

func (s *Server) GetShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 0 {
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

	page, err := buildLibraryPage(s.DB)
	if err != nil {
		s.internalError(w, "failed to build library page", err)
		return
	}

	s.render(w, r, http.StatusOK, "home.html", map[string]any{
		"Shows": page,
	})
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
	// cache the poster for a year and immutable
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(show.Poster)
}

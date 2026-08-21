package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"backend/db"
	"backend/models"
	"backend/utils"
	"backend/web"
)

func buildLibraryPage(database *db.SQL) (models.LibraryPage, error) {
	dbShows, err := database.GetAllShows()
	if err != nil {
		return models.LibraryPage{}, err
	}

	libShows := []models.LibraryShow{}
	for _, s := range dbShows {
		libShows = append(libShows, models.LibraryShow{
			ID:    s.ID,
			Title: s.Title,
			Image: "/api/poster/" + fmt.Sprint(s.ID),
		})
	}

	page := models.LibraryPage{
		TotalSeries: len(libShows),
		Shows:       libShows,
	}

	if len(dbShows) < 1 {
		return page, nil
	}

	return page, nil
}

func buildShowDetail(database *db.SQL, id int64) (models.ShowDetail, error) {
	show, err := database.GetShow(id)
	if err != nil {
		return models.ShowDetail{}, err
	}

	detail := models.ShowDetail{
		ID:    show.ID,
		Title: show.Title,
		Image: "/api/poster/" + fmt.Sprint(show.ID),
	}

	fseasons, err := utils.GetShowSeasons(show.Path)
	if err != nil {
		return models.ShowDetail{}, err
	}
	if len(fseasons) == 0 {
		return models.ShowDetail{}, models.ErrFailedToBuildShowDetails
	}

	detail.Seasons = []models.Season{}
	for i, name := range fseasons {
		eps, err := utils.GetShowEpisodes(fseasons[i], show.Path, int(show.ID))
		if err != nil {
			return models.ShowDetail{}, err
		}
		detail.EpisodeCount = len(eps)
		detail.Seasons = append(detail.Seasons, models.Season{
			Name:     name,
			Episodes: eps,
		})
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
		slog.Error("failed to build show detail", "id", id, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := web.RenderComponent(w, "home.html", "show-panel", detail); err != nil {
			slog.Error("failed to render show-panel", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	page, err := buildLibraryPage(s.DB)
	if err != nil {
		slog.Error("failed to build library page", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.render(w, r, http.StatusOK, "home.html", page)
}

func (s *Server) GetShowSeason(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	seasonIndex, err := strconv.Atoi(r.PathValue("season"))
	if err != nil || seasonIndex < 0 {
		http.NotFound(w, r)
		return
	}

	detail, err := buildShowDetail(s.DB, id)
	if err != nil {
		slog.Error("failed to build show season", "id", id, "season", seasonIndex, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.render(w, r, http.StatusOK, "season.html", detail)
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

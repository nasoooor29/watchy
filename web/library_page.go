package web

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"fmt"
)

func BuildLibraryPage(database *db.SQL) (models.LibraryPage, error) {
	dbShows, err := database.GetAllShows()
	if err != nil {
		return models.LibraryPage{}, err
	}

	libShows := make([]models.LibraryShow, 0, len(dbShows))
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

	selected, err := BuildShowDetail(database, dbShows[0].ID)
	if err != nil {
		return models.LibraryPage{}, err
	}
	page.Selected = selected

	return page, nil
}

func BuildShowDetail(database *db.SQL, id int64) (models.ShowDetail, error) {
	show, err := database.GetShow(id)
	if err != nil {
		return models.ShowDetail{}, err
	}

	detail := models.ShowDetail{
		Title: show.Title,
		Image: "/api/poster/" + fmt.Sprint(show.ID),
	}

	seasons, err := utils.GetShowSeasons(show.Path)
	if err != nil {
		return models.ShowDetail{}, err
	}
	if len(seasons) == 0 {
		return detail, nil
	}

	eps, err := utils.GetShowEpisodes(seasons[0], show.Path)
	if err != nil {
		return models.ShowDetail{}, err
	}
	detail.Season.Name = seasons[0]
	detail.Season.Episodes = eps
	detail.EpisodeCount = len(eps)

	return detail, nil
}

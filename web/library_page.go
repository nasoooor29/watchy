package web

import (
	"backend/db"
	"backend/models"
	"fmt"
)

func BuildLibraryPage(database *db.SQL) (models.LibraryPage, error) {
	dbShows, err := database.GetAllShows()
	if err != nil {
		return models.LibraryPage{}, err
	}

	var libShows []models.LibraryShow
	for _, s := range dbShows {
		libShows = append(libShows, models.LibraryShow{
			Title: s.Title,
			Image: "/api/poster/" + fmt.Sprint(s.Id),
		})
	}

	selected := models.ShowDetail{}
	if len(dbShows) > 0 {
		selected.Title = dbShows[0].Title
		selected.Image = "/api/poster/" + fmt.Sprint(dbShows[0].Id)
	}

	return models.LibraryPage{
		TotalSeries: len(libShows),
		Shows:       libShows,
		Selected:    selected,
	}, nil
}

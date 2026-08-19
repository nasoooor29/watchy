package source

import (
	"backend/db"
	"backend/models"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type Indexer struct {
	db   *db.SQL
	env  *models.EnvVars
	path string
}

func NewIndexer(path string, env *models.EnvVars, db *db.SQL) *Indexer {

	return &Indexer{
		path: path,
		env:  env,
		db:   db,
	}
}

func (i *Indexer) getShows() ([]string, error) {
	entries, err := os.ReadDir(i.path)
	if err != nil {
		return nil, err
	}

	var shows []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		shows = append(shows, entry.Name())
	}

	return shows, nil
}

func (i *Indexer) IndexShows() error {
	shows, err := i.getShows()
	if err != nil {
		return err
	}
	client := http.DefaultClient
	sources := NewSources(NewKitsuSource(i.env.KitsuBaseURL, client))
	for _, show := range shows {
		showPath := filepath.Join(i.path, show)
		_, err := i.db.GetShowByPath(showPath)
		if err == nil {
			slog.Debug("show", show, "already exists")
			continue
		}

		if err != sql.ErrNoRows {
			slog.Error("error happened", "err", err)
			continue
		}
		detail, err := sources.GetMetadata(show)
		if err != nil {
			return err
		}

		resp, err := client.Get(detail.Image)
		if err != nil {
			return err
		}

		poster, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		_, err = i.db.CreateShow(poster, filepath.Join(i.path, show), detail.Title, 0)
		if err != nil {
			return err
		}

	}
	return nil
}

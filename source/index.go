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
	db  *db.SQL
	env *models.EnvVars
}

func NewIndexer(env *models.EnvVars, db *db.SQL) *Indexer {

	return &Indexer{
		env: env,
		db:  db,
	}
}

func (i *Indexer) getShows() ([]string, error) {
	entries, err := os.ReadDir(i.env.TvDir)
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
	sources := NewSources(
		NewTenraiSource(i.env.TenraiBaseURL, client),
		NewKitsuSource(i.env.KitsuBaseURL, client),
	)
	for _, show := range shows {
		showPath := filepath.Join(i.env.TvDir, show)
		_, err := i.db.GetShowByPath(showPath)
		if err == nil {
			continue
		}

		if err != sql.ErrNoRows {
			slog.Error("error happened", "err", err)
			continue
		}
		detail, err := sources.GetMetadata(show)
		if err != nil {
			slog.Error("failed to fetch metadata", "show", show, "err", err)
			continue
		}

		resp, err := client.Get(detail.Image)
		if err != nil {
			slog.Error("failed to fetch poster", "show", show, "err", err)
			continue
		}

		poster, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			slog.Error("failed to read poster", "show", show, "err", err)
			continue
		}
		_, err = i.db.CreateShow(poster, filepath.Join(i.env.TvDir, show), detail.Title, 0)
		if err != nil {
			slog.Error("failed to insert show", "show", show, "err", err)
			continue
		}

	}
	return nil
}

func (i *Indexer) IndexShow(path string) (*db.Show, error) {
	client := http.DefaultClient
	sources := NewSources(
		NewTenraiSource(i.env.TenraiBaseURL, client),
		NewKitsuSource(i.env.KitsuBaseURL, client),
	)
	show, err := i.db.GetShowByPath(path)
	if err == nil {
		return show, nil
	}
	detail, err := sources.GetMetadata(show.Title)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(detail.Image)
	if err != nil {
		return nil, err
	}

	poster, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	show, err = i.db.CreateShow(poster, filepath.Join(i.env.TvDir, show.Path), detail.Title, 0)
	if err != nil {
		return nil, err
	}
	return show, nil
}

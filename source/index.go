package source

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Indexer struct {
	db  *db.SQL
	env *models.EnvVars
}

func NewIndexer(db *db.SQL) *Indexer {
	return &Indexer{
		env: models.GetEnv(),
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
			slog.Warn("skipping non-directory entry", "entry", entry.Name())
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			slog.Warn("skipping hidden directory", "dir", entry.Name())
			continue
		}
		if utils.IsHiddenShow(filepath.Join(i.env.TvDir, entry.Name())) {
			slog.Warn("skipping hidden show", "dir", entry.Name())
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
	for _, show := range shows {
		_, err := i.IndexShow(show)
		if err != nil {
			slog.Error("failed to index show", "show", show, "error", err)
		}
		slog.Info("successfully indexed show", "show", show)
	}
	return nil
}

func (i *Indexer) IndexShow(path string) (*db.Show, error) {
	showPath := filepath.Join(i.env.TvDir, path)
	client := http.DefaultClient
	sources := NewSources(
		NewTenraiSource(i.env.TenraiBaseURL, client),
		NewKitsuSource(i.env.KitsuBaseURL, client),
	)
	show, err := i.db.GetShowByPath(showPath)
	if err == nil {
		return show, nil
	}
	detail, err := sources.GetMetadata(filepath.Base(showPath))
	if err != nil {
		metadata, _ := json.Marshal(models.Metadata{Title: filepath.Base(showPath)})
		return i.db.CreateShow(nil, showPath, string(metadata))
	}
	resp, err := client.Get(detail.Image)
	if err != nil {
		return nil, err
	}

	poster, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	metadata, err := json.Marshal(detail)
	if err != nil {
		return nil, err
	}
	show, err = i.db.CreateShow(poster, showPath, string(metadata))
	if err != nil {
		return nil, err
	}
	return show, nil
}

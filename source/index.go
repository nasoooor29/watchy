package source

import (
	"backend/models"
	"fmt"
	"net/http"
	"os"
)

type Indexer struct {
	env  *models.EnvVars
	path string
}

func NewIndexer(path string, env *models.EnvVars) *Indexer {

	return &Indexer{
		path: path,
		env:  env,
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
	sources := NewSources(NewKitsuSource(i.env.KitsuBaseURL, http.DefaultClient))
	for _, show := range shows {
		detail, err := sources.GetMetadata(show)
		if err != nil {
			return err
		}
		fmt.Printf("show: %v\n", show)
		fmt.Printf("detail: %v\n", detail)

	}
	return nil
}

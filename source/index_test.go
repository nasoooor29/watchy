package source

import (
	"backend/models"
	"testing"
)

func TestIndexShows(t *testing.T) {
	idxr := NewIndexer("testdata/index", &models.EnvVars{
		KitsuBaseURL: "https://kitsu.io/api/edge",
	})
	err := idxr.IndexShows()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

}

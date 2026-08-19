package source

import (
	"backend/db"
	"backend/models"
	"testing"
)

func TestIndexShows(t *testing.T) {
	env := models.GetEnv()
	env.DatabaseURL = "/tmp/db.sqlite3"
	env.MigrationsDir = "../db/migrations/"
	db, err := db.GetDB(env)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	idxr := NewIndexer("testdata/index", env, db)
	err = idxr.IndexShows()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

}

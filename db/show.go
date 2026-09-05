package db

import (
	"backend/models"
	"encoding/json"
	"path/filepath"
)

type Show struct {
	ID     int64  `db:"id"`
	Poster []byte `db:"poster"`
	Path   string `db:"path"`
	Metadata string `db:"metadata"`
}

// GetMetadata converts the optional stored metadata into the application model.
func (s *Show) GetMetadata() models.Metadata {
	metadata := models.Metadata{Title: filepath.Base(s.Path)}
	if s.Metadata != "" {
		_ = json.Unmarshal([]byte(s.Metadata), &metadata)
	}
	return metadata
}

func (q *SQL) GetShow(id int64) (*Show, error) {
	show, err := Get[Show](q.db, "SELECT * FROM shows WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return show, nil
}

func (q *SQL) GetShowByPath(path string) (*Show, error) {
	return Get[Show](q.db, "SELECT * FROM shows WHERE path = ?", path)
}

func (q *SQL) GetAllShows() ([]Show, error) {
	return GetAll[Show](q.db, "SELECT * FROM shows")
}

func (q *SQL) CreateShow(poster []byte, path, metadata string) (*Show, error) {
	result, err := q.db.Exec(
		"INSERT INTO shows (poster, path, metadata) VALUES (?, ?, ?)", poster, path, metadata,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return q.GetShow(id)
}

package server

import (
	"backend/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// before: 600ms total
// after: 12ms max total
var metadataCache = make(map[string]map[string]any)
var ignoredCache = make(map[string]bool)

func getMetadataForShow(showDir string) (map[string]any, error) {
	env := models.GetEnv()
	metadataPath := filepath.Join(env.TvDir, showDir, "data.json")
	ignoredPath := filepath.Join(env.TvDir, showDir, ".hide-from-list")
	if _, ok := ignoredCache[ignoredPath]; ok {
		return nil, models.ErrHidden
	}
	if meta, ok := metadataCache[metadataPath]; ok {
		return meta, nil
	}
	if _, err := os.Stat(ignoredPath); err == nil {
		ignoredCache[ignoredPath] = true
		return nil, models.ErrHidden
	}

	meta := map[string]any{}
	_, err := os.Stat(metadataPath)
	if err != nil {
		slog.Warn("Show does not have metadata", "show", showDir)
		return nil, err
	}

	// metadata.json exists, read it and unmarshal it into a Show struct
	mdRaw, err := os.ReadFile(metadataPath)
	if err != nil {
		slog.Error("Failed to read metadata.json", "show", showDir, "err", err)
		return nil, err
	}

	err = json.Unmarshal(mdRaw, &meta)
	if err != nil {
		slog.Error("Failed to unmarshal metadata.json", "show", showDir, "err", err)
		return nil, err
	}

	metadataCache[metadataPath] = meta
	return meta, nil
}

func (s *Server) ShowsPage(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(s.ENV.TvDir)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	slog.Info("shows page", "dirs", len(dirs))
	shows := []models.Show{}
	withMD := 0
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		_, err := getMetadataForShow(dir.Name())
		if err != nil {
			slog.Error("Failed to get metadata for show", "show", dir.Name(), "err", err)
			continue
		} else {
			withMD++
		}

		shows = append(shows, models.Show{
			Title:    dir.Name(),
			Image:    "covers/" + dir.Name(),
			Episodes: []models.Episode{},
			// Metadata: utils.FlattenMap(meta), // NOTE: this will be from the metadata.json file if it exists or fetched from the DB
		})
	}
	// print the number of the shows that have metadata.json file
	slog.Info("shows page", "shows with metadata", withMD)
	s.render(w, r, http.StatusOK, "home.html", map[string]any{
		"TotalSeries": 5,
		"Shows":       shows,
	})
}

func (s *Server) GetShowCover(w http.ResponseWriter, r *http.Request) {
	coverPath := filepath.Join(s.ENV.TvDir, r.PathValue("ShowName"), "cover.jpg")

	if _, err := os.Stat(coverPath); err != nil {
		slog.Error("Failed to get cover for show", "show", r.PathValue("ShowName"), "err", err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// set the content type to image/jpeg
	w.Header().Set("Content-Type", "image/jpeg")

	// serve the cover.jpg file
	http.ServeFile(w, r, coverPath)
}

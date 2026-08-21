package utils

import (
	"backend/models"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nssteinbrenner/anitogo"
)

func GetShowSeasons(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var seasons []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		seasons = append(seasons, entry.Name())
	}

	return seasons, nil
}

func GetShowEpisodes(season, path string) ([]models.LibraryEpisode, error) {
	seasonPath := filepath.Join(path, season)
	slog.Debug("reading season directory", "path", seasonPath)

	entries, err := os.ReadDir(seasonPath)
	if err != nil {
		return nil, err
	}

	sourcesMap := make(map[int][]models.EpisodeSource)

	for _, entry := range entries {
		if !entry.IsDir() {
			slog.Error("failed to read season directory", "path", seasonPath, "err", err)
			continue
		}

		name := entry.Name()
		ele := anitogo.Parse(name, anitogo.DefaultOptions)
		if len(ele.EpisodeNumber) == 0 {
			slog.Warn("no episode number detected, skipping", "dir", name)
			continue
		}

		epNum, err := strconv.Atoi(ele.EpisodeNumber[0])
		if err != nil {
			slog.Warn("failed to parse episode number string", "dir", name, "val", ele.EpisodeNumber[0], "err", err)
			continue
		}

		sourcesMap[epNum] = append(sourcesMap[epNum], models.EpisodeSource{
			Label: entry.Name(),
		})
	}

	var episodes []models.LibraryEpisode
	for epNum, sources := range sourcesMap {
		episodes = append(episodes, models.LibraryEpisode{
			Name:    fmt.Sprintf("Episode %d", epNum),
			Watched: false,
			Sources: sources,
		})
	}

	slog.Info("loaded episodes successfully", "season", season, "count", len(episodes))
	return episodes, nil
}

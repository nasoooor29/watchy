package utils

import (
	"backend/models"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
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

func GetShowEpisodes(season, path string) ([]models.Episode, error) {
	seasonPath := filepath.Join(path, season)
	slog.Debug("reading season directory", "path", seasonPath)

	entries, err := os.ReadDir(seasonPath)
	if err != nil {
		return nil, err
	}

	pathsMap := make(map[int][]string)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ele := anitogo.Parse(name, anitogo.DefaultOptions)
		if len(ele.EpisodeNumber) == 0 {
			slog.Warn("no episode number detected, skipping", "file", name)
			continue
		}

		epNum, err := strconv.Atoi(ele.EpisodeNumber[0])
		if err != nil {
			slog.Warn("failed to parse episode number string", "file", name, "val", ele.EpisodeNumber[0], "err", err)
			continue
		}

		pathsMap[epNum] = append(pathsMap[epNum], filepath.Join(path, season, name))
	}

	keys := make([]int, 0, len(pathsMap))
	for k := range pathsMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	episodes := make([]models.Episode, 0, len(keys))
	for _, epNum := range keys {
		episodes = append(episodes, models.Episode{
			Index:   epNum,
			Watched: false,
			Paths: pathsMap[epNum],
		})
	}

	slog.Info("loaded episodes successfully", "season", season, "count", len(episodes))
	return episodes, nil
}

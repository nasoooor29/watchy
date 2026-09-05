package utils

import (
	"io/fs"
	"os"
)

func IsHiddenShow(path string) bool {
	if _, err := os.Stat(path + "/.hide-from-list"); err == nil {
		return true
	}
	return false
}

// GetLatestFileTimestamp returns the modification time of the newest regular
// file below path as Unix seconds.
func GetLatestFileTimestamp(path string) int64 {
	var latest int64
	found := false

	err := fs.WalkDir(os.DirFS(path), ".", func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		timestamp := info.ModTime().Unix()
		if !found || timestamp > latest {
			latest = timestamp
			found = true
		}
		return nil
	})
	if err != nil {
		return -1
	}
	return latest
}

package utils

import "os"

func IsHiddenShow(path string) bool {
	if _, err := os.Stat(path + "/.hide-from-list"); err == nil {
		return true
	}
	return false
}

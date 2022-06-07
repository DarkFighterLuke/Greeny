package utils

import (
	"fmt"
	"os"
	"sort"
)

func GetUserFolderPath() (string, error) {
	dirs, err := os.ReadDir("data")
	if err != nil {
		return "", err
	}

	if len(dirs) < 1 {
		return "", fmt.Errorf("no directories are present")
	}

	sort.Slice(dirs, func(i, j int) bool {
		infoI, err := dirs[i].Info()
		if err != nil {
			return false
		}
		infoJ, err := dirs[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	return dirs[0].Name(), nil
}

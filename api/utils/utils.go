package utils

import (
	"fmt"
	"os"
	"sort"
)

func GetUserFolderPath() (string, error) {
	nodes, err := os.ReadDir("data")
	if err != nil {
		return "", err
	}

	if len(nodes) < 1 {
		return "", fmt.Errorf("no directories are present")
	}

	sort.Slice(nodes, func(i, j int) bool {
		infoI, err := nodes[i].Info()
		if err != nil {
			return false
		}
		infoJ, err := nodes[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	for _, node := range nodes {
		if node.IsDir() {
			return node.Name(), nil
		}
	}
	return "", fmt.Errorf("no directories are present")
}

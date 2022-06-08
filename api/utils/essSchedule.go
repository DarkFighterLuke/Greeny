package utils

import (
	"encoding/csv"
	"os"
	"strconv"
)

type ESSSchedule = []float32

func parseESSSchedule(entries [][]string) (ESSSchedule, error) {
	var essSchedule []float32
	for i, entry := range entries {
		if i == 0 {
			continue
		}
		charge, err := strconv.ParseFloat(entry[1], 32)
		if err != nil {
			return nil, err
		}
		essSchedule = append(essSchedule, float32(charge))
	}
	return essSchedule, nil
}

func ReadESSSchedule(pathToESSFile string) (ESSSchedule, error) {
	file, err := os.Open(pathToESSFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	essSchedule, err := parseESSSchedule(entries)
	return essSchedule, err
}

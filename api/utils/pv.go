package utils

import (
	"encoding/csv"
	"os"
	"strconv"
)

type PVProduction = []float32

func parsePV(entries [][]string) (PVProduction, error) {
	var pv []float32
	for _, hourPV := range entries[1] {
		pvh, err := strconv.ParseFloat(hourPV, 32)
		if err != nil {
			return nil, err
		}
		pv = append(pv, float32(pvh))
	}
	return pv, nil
}

func ReadPV(pathToPVFile string) (PVProduction, error) {
	file, err := os.Open(pathToPVFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	pv, err := parsePV(entries)
	return pv, err
}

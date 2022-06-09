package utils

import (
	"encoding/csv"
	"os"
	"strconv"
)

func ReadInternalTemperature() (float32, error) {
	entries, err := readTemperatures()
	if err != nil {
		return 0.0, nil
	}

	internalTemperature, err := strconv.ParseFloat(entries[1][1], 32)
	if err != nil {
		return 0.0, err
	}
	return float32(internalTemperature), nil
}

func ReadExternalTemperature() (float32, error) {
	entries, err := readTemperatures()
	if err != nil {
		return 0.0, err
	}

	externalTemperature, err := strconv.ParseFloat(entries[1][0], 32)
	if err != nil {
		return 0.0, err
	}
	return float32(externalTemperature), nil
}

func readTemperatures() ([][]string, error) {
	temperatures, err := os.Open("data/simulated_temperatures.csv")
	if err != nil {
		return nil, err
	}

	entries, err := csv.NewReader(temperatures).ReadAll()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

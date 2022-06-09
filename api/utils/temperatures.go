package utils

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Temperatures struct {
	InternalTemperatures []float32
	ExternalTemperatures []float32
}

func ReadTemperatures() (Temperatures, error) {
	temperaturesFile, err := os.Open("data/simulated_temperatures.csv")
	if err != nil {
		return Temperatures{}, err
	}

	entries, err := csv.NewReader(temperaturesFile).ReadAll()
	if err != nil {
		return Temperatures{}, err
	}

	var temperatures Temperatures
	for _, externalTemperature := range entries[1][1:] {
		tempF, err := strconv.ParseFloat(externalTemperature, 32)
		if err != nil {
			return Temperatures{}, err
		}
		temperatures.ExternalTemperatures = append(temperatures.ExternalTemperatures, float32(tempF))
	}

	for _, internalTemperature := range entries[2][1:] {
		tempF, err := strconv.ParseFloat(internalTemperature, 32)
		if err != nil {
			return Temperatures{}, err
		}
		temperatures.InternalTemperatures = append(temperatures.InternalTemperatures, float32(tempF))
	}

	return temperatures, nil
}

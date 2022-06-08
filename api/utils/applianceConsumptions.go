package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Consumptions = []ConsumptionEntry

type ConsumptionEntry struct {
	ApplianceName      string
	HourlyConsumptions []float32
}

func parseConsumptions(entries [][]string) (Consumptions, error) {
	consumptions := Consumptions{}
	for _, entry := range entries {
		var hourlyConsumptions []float32
		for _, hourConsumption := range entry[1:] {
			consumption, err := strconv.ParseFloat(hourConsumption, 32)
			if err != nil {
				return nil, err
			}
			hourlyConsumptions = append(hourlyConsumptions, float32(consumption))
		}
		consumptionEntry := ConsumptionEntry{
			ApplianceName:      entry[0],
			HourlyConsumptions: hourlyConsumptions,
		}
		consumptions = append(consumptions, consumptionEntry)
	}
	return consumptions, nil
}

func ReadConsumptions(pathToConsumptionsFile string) (Consumptions, error) {
	file, err := os.Open(pathToConsumptionsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	consumptions, err := parseConsumptions(entries)
	return consumptions, err
}

func (ce *ConsumptionEntry) isTurnedOn(hour int) (bool, error) {
	if hour < 0 || hour > 23 {
		return false, fmt.Errorf("hour out of range")
	}
	return ce.HourlyConsumptions[hour] != 0, nil
}

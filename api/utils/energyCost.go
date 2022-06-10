package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	strconv "strconv"
)

type EnergyCost = [][]float32

func parseEnergyCost(entries [][]string) (EnergyCost, error) {
	var energyCost EnergyCost

	// Sunday energy cost parsing
	var dayCostsFloat []float32
	for j, hourCost := range entries[7] {
		if j == 0 {
			continue
		}
		cost, err := strconv.ParseFloat(hourCost, 32)
		if err != nil {
			return nil, err
		}
		dayCostsFloat = append(dayCostsFloat, float32(cost))
	}
	energyCost = append(energyCost, dayCostsFloat)

	for i := 1; i < 7; i++ {
		var dayCostsFloat []float32
		for j, hourCost := range entries[i] {
			if j == 0 {
				continue
			}
			cost, err := strconv.ParseFloat(hourCost, 32)
			if err != nil {
				return nil, err
			}
			dayCostsFloat = append(dayCostsFloat, float32(cost))
		}
		energyCost = append(energyCost, dayCostsFloat)
	}
	return energyCost, nil
}

func ReadEnergyCost(pathToEnergyCostFile string) (EnergyCost, error) {
	file, err := os.Open(pathToEnergyCostFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	energyCost, err := parseEnergyCost(entries)
	return energyCost, err
}

func GetEnergyCostByDayAndHour(energyCost *EnergyCost, dayOfWeek, hour int) (float32, error) {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return 0, fmt.Errorf("day of week out of range")
	}

	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("hour out of range")
	}

	return (*energyCost)[dayOfWeek][hour], nil
}

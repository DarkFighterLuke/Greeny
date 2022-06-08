package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Consumptions = []ConsumptionEntry

type ConsumptionEntry struct {
	ApplianceName      string
	HourlyConsumptions []float32
}

func parseConsumptions(entries [][]string) (Consumptions, error) {
	consumptions := Consumptions{}
	for i, entry := range entries {
		if i == 0 {
			continue
		}
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

func FindConsumptionsByApplianceName(consumptions *Consumptions, commonName string) (*ConsumptionEntry, error) {
	for i, entry := range *consumptions {
		if strings.ToLower(entry.ApplianceName) == strings.ToLower(commonName) {
			return &(*consumptions)[i], nil
		}
	}
	return nil, fmt.Errorf("no appliance found with the given name")
}

func WriteConsumptionsToCsv(consumptions *Consumptions, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	entries := consumptionsTo2DArray(consumptions)
	err = w.WriteAll(entries)
	return err
}

func consumptionsTo2DArray(consumptions *Consumptions) [][]string {
	var array [][]string
	array = append(array, getTimeHeaderArray())
	for _, entry := range *consumptions {
		var arrayEntry []string
		arrayEntry = append(arrayEntry, entry.ApplianceName)
		for i := 0; i < 24; i++ {
			arrayEntry = append(arrayEntry, fmt.Sprintf("%.2f", entry.HourlyConsumptions[i]))
		}
		array = append(array, arrayEntry)
	}
	return array
}

func getTimeHeaderArray() []string {
	timeArray := []string{"Orario"}
	for i := 0; i < 24; i++ {
		timeArray = append(timeArray, fmt.Sprintf("%d", i))
	}
	return timeArray
}

func (ce *ConsumptionEntry) IsTurnedOn(hour int) (bool, error) {
	if hour < 0 || hour > 23 {
		return false, fmt.Errorf("hour out of range")
	}
	return ce.HourlyConsumptions[hour] != 0, nil
}

func (ce *ConsumptionEntry) GetPowerOnInterval() (startHour, endHour int, err error) {
	startHour = -1
	endHour = -1
	for hour, hourConsumption := range ce.HourlyConsumptions {
		if hourConsumption != 0 {
			startHour = hour
			break
		}
	}

	if startHour == -1 {
		err = fmt.Errorf("appliance has no consumptions")
		return
	}

	for i := startHour; i < len(ce.HourlyConsumptions); i++ {
		if ce.HourlyConsumptions[i] == 0 {
			endHour = i
		}
	}
	if endHour == -1 {
		endHour = 23
	}

	return
}

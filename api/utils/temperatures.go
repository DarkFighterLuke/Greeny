package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	deltaTemperatureHealthHazardDegrees = 15
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

func IsColderOutside(hour int) (bool, error) {
	temperatures, err := ReadTemperatures()
	if err != nil {
		return false, err
	}

	if temperatures.ExternalTemperatures[hour] < temperatures.InternalTemperatures[hour] {
		return true, nil
	} else {
		return false, nil
	}
}

func IsWantedTemperatureLowerThanOutside(wantedTemperature float32, hour int) (bool, error) {
	temperatures, err := ReadTemperatures()
	if err != nil {
		return false, err
	}

	return wantedTemperature < temperatures.ExternalTemperatures[hour], nil
}

func IsWantedTemperatureLowerThanInside(wantedTemperature float32, hour int) (bool, error) {
	temperatures, err := ReadTemperatures()
	if err != nil {
		return false, err
	}

	return wantedTemperature < temperatures.InternalTemperatures[hour], nil
}

func GetInternalTemperatureByHour(hour int) (float32, error) {
	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("hour out of range")
	}

	temperatures, err := ReadTemperatures()
	if err != nil {
		return 0, err
	}

	return temperatures.InternalTemperatures[hour], nil
}

func GetExternalTemperatureByHour(hour int) (float32, error) {
	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("hour out of range")
	}

	temperatures, err := ReadTemperatures()
	if err != nil {
		return 0, err
	}

	return temperatures.ExternalTemperatures[hour], nil
}

func IsDeltaTemperatureDangerous(wantedTemperature float32, hour int) (bool, bool, error) {
	if hour < 0 || hour > 23 {
		return false, false, fmt.Errorf("hour out of range")
	}

	externalTemperature, err := GetExternalTemperatureByHour(hour)
	if err != nil {
		return false, false, err
	}

	isWantedTooLow := externalTemperature > wantedTemperature

	return math.Abs(float64(externalTemperature-wantedTemperature)) >= deltaTemperatureHealthHazardDegrees, isWantedTooLow, nil
}

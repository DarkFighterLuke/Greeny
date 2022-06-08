package utils

import (
	"errors"
	"os"
	"os/exec"
)

func GenerateConsumptionsFiles() error {
	userFolderName, err := GetUserFolderPath()
	if err != nil {
		return err
	}

	basePath := "data/" + userFolderName + "/"
	summary, err := ReadSummaryFile(basePath + "summary.csv")
	if err != nil {
		return err
	}

	simulationConsumptions, err := ReadConsumptions("data/simulation.csv")
	if err != nil {
		return err
	}

	shiftableEntries := GetShiftableEntries(&summary)
	var shiftableConsumptions Consumptions
	for _, entry := range shiftableEntries {
		entryConsumptions, err := FindConsumptionsByApplianceName(&simulationConsumptions, entry.CommonName)
		if err != nil {
			continue
		}
		shiftableConsumptions = append(shiftableConsumptions, *entryConsumptions)
	}
	err = WriteConsumptionsToCsv(&shiftableConsumptions, basePath+"shiftable.csv")
	if err != nil {
		return err
	}

	nonShiftableEntries := GetNonShiftableEntries(&summary)
	var nonShiftableConsumptions Consumptions
	for _, entry := range nonShiftableEntries {
		entryConsumptions, err := FindConsumptionsByApplianceName(&simulationConsumptions, entry.CommonName)
		if err != nil {
			continue
		}
		nonShiftableConsumptions = append(nonShiftableConsumptions, *entryConsumptions)
	}
	nonShiftableConsumptions = append(nonShiftableConsumptions, simulationConsumptions[len(simulationConsumptions)-2:]...)
	err = WriteConsumptionsToCsv(&nonShiftableConsumptions, basePath+"non-shiftable.csv")
	if err != nil {
		return err
	}

	return nil
}

func GenerateOptimalSchedule(shiftablePath, nonShiftablePath, optimizedSchedulePath string) error {
	if _, err := os.Stat("../GAEnergyOptimizer/venv"); errors.Is(err, os.ErrNotExist) {
		_, err := exec.Command("python3", "-m", "venv", "../GAEnergyOptimizer/venv").Output()
		if err != nil {
			return err
		}
		_, err = exec.Command("../GAEnergyOptimizer/venv/bin/pip3", "install", "-r", "../GAEnergyOptimizer/requirements.txt").Output()
		if err != nil {
			return err
		}
	}

	_, err := exec.Command("../GAEnergyOptimizer/venv/bin/python3", "../GAEnergyOptimizer/main.py", shiftablePath, nonShiftablePath, "data/user_preferred_time_slots.csv", "-o", optimizedSchedulePath).Output()
	if err != nil {
		return err
	}

	return nil
}

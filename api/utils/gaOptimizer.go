package utils

import (
	"encoding/csv"
	"os"
)

const (
	nonShiftableStart = 1
	nonShiftableEnd   = 9
	shiftableStart    = 9
	shiftableEnd      = 13
)

func GenerateConsumptionsFiles() error {
	userFolderName, err := GetUserFolderPath()
	if err != nil {
		return err
	}

	path := "data/" + userFolderName + "/" + userFolderName + ".csv"
	summary, err := ReadSummaryFile(path)
	if err != nil {
		return err
	}

	var shiftableSummary Summary
	var nonShiftableSummary Summary
	for _, entry := range summary {
		if entry.Shiftable {
			shiftableSummary = append(shiftableSummary, entry)
		} else {
			nonShiftableSummary = append(nonShiftableSummary, entry)
		}
	}

	file, err := os.Open("data/simulation.csv")
	if err != nil {
		return err
	}

	r := csv.NewReader(file)
	simulationEntries, err := r.ReadAll()
	if err != nil {
		return err
	}

	path = "data/" + userFolderName
	shiftableFile, err := os.OpenFile(path+"/shiftable.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer shiftableFile.Close()
	w := csv.NewWriter(shiftableFile)
	defer w.Flush()
	err = w.WriteAll(simulationEntries[shiftableStart:shiftableEnd])
	if err != nil {
		return err
	}

	nonShiftableFile, err := os.OpenFile(path+"/non-shiftable.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer nonShiftableFile.Close()
	w = csv.NewWriter(nonShiftableFile)
	defer w.Flush()
	err = w.WriteAll(simulationEntries[nonShiftableStart:nonShiftableEnd])
	if err != nil {
		return err
	}

	return nil
}

package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Summary = []SummaryEntry

type SummaryEntry struct {
	Appliance  string
	Shiftable  bool
	Priority   int
	SetupDone  bool
	CommonName string
}

func ParseSummary(entries [][]string) (Summary, error) {
	summary := Summary{}
	for _, entry := range entries {
		shiftable, err := strconv.ParseBool(entry[1])
		if err != nil {
			return nil, err
		}
		priority, err := strconv.ParseInt(entry[2], 10, 5)
		if err != nil {
			return nil, err
		}
		if priority < 1 || priority > 10 {
			return nil, fmt.Errorf("priority out of range")
		}
		setupDone, err := strconv.ParseBool(entry[3])
		if err != nil {
			return nil, err
		}
		summaryEntry := SummaryEntry{
			Appliance:  entry[0],
			Shiftable:  shiftable,
			Priority:   int(priority),
			SetupDone:  setupDone,
			CommonName: entry[4],
		}
		summary = append(summary, summaryEntry)
	}
	return summary, nil
}

func FindFirstUnconfigured(summary *Summary) (*SummaryEntry, error) {
	for _, entry := range *summary {
		if !entry.SetupDone {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("all entries are configured")
}

func ReadSummaryFile(pathToSummaryFile string) (Summary, error) {
	file, err := os.Open(pathToSummaryFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	summary, err := ParseSummary(entries)
	return summary, err
}

func WriteToCsv(summary *Summary, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	entries := To2DArray(summary)
	err = w.WriteAll(entries)
	return err
}

func To2DArray(summary *Summary) [][]string {
	var array [][]string
	for _, entry := range *summary {
		arrayEntry := make([]string, 5)
		arrayEntry[0] = entry.Appliance
		arrayEntry[1] = strconv.FormatBool(entry.Shiftable)
		arrayEntry[2] = strconv.Itoa(entry.Priority)
		arrayEntry[3] = strconv.FormatBool(entry.SetupDone)
		arrayEntry[4] = entry.CommonName
		array = append(array, arrayEntry)
	}
	return array
}

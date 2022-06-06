package controllers

import (
	"encoding/csv"
	"errors"
	"greeny/main/utils"
	"os"
)

func CreateUser(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	path := "main/data"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
	}

	path = "main/data/" + request.QueryResult.Parameters["username"]["name"]
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
	}

	f, err := os.Create("main/data/" + request.QueryResult.Parameters["username"]["name"] + "/" +
		request.QueryResult.Parameters["username"]["name"] + ".csv")
	defer f.Close()
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	records := [][]string{
		{"appliance", "shiftable", "priority", "setup", "common_name"},
		{"piano_cottura", "", "", "", "Piano cottura"},
		{"forno", "", "", "", "Forno"},
		{"microonde", "", "", "", "Microonde"},
		{"computer_desktop", "", "", "", "Computer Desktop"},
		{"caricabatterie_cellulare", "", "", "", "Caricabatterie per cellulare"},
		{"televisione", "", "", "", "Televisore"},
		{"scaldabagno", "", "", "", "Scaldabagno"},
		{"condizionatore", "", "", "", "Condizionatore"},
		{"lavatrice", "", "", "", "Lavatrice"},
		{"asciugatrice", "", "", "", "Asciugatrice"},
		{"lavastoviglie", "", "", "", "Lavastoviglie"},
		{"aspirapolvere", "", "", "", "Aspirapolvere"},
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	if err = w.WriteAll(records); err != nil {
		return utils.WebhookResponse{}, err
	}

	response := utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Username creato"},
				},
			},
		},
	}
	return response, nil
}

func imReadyForSetup(request utils.WebhookRequest) (utils.WebhookResponse, error) {

	return utils.WebhookResponse{}, nil
}

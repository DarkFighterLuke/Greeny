package controllers

import (
	"encoding/csv"
	"errors"
	"greeny/utils"
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

func AmIReadyForSetup(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	var response utils.WebhookResponse
	/*if request.QueryResult.Parameters[""] != "" {
		response = utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Allora cominciamo!\n\n" +
							"Sono riuscito a recuperare la lista degli elettrodomestici della casa, adesso ti" +
							"chiederò per ciascuno di essi di esprimere che priorità hanno in relazione agli altri" +
							"su una scala da 1 a 10 e se posso ripianificare la loro accensione in altre ore del" +
							"giorno per farti risparmiare.\n" +
							"Nella lista ho trovato l’elettrodomestico " + records[1][4] + ". Che priorità ha per te da 1 a 10?"},
					},
				},
			},
			OutputContext: []utils.Context{
				{
					Name:          "setup",
					LifespanCount: 1,
				},
				{
					Name:          "priority",
					LifespanCount: 1,
				},
			},
		}
	} else {
		response = utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Va bene alla prossima!"},
					},
				},
			},
			OutputContext: []utils.Context{},
		}
	}
	*/
	return response, nil
}

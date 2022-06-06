package controllers

import (
	"encoding/csv"
	"errors"
	"greeny/utils"
	"os"
	"reflect"
)

func CreateUser(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	response := utils.WebhookResponse{}
	path := "data/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
	}

	v := reflect.ValueOf(request.QueryResult.Parameters["username"])
	if v.Kind() == reflect.Map {
		path = path + request.QueryResult.Parameters["username"].(map[string]interface{})["name"].(string)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return utils.WebhookResponse{}, err
			}
		}

		f, err := os.Create(path + "/" + request.QueryResult.Parameters["username"].(map[string]interface{})["name"].(string) + ".csv")
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

		response = utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Hai davvero un bel nome! Penso che rimarrà impresso nei miei circuiti." +
							"\nEssendo un nuovo arrivato ho bisogno di farti alcune domande per iniziare ad aiutarti." +
							"\nSei pronto?"},
					},
				},
			},
			OutputContext: []utils.Context{
				{
					Name:          "setup",
					LifespanCount: 5,
				},
				{
					Name:          "ready_request",
					LifespanCount: 1,
				},
				{
					Name:          "Setup-Name-followup",
					LifespanCount: 2,
				},
			},
		}
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

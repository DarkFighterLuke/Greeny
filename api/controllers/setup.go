package controllers

import (
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

		summary := utils.Summary{
			{
				Appliance:  "piano_cottura",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Piano cottura",
			},
			{
				Appliance:  "forno",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Forno",
			},
			{
				Appliance:  "microonde",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Microonde",
			},
			{
				Appliance:  "computer",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Computer Desktop",
			},
			{
				Appliance:  "caricabatterie_cellulare",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Caricabatterie per cellulare",
			},
			{
				Appliance:  "televisione",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Televisore",
			},
			{
				Appliance:  "scaldabagno",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Scaldabagno",
			},
			{
				Appliance:  "condizionatore",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Condizionatore",
			},
			{
				Appliance:  "lavatrice",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Lavatrice",
			},
			{
				Appliance:  "asciugatrice",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Asciugatrice",
			},
			{
				Appliance:  "lavastoviglie",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Lavastoviglie",
			},
			{
				Appliance:  "aspirapolvere",
				Shiftable:  false,
				Priority:   1,
				SetupDone:  false,
				CommonName: "Aspirapolvere",
			},
		}

		if err := utils.WriteToCsv(&summary, path+"/"+request.QueryResult.Parameters["username"].(map[string]interface{})["name"].(string)+".csv"); err != nil {
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
	if request.QueryResult.Parameters[""] != "" && request.QueryResult.Parameters["false"] == "" {
		// TODO: Read summary file and check for the first unconfigured appliance
		/*response = utils.WebhookResponse{
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
					Name:          "appliance_priority_request",
					LifespanCount: 1,
				},
			},
		}*/
	} else {
		response = utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Va bene, alla prossima!"},
					},
				},
			},
			OutputContext: []utils.Context{},
		}
	}

	return response, nil
}

package controllers

import (
	"errors"
	"fmt"
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

		if err := utils.WriteToCsv(&summary, path+"/"+"summary.csv"); err != nil {
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
	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		userFolderName, err := utils.GetUserFolderPath()
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		summary, err := utils.ReadSummaryFile("data/" + userFolderName + "/" + "summary.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		unconfigured, err := utils.FindFirstUnconfigured(&summary)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		response = utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Allora cominciamo!\n\n" +
							"Sono riuscito a recuperare la lista degli elettrodomestici della casa, adesso ti " +
							"chiederò per ciascuno di essi di esprimere che priorità hanno in relazione agli altri " +
							"su una scala da 1 a 10 e se posso ripianificare la loro accensione in altre ore del " +
							"giorno per farti risparmiare.\n" +
							"Nella lista ho trovato l’elettrodomestico " + unconfigured.CommonName + ". Che priorità ha per te da 1 a 10?"},
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
		}
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

func AppliancePriority(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	priority := request.QueryResult.Parameters["priority"].(float64)
	if priority < 1 || priority > 10 {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"I miei circuiti non riescono a comprendere il numero che hai specificato. " +
							"Potresti ripetere la priorità su una scala da 1 a 10?"},
					},
				},
			},
			OutputContext: []utils.Context{
				{
					Name:          "setup",
					LifespanCount: 1,
				},
				{
					Name:          "Setup-ReadyAnswer-followup",
					LifespanCount: 1,
				},
				{
					Name:          "appliance_priority_request",
					LifespanCount: 1,
				},
			},
		}, fmt.Errorf("priority number out of allowed range")
	}

	userFolderName, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	path := "data/" + userFolderName + "/" + "summary.csv"
	summary, err := utils.ReadSummaryFile(path)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	unconfigured, err := utils.FindFirstUnconfigured(&summary)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	unconfigured.Priority = int(priority)
	err = utils.WriteToCsv(&summary, path)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Bene.\nPosso ripianificare la sua accensione in altre ore del giorno?"},
				},
			},
		},
		OutputContext: []utils.Context{
			{
				Name:          "setup",
				LifespanCount: 1,
			},
			{
				Name:          "appliance_shiftability_request",
				LifespanCount: 1,
			},
		},
	}, nil
}

func ApplianceShiftability(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	userFolderName, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	path := "data/" + userFolderName + "/" + "summary.csv"
	summary, err := utils.ReadSummaryFile(path)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	unconfigured, err := utils.FindFirstUnconfigured(&summary)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	if request.QueryResult.Parameters["true"] == nil || request.QueryResult.Parameters["false"] == nil {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Non ho capito ciò che hai detto...\n" +
							"Posso ripianificare la sua accensione in altre ore del giorno?"},
					},
				},
			},
			OutputContext: []utils.Context{
				{
					Name:          "setup",
					LifespanCount: 1,
				},
				{
					Name:          "Setup-ReadyAnswer-followup",
					LifespanCount: 1,
				},
				{
					Name:          "Setup-AppliancePriority-followup",
					LifespanCount: 1,
				},
				{
					Name:          "appliance_shiftability_request",
					LifespanCount: 1,
				},
			},
		}, nil
	} else if request.QueryResult.Parameters["false"] == "" {
		unconfigured.Shiftable = true
		unconfigured.SetupDone = true
	} else {
		unconfigured.SetupDone = true
	}

	err = utils.WriteToCsv(&summary, path)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	unconfigured, err = utils.FindFirstUnconfigured(&summary)
	if err != nil {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Grazie per la pazienza.\n" +
							"Potresti dirmi ora quali elettrodomestici tra quelli citati usi " +
							"per regolare la temperatura della casa?"},
					},
				},
			},
			OutputContext: []utils.Context{
				{
					Name:          "setup",
					LifespanCount: 1,
				},
				{
					Name:          "temperature_setters_request",
					LifespanCount: 1,
				},
			},
		}, nil
	}
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Afferrato!\n" +
						"Ho trovato l’elettrodomestico " + unconfigured.CommonName + ". Che priorità ha per te da 1 a 10?"},
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
	}, nil
}

func TemperatureSetters(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	if request.QueryResult.Parameters["appliances"] != nil && len(request.QueryResult.Parameters["appliances"].([]interface{})) > 0 {
		userFolderName, err := utils.GetUserFolderPath()
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		path := "data/" + userFolderName + "/" + "summary.csv"
		summary, err := utils.ReadSummaryFile(path)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		appliances := request.QueryResult.Parameters["appliances"].([]interface{})
		for _, commonName := range appliances {
			entry, err := utils.FindEntryByCommonName(&summary, commonName.(string))
			if err != nil {
				return utils.WebhookResponse{
					FulfillmentMessages: []utils.Message{
						{
							Text: utils.Text{
								Text: []string{"Non ho trovato gli elettrodomestici di cui parli. Se vuoi che ti ripeta la lista basta dirlo!"},
							},
						},
					},
					OutputContext: []utils.Context{
						{
							Name:          "setup",
							LifespanCount: 1,
						},
						{
							Name:          "temperature_setters_request",
							LifespanCount: 1,
						},
					},
				}, nil
			}
			entry.TemperatureSetter = true
		}

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Afferrato!" +
							"\nGrazie per tua pazienza Elisabetta e scusami se sono stato chiacchierone." +
							"\nTi ricordo che i tuoi dati sono al sicuro con me, non dirò a nessuno del nostro segreto per risparmiare."},
					},
				},
			},
		}, nil
	} else {
		return utils.WebhookResponse{}, fmt.Errorf("no parameters have been supplied")
	}
}

func RepeatAppliances(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	userFolderName, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	path := "data/" + userFolderName + "/" + "summary.csv"
	summary, err := utils.ReadSummaryFile(path)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	responseMessage := "Non c'è problema! Gli elettrodomestici nella tua casa sono: "
	for i, entry := range summary {
		if i != 0 {
			responseMessage += ", "
		}
		responseMessage += entry.CommonName
	}

	responseMessage += "\nQuale di questi usi per regolare la temperatura della casa?"

	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{responseMessage},
				},
			},
		},
		OutputContext: []utils.Context{
			{
				Name:          "setup",
				LifespanCount: 1,
			},
			{
				Name:          "temperature_setters_request",
				LifespanCount: 1,
			},
		},
	}, nil
}

// TODO: Generate shiftable and non-shiftable files
// TODO: Generate optimal schedule
// TODO: Copy file with random internal and external temperatures
// TODO: Copy file with user preferred time slots
// TODO: Copy file with PV production

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
				Appliance:                 "piano_cottura",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Piano cottura",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "forno",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Forno",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "microonde",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Microonde",
				NeedsTemperatureToPowerOn: true,
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
				CommonName: "Televisione",
			},
			{
				Appliance:                 "scaldabagno",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Scaldabagno",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "condizionatore",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Condizionatore",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "lavatrice",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Lavatrice",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "asciugatrice",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Asciugatrice",
				NeedsTemperatureToPowerOn: true,
			},
			{
				Appliance:                 "lavastoviglie",
				Shiftable:                 false,
				Priority:                  1,
				SetupDone:                 false,
				CommonName:                "Lavastoviglie",
				NeedsTemperatureToPowerOn: true,
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
						Text: []string{"Hai davvero un bel nome! Penso che rimarr?? impresso nei miei circuiti." +
							"\nEssendo un nuovo arrivato ho bisogno di farti alcune domande per iniziare ad aiutarti." +
							"\nSei pronto?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          request.Session + "/contexts/setup",
					LifespanCount: 2,
				},
				{
					Name:          request.Session + "/contexts/ready_request",
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
							"chieder?? per ciascuno di essi di esprimere che priorit?? hanno in relazione agli altri " +
							"su una scala da 1 a 10 e se posso ripianificare la loro accensione in altre ore del " +
							"giorno per farti risparmiare.\n" +
							"Nella lista ho trovato l???elettrodomestico " + unconfigured.CommonName + ". Che priorit?? ha per te da 1 a 10?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          request.Session + "/contexts/setup",
					LifespanCount: 2,
				},
				{
					Name:          request.Session + "/contexts/appliance_priority_request",
					LifespanCount: 2,
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
			OutputContexts: []utils.Context{},
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
							"Potresti ripetere la priorit?? su una scala da 1 a 10?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          request.Session + "/contexts/setup",
					LifespanCount: 2,
				},
				{
					Name:          request.Session + "/contexts/appliance_priority_request",
					LifespanCount: 2,
				},
			},
		}, nil
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
		OutputContexts: []utils.Context{
			{
				Name:          request.Session + "/contexts/setup",
				LifespanCount: 2,
			},
			{
				Name:          request.Session + "/contexts/appliance_shiftability_request",
				LifespanCount: 2,
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
						Text: []string{"Non ho capito ci?? che hai detto...\n" +
							"Posso ripianificare la sua accensione in altre ore del giorno?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          request.Session + "/contexts/setup",
					LifespanCount: 2,
				},
				{
					Name:          request.Session + "/contexts/appliance_shiftability_request",
					LifespanCount: 2,
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
			OutputContexts: []utils.Context{
				{
					Name:          request.Session + "/contexts/setup",
					LifespanCount: 2,
				},
				{
					Name:          request.Session + "/contexts/temperature_setters_request",
					LifespanCount: 2,
				},
			},
		}, nil
	}
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Afferrato!\n" +
						"Ho trovato l???elettrodomestico " + unconfigured.CommonName + ". Che priorit?? ha per te da 1 a 10?"},
				},
			},
		},
		OutputContexts: []utils.Context{
			{
				Name:          request.Session + "/contexts/setup",
				LifespanCount: 2,
			},
			{
				Name:          request.Session + "/contexts/appliance_priority_request",
				LifespanCount: 2,
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
			entry, err := utils.FindSummaryEntryByCommonName(&summary, commonName.(string))
			if err != nil {
				return utils.WebhookResponse{
					FulfillmentMessages: []utils.Message{
						{
							Text: utils.Text{
								Text: []string{"Non ho trovato gli elettrodomestici di cui parli. Se vuoi che ti ripeta la lista basta dirlo!"},
							},
						},
					},
					OutputContexts: []utils.Context{
						{
							Name:          request.Session + "/contexts/setup",
							LifespanCount: 2,
						},
						{
							Name:          request.Session + "/contexts/temperature_setters_request",
							LifespanCount: 2,
						},
					},
				}, nil
			}
			entry.TemperatureSetter = true
		}

		err = utils.WriteToCsv(&summary, path)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = utils.GenerateConsumptionsFiles()
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		err = utils.GenerateOptimalSchedule("data/"+userFolderName+"/shiftable.csv", "data/"+userFolderName+"/non-shiftable.csv", "data/"+userFolderName+"/optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Afferrato!" +
							"\nGrazie per tua pazienza " + userFolderName + " e scusami se sono stato chiacchierone." +
							"\nTi ricordo che i tuoi dati sono al sicuro con me, non dir?? a nessuno del nostro segreto per risparmiare." +
							"\nCosa posso fare per te?"},
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

	responseMessage := "Non c'?? problema! Gli elettrodomestici nella tua casa sono: "
	for i, entry := range summary {
		if i != 0 {
			responseMessage += ", "
		}
		responseMessage += entry.CommonName
	}

	responseMessage += ".\nQuale di questi usi per regolare la temperatura della casa?"

	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{responseMessage},
				},
			},
		},
		OutputContexts: []utils.Context{
			{
				Name:          request.Session + "/contexts/setup",
				LifespanCount: 2,
			},
			{
				Name:          request.Session + "/contexts/temperature_setters_request",
				LifespanCount: 2,
			},
		},
	}, nil
}

// TODO: Copy file with PV production

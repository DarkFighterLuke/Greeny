package controllers

import (
	"fmt"
	"greeny/utils"
	"os"
	"sort"
	"strings"
	"time"
)

func AppliancePowerOn(request utils.WebhookRequest, doTemperatureCheck bool) (utils.WebhookResponse, error) {
	currentHour := time.Now().Hour()
	currentDayOfWeek := int(time.Now().Weekday())
	var applianceName string
	if request.QueryResult.Parameters["appliance"] == nil {
		context, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
			request.Session, "power_on_request")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceName = context.Parameters["appliance"].(string)
	} else {
		applianceName = request.QueryResult.Parameters["appliance"].(string)
	}

	var temperatureParameters map[string]interface{}
	var temperature float32
	if request.QueryResult.Parameters["temperature"] == nil {
		context, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
			request.Session, "power_on_request")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		temperatureParameters = context.Parameters["temperature"].(map[string]interface{})
	} else {
		switch request.QueryResult.Parameters["temperature"].(type) {
		case string:
			temperatureParameters = nil
		default:
			temperatureParameters = request.QueryResult.Parameters["temperature"].(map[string]interface{})
			temperature = float32(temperatureParameters["amount"].(float64))
		}
	}

	userFolderName, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	basePath := "data/" + userFolderName + "/"
	summary, err := utils.ReadSummaryFile(basePath + "summary.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	if !utils.IsSetupCompleted(&summary) {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Ho rilevato che non hai ancora completato la fase di setup iniziale. " +
							"Per far girare i miei ingranaggi ho bisogno di quelle importanti informazioni!\n" +
							"Per favore, avvia la fase di setup dicendo \"setup\"."},
					},
				},
			},
		}, nil
	}

	summaryAppliance, err := utils.FindSummaryEntryByCommonName(&summary, applianceName)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	if summaryAppliance.NeedsTemperatureToPowerOn && temperatureParameters == nil {
		outputContexts := request.QueryResult.OutputContexts
		outputContexts = append(outputContexts, utils.Context{
			Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "temperature_request"),
			LifespanCount: 1,
		})

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{fmt.Sprintf("A che temperatura vuoi che accenda l'elettrodomestico %s?", applianceName)},
					},
				},
			},
			OutputContexts: outputContexts,
		}, nil
	} else if !summaryAppliance.NeedsTemperatureToPowerOn && temperatureParameters != nil {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{fmt.Sprintf("Mi dispiace %s, ma questo elettrodomestico non prevede l'impostazione di una "+
							"temperatura...", userFolderName)},
					},
				},
			},
		}, nil
	}

	consumptions, err := utils.ReadConsumptions(basePath + "optimal-schedule.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	applianceConsumptions, _ := utils.FindConsumptionsByApplianceName(&consumptions, applianceName)
	if err != nil {
		return utils.WebhookResponse{}, nil
	}

	startHour, endHour, err := applianceConsumptions.GetPowerOnInterval()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	reAmount, err := calculateCurrentRenewableEnergyAmount(currentHour, userFolderName)
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	currentTotalConsumption, err := calculateCurrentTotalConsumption(currentHour, &consumptions)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	energyCost, err := utils.ReadEnergyCost("data/italian_electricity_cost.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	currentEnergyCost, err := utils.GetEnergyCostByDayAndHour(&energyCost, currentDayOfWeek, currentHour)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	powerOnContext, _ := utils.FindContextByName(&request.QueryResult.OutputContexts, request.Session, "power_on_request")
	powerOnContext.LifespanCount = 1

	if summaryAppliance.TemperatureSetter && doTemperatureCheck {
		isDangerous, isWantedTooLow, err := utils.IsDeltaTemperatureDangerous(temperature, currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		healthAlertMessage := ""
		if isDangerous {
			if isWantedTooLow {
				healthAlertMessage = "La temperatura esterna è di gran lunga superiore a quella che mi chiedi, questo " +
					"potrebbe nuocere alla tua salute."
			} else {
				healthAlertMessage = "La temperatura interna è di gran lunga superiore a quella che mi chiedi, questo " +
					"potrebbe nuocere alla tua salute."
			}
		}

		isTemperatureLowerThanInside, err := utils.IsWantedTemperatureLowerThanInside(temperature, currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		isColderOutside, err := utils.IsColderOutside(currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		internalTemperature, err := utils.GetInternalTemperatureByHour(currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		externalTemperature, err := utils.GetExternalTemperatureByHour(currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		nreAmount := currentTotalConsumption + applianceConsumptions.HourlyConsumptions[startHour] - reAmount
		nreCost := nreAmount * currentEnergyCost

		var savingMessage string
		if nreAmount > 0 {
			savingMessage = fmt.Sprintf("In questa maniera risparmieresti %.2f€ e l'immissione di CO2 nell'ambiente.", nreCost)
		}

		if isTemperatureLowerThanInside {
			if isColderOutside {
				return utils.WebhookResponse{
					FulfillmentMessages: []utils.Message{
						{
							Text: utils.Text{
								Text: []string{fmt.Sprintf("%s \nLa temperatura interna attualmente è di %.1f°C, "+
									"quella esterna di %1.f. Se vuoi un po’ di fresco il mio consiglio è di aprire le "+
									"finestre invece di accendere l'elettrodomestico %s. %s\nChe ne dici?", healthAlertMessage,
									internalTemperature, externalTemperature, applianceName, savingMessage)},
							},
						},
					},
					OutputContexts: []utils.Context{
						{
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "advice_windows_open_request"),
							LifespanCount: 1,
						},
						powerOnContext,
					},
				}, nil
			}
		} else {
			if !isColderOutside {
				return utils.WebhookResponse{
					FulfillmentMessages: []utils.Message{
						{
							Text: utils.Text{
								Text: []string{fmt.Sprintf("%s \nLa temperatura interna attualmente è di %.1f°C, "+
									"quella esterna di %1.f. Se vuoi un po’ di caldo il mio consiglio è di aprire le "+
									"finestre invece di accendere l'elettrodomestico %s. In questa maniera risparmieresti "+
									"%.2f€ e l'immissione di CO2 nell'ambiente.\nChe ne dici?", healthAlertMessage,
									internalTemperature, externalTemperature, applianceName, nreCost)},
							},
						},
					},
					OutputContexts: []utils.Context{
						{
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "advice_windows_open_request"),
							LifespanCount: 1,
						},
						powerOnContext,
					},
				}, nil
			}
		}
	}

	isOn, err := applianceConsumptions.IsTurnedOn(currentHour)
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	if isOn {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"L'elettrodomestico che mi hai chiesto di accendere è già acceso"},
					},
				},
			},
		}, nil
	} else {
		if summaryAppliance.Shiftable {
			currentApplianceCost := currentEnergyCost * applianceConsumptions.HourlyConsumptions[startHour]
			scheduledEnergyCost, err := utils.GetEnergyCostByDayAndHour(&energyCost, currentDayOfWeek, startHour)
			if err != nil {
				return utils.WebhookResponse{}, err
			}
			scheduledApplianceCost := scheduledEnergyCost * applianceConsumptions.HourlyConsumptions[startHour]
			savingPercentage := (currentApplianceCost - scheduledApplianceCost) * 100 / currentApplianceCost
			var responseMessage string
			if currentApplianceCost == scheduledApplianceCost {
				responseMessage += fmt.Sprintf("Se accendessi l'elettrodomestico %s adesso la pianficazione di tutti gli altri elettrodomestici "+
					"potrebbe cambiare.\n", applianceName)
			} else {
				responseMessage += fmt.Sprintf("%s ho fatto due calcoli per te. "+
					"Se accendessi l'elettrodomestico %s adesso spenderesti %.2f€, mentre, se lo accendessi "+
					"dalle ore %d alle ore %d, spenderesti soltanto %.2f€, un risparmio del %.2f%%.\n",
					userFolderName, applianceName, currentApplianceCost, startHour, endHour, scheduledApplianceCost, savingPercentage)
			}

			if currentTotalConsumption+applianceConsumptions.HourlyConsumptions[startHour] > reAmount {
				responseMessage += "Inoltre sforeresti la produzione attuale di energia rinnovabile della casa, " +
					"eliminando le possibilità di risparmio per tutti gli altri dispositivi.\n"
			}
			responseMessage += "Vuoi comunque accendere l'elettrodomestico " + applianceName + "?"

			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							Text: []string{responseMessage},
						},
					},
				},
				OutputContexts: []utils.Context{
					powerOnContext,
					{
						Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "shiftable_power_on_confirm_request"),
						LifespanCount: 1,
					},
				},
			}, nil
		} else {
			nreAmount := currentTotalConsumption + applianceConsumptions.HourlyConsumptions[startHour] - reAmount
			nreCost := nreAmount * currentEnergyCost
			if nreAmount > 0 {
				responseMessage := fmt.Sprintf("%s, se accendi l'elettrodomestico %s adesso sforerai la "+
					"quantità di energia rinnovabile a tua disposizione, dovendo così prelevarla dalla rete "+
					"elettrica, per un costo di %.2f€.\n", userFolderName, applianceName, nreCost)

				appliancesToPowerOff, err := calculateAppliancesToPowerOff(&summary, &consumptions, currentHour, summaryAppliance, nreAmount)
				if err != nil {
					return utils.WebhookResponse{}, err
				}
				if appliancesToPowerOff == nil || len(appliancesToPowerOff) == 0 {
					responseMessage += "Sfortunatamente non sono riuscito a trovare un modo possibile per risparmiare. " +
						"Vuoi procedere lo stesso?"
					return utils.WebhookResponse{
						FulfillmentMessages: []utils.Message{
							{
								Text: utils.Text{
									Text: []string{responseMessage},
								},
							},
						},
						OutputContexts: []utils.Context{
							powerOnContext,
							{
								Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "nre_turn_on_request"),
								LifespanCount: 1,
							},
						},
					}, nil
				}
				var appliancesNames []string
				for _, applianceConsumption := range appliancesToPowerOff {
					appliancesNames = append(appliancesNames, applianceConsumption.ApplianceName)
				}
				responseMessage += fmt.Sprintf("Tuttavia, se spegnessi %s, non dovresti "+
					"acquistare energia dalla rete elettrica.\nVuoi che li spenga?", strings.Join(appliancesNames, ", "))

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
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "power_on_request"),
							LifespanCount: 1,
						},
						{
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "power_off_appliances_request"),
							LifespanCount: 1,
							Parameters: map[string]interface{}{
								"appliances": appliancesNames,
							},
						},
					},
				}, nil
			} else {
				err = utils.PowerOnNonShiftable(userFolderName, applianceName, currentHour)
				if err != nil {
					return utils.WebhookResponse{}, err
				}

				err = utils.GenerateOptimalSchedule(basePath+"shiftable.csv", basePath+"non-shiftable_temp.csv", basePath+"optimal-schedule.csv")
				if err != nil {
					return utils.WebhookResponse{}, err
				}

				err = os.Remove(basePath + "/non-shiftable_temp.csv")
				if err != nil {
					return utils.WebhookResponse{}, err
				}

				return utils.WebhookResponse{
					FulfillmentMessages: []utils.Message{
						{
							Text: utils.Text{
								Text: []string{"Sfortunatamente non sono riuscito a trovare un modo possibile per risparmiare. Procedo all'accensione.\nPosso fare altro per te?"},
							},
						},
					},
				}, nil
			}
		}
	}
}

func calculateCurrentRenewableEnergyAmount(currentHour int, userFolderName string) (float32, error) {
	if currentHour < 0 || currentHour > 23 {
		return 0, fmt.Errorf("hour out of range")
	}

	pv, err := utils.ReadPV("data/pv.csv")
	if err != nil {
		return 0, err
	}

	ess, err := utils.ReadESSSchedule("data/" + userFolderName + "/optimal-schedule_ess.csv")
	if err != nil {
		return 0, err
	}

	return pv[currentHour] + ess[currentHour], nil
}

func calculateCurrentTotalConsumption(currentHour int, consumptions *utils.Consumptions) (float32, error) {
	if currentHour < 0 || currentHour > 23 {
		return 0, fmt.Errorf("hour out of range")
	}

	var currentTotalConsumption float32
	for _, applianceConsumption := range *consumptions {
		currentTotalConsumption += applianceConsumption.HourlyConsumptions[currentHour]
	}

	return currentTotalConsumption, nil
}

func calculateAppliancesToPowerOff(summary *utils.Summary, consumptions *utils.Consumptions, currentHour int, applianceToPowerOn *utils.SummaryEntry, energyToSave float32) (utils.Consumptions, error) {
	if applianceToPowerOn.Priority < 1 || applianceToPowerOn.Priority > 10 {
		return nil, fmt.Errorf("priority out of range")
	}

	if currentHour < 0 || currentHour > 23 {
		return nil, fmt.Errorf("hour out of range")
	}

	if energyToSave <= 0 {
		return nil, fmt.Errorf("nothing to power off")
	}

	var lowerEqualPriorityConsumptions []*utils.ConsumptionEntry
	for _, entry := range *summary {
		if entry.Priority <= applianceToPowerOn.Priority && entry.CommonName != applianceToPowerOn.CommonName {
			consumptionEntry, err := utils.FindConsumptionsByApplianceName(consumptions, entry.CommonName)
			if err != nil {
				continue
			}
			lowerEqualPriorityConsumptions = append(lowerEqualPriorityConsumptions, consumptionEntry)
		}
	}

	if len(lowerEqualPriorityConsumptions) == 0 {
		return utils.Consumptions{}, nil
	}

	sort.Slice(lowerEqualPriorityConsumptions, func(i, j int) bool {
		return lowerEqualPriorityConsumptions[i].HourlyConsumptions[currentHour] < lowerEqualPriorityConsumptions[j].HourlyConsumptions[currentHour]
	})

	var sum float32
	var appliancesToPowerOff utils.Consumptions
	i := 0
	for sum < energyToSave && i < len(lowerEqualPriorityConsumptions) {
		if lowerEqualPriorityConsumptions[i].HourlyConsumptions[currentHour] > 0 {
			sum += lowerEqualPriorityConsumptions[i].HourlyConsumptions[currentHour]
			appliancesToPowerOff = append(appliancesToPowerOff, *lowerEqualPriorityConsumptions[i])
		}
		i++
	}

	if sum < energyToSave {
		return nil, nil
	}

	return appliancesToPowerOff, nil
}

func NREUsageConfirmation(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		currentHour := time.Now().Hour()

		userFolderName, err := utils.GetUserFolderPath()
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		context, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
			request.Session, "power_on_request")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceName := context.Parameters["appliance"].(string)

		basePath := "data/" + userFolderName + "/"

		nonShiftableConsumptions, err := utils.ReadConsumptions(basePath + "/non-shiftable.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceConsumptions, err := utils.FindConsumptionsByApplianceName(&nonShiftableConsumptions, applianceName)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		hourlyConsumptionIndex, _, err := applianceConsumptions.GetPowerOnInterval()
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceConsumptions.HourlyConsumptions[currentHour] = applianceConsumptions.HourlyConsumptions[hourlyConsumptionIndex]
		nonShiftableConsumptions = *utils.ReplaceConsumptionsEntry(&nonShiftableConsumptions, applianceConsumptions)

		err = utils.WriteConsumptionsToCsv(&nonShiftableConsumptions, basePath+"/non-shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = utils.GenerateOptimalSchedule(basePath+"/shiftable.csv", basePath+"/non-shiftable_temp.csv",
			basePath+"/optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = os.Remove(basePath + "/non-shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Ok " + userFolderName + ", procedo all'accensione."},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	} else {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Grazie per la tua scelta green, il pianeta te ne è grato.\n" +
							"Posso fare altro per te?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	}
}

func RecommendedPowerOffConfirmation(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	powerOnContext, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
		request.Session, "power_on_request")
	powerOnContext.LifespanCount = 1
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		currentHour := time.Now().Hour()

		userFolderName, err := utils.GetUserFolderPath()
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		applianceName := powerOnContext.Parameters["appliance"].(string)

		pwOffContext, err := utils.FindContextByName(&request.QueryResult.OutputContexts, request.Session,
			"power_off_appliances_request")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		appliancesToPowerOff := pwOffContext.Parameters["appliances"].([]interface{})

		basePath := "data/" + userFolderName + "/"

		summary, err := utils.ReadSummaryFile(basePath + "summary.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		for _, appliance := range appliancesToPowerOff {
			summaryEntry, err := utils.FindSummaryEntryByCommonName(&summary, appliance.(string))
			if err != nil {
				return utils.WebhookResponse{}, err
			}
			if summaryEntry.Shiftable {
				err = utils.PowerOffShiftable(userFolderName, appliance.(string), currentHour)
				if err != nil {
					return utils.WebhookResponse{}, err
				}
			} else {
				err = utils.PowerOffNonShiftable(userFolderName, appliance.(string), currentHour)
				if err != nil {
					return utils.WebhookResponse{}, err
				}
			}
		}
		err = utils.PowerOnNonShiftable(userFolderName, applianceName, currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = utils.GenerateOptimalSchedule(basePath+"shiftable_temp.csv", basePath+"non-shiftable_temp.csv", basePath+"optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = os.Remove(basePath + "/shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		err = os.Remove(basePath + "/non-shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Grazie per la tua scelta green, il pianeta te ne è grato.\n" +
							"Posso fare altro per te?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	} else {
		var outputContexts []utils.Context
		outputContexts = append(outputContexts, powerOnContext)
		outputContexts = append(outputContexts, utils.Context{
			Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "nre_turn_on_request"),
			LifespanCount: 1,
		})
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Vuoi procedere comunque con l'accensione?"},
					},
				},
			},
			OutputContexts: outputContexts,
		}, nil
	}
}

func ProceedToShiftablePowerOn(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		powerOnContext, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
			request.Session, "power_on_request")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		currentHour := time.Now().Hour()

		userFolderName, err := utils.GetUserFolderPath()
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		applianceName := powerOnContext.Parameters["appliance"].(string)

		err = utils.PowerOnShiftable(userFolderName, applianceName, currentHour)
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		basePath := "data/" + userFolderName + "/"
		err = utils.GenerateOptimalSchedule(basePath+"shiftable_temp.csv", basePath+"non-shiftable_temp.csv", basePath+"optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = os.Remove(basePath + "/shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = os.Remove(basePath + "/non-shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{fmt.Sprintf("Ok %s, procedo.\nPosso fare altro per te?", userFolderName)},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	} else {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Grazie per la tua scelta green, il pianeta te ne è grato.\n" +
							"Posso fare altro per te?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	}
}

func AdviceWindows(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Grazie per la tua scelta green, il pianeta te ne è grato.\n" +
							"Posso fare altro per te?"},
					},
				},
			},
			OutputContexts: []utils.Context{
				{
					Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "can_i_do_something_else_request"),
					LifespanCount: 1,
				},
			},
		}, nil
	} else {
		return AppliancePowerOn(request, false)
	}
}

package controllers

import (
	"fmt"
	"greeny/utils"
	"os"
	"sort"
	"strings"
	"time"
)

func AppliancePowerOn(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	currentHour := time.Now().Hour()
	currentDayOfWeek := int(time.Now().Weekday())
	var applianceName string
	if request.QueryResult.Parameters["appliance"] == nil {
		context, err := utils.FindContextByName(&request.QueryResult.OutputContexts,
			fmt.Sprintf(utils.ContextsBase, request.Session, "power_on_request"))
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceName = context.Parameters["appliance"].(string)
	} else {
		applianceName = request.QueryResult.Parameters["appliance"].(string)
	}
	temperature := request.QueryResult.Parameters["temperature"]

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

	if summaryAppliance.NeedsTemperatureToPowerOn && temperature == "" {
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
	} else if !summaryAppliance.NeedsTemperatureToPowerOn && temperature != "" {
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

	if summaryAppliance.TemperatureSetter {
		// TODO: Implement appliance temperature setter case
		return utils.WebhookResponse{}, fmt.Errorf("non implemented yet")
	} else {
		consumptions, err := utils.ReadConsumptions(basePath + "optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		applianceConsumptions, _ := utils.FindConsumptionsByApplianceName(&consumptions, applianceName)
		if err != nil {
			return utils.WebhookResponse{}, nil
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

			if summaryAppliance.Shiftable {
				currentApplianceCost := currentEnergyCost * applianceConsumptions.HourlyConsumptions[startHour]
				scheduledEnergyCost, err := utils.GetEnergyCostByDayAndHour(&energyCost, currentDayOfWeek, startHour)
				if err != nil {
					return utils.WebhookResponse{}, err
				}
				scheduledApplianceCost := scheduledEnergyCost * applianceConsumptions.HourlyConsumptions[startHour]
				savingPercentage := (currentApplianceCost - scheduledApplianceCost) * 100 / currentApplianceCost
				responseMessage := fmt.Sprintf("%s ho fatto due calcoli per te. "+
					"Se accendessi l'elettrodomestico %s adesso spenderesti %.2f€, mentre, se lo accendessi"+
					"dalle ore %d alle ore %d, spenderesti soltanto %.2f€, un risparmio del %.2f%%.\n",
					strings.ToTitle(strings.ToLower(userFolderName)),
					applianceName, currentApplianceCost, startHour, endHour, scheduledApplianceCost, savingPercentage)

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
						{
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "power_on_request"),
							LifespanCount: 1,
						},
						{
							Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "shiftable_proceed_request_step_1"),
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
					if len(appliancesToPowerOff) == 0 {
						responseMessage += "Vuoi procedere lo stesso?"
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
									Name:          fmt.Sprintf(utils.ContextsBase, request.Session, "nre_turn_on_request"),
									LifespanCount: 1,
								},
							},
						}, nil
					}
					var appliancesMessage string
					for _, applianceConsumption := range appliancesToPowerOff {
						appliancesMessage += applianceConsumption.ApplianceName + ", "
					}
					responseMessage += fmt.Sprintf("Tuttavia, se spegnessi %s non dovresti "+
						"acquistare energia dalla rete elettrica.\nVuoi che li spenga?", appliancesMessage)

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
									"appliances": appliancesToPowerOff,
								},
							},
						},
					}, nil
				}
			}
		}
	}
	return utils.WebhookResponse{}, fmt.Errorf("")
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
		sum += lowerEqualPriorityConsumptions[i].HourlyConsumptions[currentHour]
		appliancesToPowerOff = append(appliancesToPowerOff, *lowerEqualPriorityConsumptions[i])
		i++
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
			fmt.Sprintf(utils.ContextsBase, request.Session, "power_on_request"))
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
		fmt.Sprintf(utils.ContextsBase, request.Session, "power_on_request"))
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

		pwOffContext, err := utils.FindContextByName(&request.QueryResult.OutputContexts, fmt.Sprintf(utils.ContextsBase, request.Session,
			"power_off_appliances_request"))
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
			Name:          "nre_turn_on_request",
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

package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

func WhatIsTheTemperature(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	temperatures, err := utils.ReadTemperatures()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := time.Now().Hour()
	strTemperatures := "La temperatura interna è di " + fmt.Sprintf("%.1f", temperatures.InternalTemperatures[hour]) +
		"°C mentre quella esterna è di " +
		fmt.Sprintf("%.1f", temperatures.ExternalTemperatures[hour]) + "°C"

	if request.QueryResult.Parameters["perceived-temperature-feeling"] != nil {
		if request.QueryResult.Parameters["perceived-temperature-feeling"].(string) == "freddo" {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							Text: []string{
								"Brrrr! Mi si potrebbero congelare i circuiti. " + strTemperatures,
							},
						},
					},
				},
			}, nil
		} else if request.QueryResult.Parameters["perceived-temperature-feeling"].(string) == "caldo" {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							Text: []string{
								"Brrrr! Mi si potrebbero sciogliere i circuiti. " + strTemperatures,
							},
						},
					},
				},
			}, nil
		}
	}
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{
						strTemperatures,
					},
				},
			},
		},
	}, nil
}

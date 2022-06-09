package controllers

import (
	"fmt"
	"greeny/utils"
)

func WhatIsTheTemperature(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	internalTemperature, err := utils.ReadInternalTemperature()
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	exrternalTemperature, err := utils.ReadExternalTemperature()
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	strTemperatures := "La temperatura interna è di " + fmt.Sprintf("%.1f", internalTemperature) +
		"°C mentre quella esterna è di " +
		fmt.Sprintf("%.1f", exrternalTemperature) + "°C"

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

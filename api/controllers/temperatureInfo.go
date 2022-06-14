package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

// WhatIsTheTemperature returns the internal and the external temperature.
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
								"Brrrr! Mi si potrebbero congelare i circuiti. " + strTemperatures +
									"\nPosso fare altro per te?",
							},
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
		} else if request.QueryResult.Parameters["perceived-temperature-feeling"].(string) == "caldo" {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							Text: []string{
								"Che caldo! Mi si potrebbero sciogliere i circuiti. " + strTemperatures +
									"\nPosso fare altro per te?",
							},
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
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{
						strTemperatures + "\nPosso fare altro per te?",
					},
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

package controllers

import (
	"greeny/utils"
)

// getAgentName creates a response for the get-agent-name intent.
func GetAgentName(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	response := utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"My name is Greeny"},
				},
			},
		},
	}
	return response, nil
}

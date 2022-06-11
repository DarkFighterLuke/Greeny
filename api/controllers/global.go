package controllers

import "greeny/utils"

func CanIDoSomethingElse(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	if request.QueryResult.Parameters["false"] != nil && request.QueryResult.Parameters["false"] == "" {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Cosa posso fare per te?"},
					},
				},
			},
		}, nil
	} else {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Ok vado in stand-by, chiamami quando hai bisogno di me!"},
					},
				},
			},
		}, nil
	}
}

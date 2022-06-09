package controllers

import "greeny/utils"

func GlobalFallback(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	var outputContexts []utils.Context
	for _, context := range request.QueryResult.OutputContexts {
		context.LifespanCount = 2
		outputContexts = append(outputContexts, context)
	}

	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Non ho capito bene, potresti ripetere?"},
				},
			},
		},
		OutputContexts: outputContexts,
	}, nil
}

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

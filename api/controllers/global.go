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

package controllers

import (
	"greeny/utils"
	"time"
)

func OpenWindows(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	colder, err := utils.IsColderOutside(time.Now().Hour())
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	if request.QueryResult.Parameters["perceived-temperature-feeling"].(string) == "caldo" {
		if colder {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							[]string{"Fuori è più freddo, aprendo le finestre stai attuando una scelta green." +
								" Complimenti!"},
						},
					},
				},
			}, nil
		} else {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							[]string{"Ti consiglio di non aprire perchè fuori la temperatura è più alta" +
								" e la casa potrebbe riscaldarsi!"},
						},
					},
				},
			}, nil
		}
	} else {
		if !colder {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							[]string{"Fuori è più caldo, aprendo le finestre stai attuando una scelta green." +
								" Complimenti!"},
						},
					},
				},
			}, nil
		} else {
			return utils.WebhookResponse{
				FulfillmentMessages: []utils.Message{
					{
						Text: utils.Text{
							[]string{"Ti consiglio di non aprire perchè fuori la temperatura è più bassa" +
								" e la casa potrebbe raffreddarsi!"},
						},
					},
				},
			}, nil
		}
	}
}

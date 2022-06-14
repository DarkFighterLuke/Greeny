package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

var pathToPVFile = "data/pv.csv"

// CurrentEnergyProduction this function return an utils.WebhookRequest with the current energy production, where
// current means time.Now()
func CurrentEnergyProduction(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	pvs, err := utils.ReadPV(pathToPVFile)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := time.Now().Hour()
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					[]string{"La produzione corrente di energia è di " + fmt.Sprintf("%.2f", pvs[hour]) + "kW" +
						"\nPosso fare altro per te?"},
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

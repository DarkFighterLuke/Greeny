package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

func AvailableEnergyAmount(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	pvs, err := utils.ReadPV(pathToPVFile)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := time.Now().Hour()
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					[]string{"La produzione corrente di energia è di " + fmt.Sprintf("%.2f", pvs[hour]) + "kW"},
				},
			},
		},
	}, nil
}

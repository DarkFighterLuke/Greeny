package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

func AvailableEnergyAmount(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	user, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	pathToESSFile := "data/" + user + "/optimal-schedule_ess.csv"
	ess, err := utils.ReadESSSchedule(pathToESSFile)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := time.Now().Hour()
	return utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					[]string{"La produzione corrente di energia è di " + fmt.Sprintf("%.2f", ess[hour]) + "kW"},
				},
			},
		},
	}, nil
}

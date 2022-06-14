package controllers

import (
	"fmt"
	"greeny/utils"
	"time"
)

// AvailableEnergyAmount this function return an utils.WebhookRequest with the current battery residual energy, where
// current means time.Now()
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
					[]string{"Le batterie hanno " + fmt.Sprintf("%.2f", ess[hour]) + "kW di energia disponibili" +
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

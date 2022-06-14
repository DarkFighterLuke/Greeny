package controllers

import (
	"fmt"
	"greeny/utils"
	"os"
	"time"
)

// AppliancePowerOff this function switch off an appliance and then calculates the new optimal-schedule file.
// It accepts an utils.WebhookRequest and return an utils.WebhookResponse and an error if it is present.
func AppliancePowerOff(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	user, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	appliance := request.QueryResult.Parameters["appliance"].(string)

	basePath := "data/" + user + "/"
	summary, err := utils.ReadSummaryFile(basePath + "summary.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	if !utils.IsSetupCompleted(&summary) {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"Ho rilevato che non hai ancora completato la fase di setup iniziale. " +
							"Per far girare i miei ingranaggi ho bisogno di quelle importanti informazioni!\n" +
							"Per favore, avvia la fase di setup dicendo \"setup\"."},
					},
				},
			},
		}, nil
	}

	consumptions, err := utils.ReadConsumptions("data/" + user + "/optimal-schedule.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	consumption, err := utils.FindConsumptionsByApplianceName(&consumptions, appliance)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := time.Now().Hour()
	on, err := consumption.IsTurnedOn(hour)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	if on {
		shiftable, err := utils.IsApplianceShiftable(user, appliance)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		if shiftable {
			err := utils.PowerOffShiftable(user, appliance, hour)
			if err != nil {
				return utils.WebhookResponse{}, err
			}
		} else {
			err := utils.PowerOffNonShiftable(user, appliance, hour)
			if err != nil {
				return utils.WebhookResponse{}, err
			}
		}

		err = utils.GenerateOptimalSchedule("data/"+user+"/shiftable_temp.csv", "data/"+user+
			"/non-shiftable_temp.csv", "data/"+user+"/optimal-schedule.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}

		err = os.Remove("data/" + user + "/shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		err = os.Remove("data/" + user + "/non-shiftable_temp.csv")
		if err != nil {
			return utils.WebhookResponse{}, err
		}
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"L'elettrodomestico " + appliance + " è stato spento.\n" +
							"Posso fare altro per te?"},
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
	} else {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"L'elettrodomestico " + appliance + " è già spento.\n" +
							"Posso fare altro per te?"},
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

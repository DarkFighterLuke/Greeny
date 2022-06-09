package controllers

import (
	"greeny/utils"
	"os"
)

func PowerOff(response utils.WebhookRequest) (utils.WebhookResponse, error) {
	user, err := utils.GetUserFolderPath()
	if err != nil {
		return utils.WebhookResponse{}, err
	}
	appliance := response.QueryResult.Parameters["appliance"].(string)

	consumptions, err := utils.ReadConsumptions("data/" + user + "/optimal-schedule.csv")
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	consumption, err := utils.FindConsumptionsByApplianceName(&consumptions, appliance)
	if err != nil {
		return utils.WebhookResponse{}, err
	}

	hour := 21 //time.Now().Hour()
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
						Text: []string{"L'elettrodomestico " + appliance + " è stato spento"},
					},
				},
			},
		}, nil
	} else {
		return utils.WebhookResponse{
			FulfillmentMessages: []utils.Message{
				{
					Text: utils.Text{
						Text: []string{"L'elettrodomestico " + appliance + " è già spento"},
					},
				},
			},
		}, nil
	}
}

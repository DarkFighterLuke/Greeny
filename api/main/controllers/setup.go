package controllers

import (
	"errors"
	"greeny/main/utils"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func CreateUser(request utils.WebhookRequest) (utils.WebhookResponse, error) {
	path := "main/data"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
	}

	rand.Seed(time.Now().Unix())
	path = "main/data/" + request.QueryResult.Parameters["username"]["name"] + "_" + strconv.Itoa(rand.Int())
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return utils.WebhookResponse{}, err
		}
	}
	response := utils.WebhookResponse{
		FulfillmentMessages: []utils.Message{
			{
				Text: utils.Text{
					Text: []string{"Username creato"},
				},
			},
		},
	}
	return response, nil
}

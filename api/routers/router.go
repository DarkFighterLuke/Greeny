package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"greeny/controllers"
	"greeny/utils"
	"log"
	"net/http"
	"strings"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/Greeny/api", handleWebhookRequest).Methods("POST")
	return r
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func handleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var request utils.WebhookRequest
	var response utils.WebhookResponse
	var err error

	// Read input JSON
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Request: %+v", request)

	// Call intent handler
	intent := strings.Split(request.QueryResult.Intent.Name, "/")
	switch intent[len(intent)-1] {
	// Use intent-id to identify it
	case "e7a823e2-b2ba-49c2-9caa-c7946ff647c2":
		// Setup - Name
		response, err = controllers.CreateUser(request)
	case "063b3da9-3a2f-45bf-bda3-63eeb1939c82":
		// Setup - Ready answer
		response, err = controllers.AmIReadyForSetup(request)
	case "833bdcb2-83f7-4224-92b5-8d1de3319660":
		// Setup - Appliance priority
		response, err = controllers.AppliancePriority(request)
	case "7b1a49f1-8243-4199-b1d6-84a0f2587f38":
		// Setup - Appliance shiftability
		response, err = controllers.ApplianceShiftability(request)
	case "3cd37b30-e6f4-4730-94df-f0ab8cbf1dd1":
		// Setup - Temperature setters
		response, err = controllers.TemperatureSetters(request)
	case "cc0fd5cb-7568-4f78-97e7-81b931284019":
		// Setup - Repeat appliances
		response, err = controllers.RepeatAppliances(request)
	case "bef393dc-2c73-45e1-a72f-6a9235d2c92e", "7b78b8d8-fef2-43f7-8e12-c0bddf7e57fe":
		// Appliance power on
		// Appliance power on - Temperature request
		response, err = controllers.AppliancePowerOn(request)
	case "132187da-5035-45ea-ab71-6b82e04e03ae":
		// Appliance power off
		response, err = controllers.PowerOff(request)
	case "d088ed94-9df9-40c8-99b1-93ee1cab4fe7":
		// Appliance power on - NRE usage confirmation
		response, err = controllers.NREUsageConfirmation(request)
	case "59c9f470-ebb6-48c2-87db-e1bb5fd1f7d5":
		// Appliance power on - Power off lower priority appliances
		response, err = controllers.RecommendedPowerOffConfirmation(request)
	case "660e8d42-a9d1-4f29-aa78-1b700771b640":
		// Temperature info
		response, err = controllers.WhatIsTheTemperature(request)
	case "d40cc55f-c73d-4b7b-bc6d-62248d4c94c6":
		// Currently energy production
		response, err = controllers.CurrentlyEnergyProduction(request)
	case "c86594a7-de4a-46fe-a1fb-0d66546de2b1":
		// Available energy amount
		response, err = controllers.AvailableEnergyAmount(request)
	case "187ee8f9-b9eb-4a03-86bc-6344d127490c":
		// Open windows
		response, err = controllers.OpenWindows(request)
	default:
		err = fmt.Errorf("Unknown intent: %s", intent)
	}
	if err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Response: %+v", response)

	// Send response
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		handleError(w, err)
		return
	}
}

// handleError handles internal errors.
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ERROR: %v", err)
}

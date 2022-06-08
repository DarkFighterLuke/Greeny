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

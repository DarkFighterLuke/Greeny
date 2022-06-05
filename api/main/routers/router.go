package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
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
	var request *dialogflow.WebhookRequest
	var response *dialogflow.WebhookResponse
	var err error

	// Read input JSON
	print(r.Body)
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Request: %+v", request)

	// Call intent handler
	intent := strings.Split(request.QueryResult.Intent.Name, "/")
	switch intent[len(intent)-1] {
	// Use intent-id to identify it
	case "<intent-id>":
		response, err = getAgentName(request)
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

// getAgentName creates a response for the get-agent-name intent.
func getAgentName(request *dialogflow.WebhookRequest) (*dialogflow.WebhookResponse, error) {
	response := &dialogflow.WebhookResponse{
		FulfillmentMessages: []*dialogflow.Intent_Message{
			{
				Message: &dialogflow.Intent_Message_Text_{
					Text: &dialogflow.Intent_Message_Text{
						Text: []string{request.String()},
					},
				},
			},
		},
	}
	return response, nil
}

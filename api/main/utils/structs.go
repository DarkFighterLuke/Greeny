package utils

type Intent struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

type QueryResult struct {
	Intent     Intent                       `json:"intent"`
	Parameters map[string]map[string]string `json:"parameters"`
}

type Text struct {
	Text []string `json:"text"`
}

type Message struct {
	Text Text `json:"text"`
}

// webhookRequest is used to unmarshal a WebhookRequest JSON object.
type WebhookRequest struct {
	Session     string      `json:"session"`
	ResponseID  string      `json:"responseId"`
	QueryResult QueryResult `json:"queryResult"`
}

// webhookResponse is used to marshal a WebhookResponse JSON object.
type WebhookResponse struct {
	FulfillmentMessages []Message `json:"fulfillmentMessages"`
}

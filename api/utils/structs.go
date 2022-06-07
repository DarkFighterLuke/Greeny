package utils

type Intent struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

type QueryResult struct {
	Intent     Intent                 `json:"intent"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Text struct {
	Text []string `json:"text"`
}

type Message struct {
	Text Text `json:"text"`
}

type Context struct {
	Name          string            `json:"name"`
	LifespanCount int               `json:"lifespanCount"`
	Parameters    map[string]string `json:"parameters"`
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
	OutputContext       []Context `json:"outputContexts"`
}

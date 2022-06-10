package utils

import "fmt"

const (
	ProjectId    = "greeny-ccqg"
	ContextsBase = "%s/contexts/%s"
)

type Intent struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

type QueryResult struct {
	Intent         Intent                 `json:"intent"`
	Parameters     map[string]interface{} `json:"parameters"`
	OutputContexts []Context              `json:"outputContexts"`
}

type Text struct {
	Text []string `json:"text"`
}

type Message struct {
	Text Text `json:"text"`
}

type Context struct {
	Name          string                 `json:"name"`
	LifespanCount int                    `json:"lifespanCount"`
	Parameters    map[string]interface{} `json:"parameters"`
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
	OutputContexts      []Context `json:"outputContexts"`
}

func FindContextByName(contexts *[]Context, session, contextName string) (Context, error) {
	for _, context := range *contexts {
		temp := fmt.Sprintf(ContextsBase, session, contextName)
		if context.Name == temp {
			return context, nil
		}
	}
	return Context{}, fmt.Errorf("no context found with the given name")
}

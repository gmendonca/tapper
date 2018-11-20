package datadog

import (
	datadogClient "gopkg.in/zorkian/go-datadog-api.v2"
)

// Datadog Struct that holds APIKey and ApplicationID
type Datadog struct {
	APIKey        string
	ApplicationID string
}

// GetClient a DataDog Client that make available to users to send metrics or interact with Datadog API
func (datadog *Datadog) GetClient() *datadogClient.Client {
	client := datadogClient.NewClient(datadog.APIKey, datadog.ApplicationID)
	return client
}

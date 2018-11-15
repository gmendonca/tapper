package datadog

import (
	datadogClient "gopkg.in/zorkian/go-datadog-api.v2"
)

type Datadog struct {
	APIKey        string
	ApplicationID string
}

func (datadog *Datadog) GetClient() *datadogClient.Client {
	client := datadogClient.NewClient(datadog.APIKey, datadog.ApplicationID)
	return client
}

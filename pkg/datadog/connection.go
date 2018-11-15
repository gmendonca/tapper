package datadog

import (
	datadogClient "gopkg.in/zorkian/go-datadog-api.v2"
)

type Datadog struct {
	ApiKey        string
	ApplicationId string
}

func (datadog *Datadog) GetClient() *datadogClient.Client {
	client := datadogClient.NewClient(datadog.ApiKey, datadog.ApplicationId)
	return client
}

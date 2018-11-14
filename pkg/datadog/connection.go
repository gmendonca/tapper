package datadog

import (
	"fmt"
	"strconv"
)

type Datadog struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

func (datadog *Datadog) PostMetric)

func (datadog *Datadog) GetURL() string {
	protocol := "http"
	if datadog.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, datadog.Host, strconv.Itoa(datadog.Port))
}

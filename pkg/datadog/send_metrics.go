package datadog

import (
	"fmt"
	"time"

	datadogClient "gopkg.in/zorkian/go-datadog-api.v2"
)

//
func (datadog *Datadog) PostMetrics(metricName string, point float64, hostname string, metricType string, tags []string) (bool, error) {
	dataPoint := datadogClient.DataPoint{float64(time.Now().Unix()), point}
	dataPoints := []datadogClient.DataPoint{dataPoint}

	fmt.Println(metricName)
	m := &datadogClient.Metric{
		Metric: &metricName,
		Points: dataPoints,
		Type:   &metricType,
		Host:   &hostname,
		Tags:   tags,
	}

	series := []datadogClient.Metric{*m}

	err := datadog.GetClient().PostMetrics(series)
	if err != nil {
		return false, err
	}

	return true, nil
}

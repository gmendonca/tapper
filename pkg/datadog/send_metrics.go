package datadog

func (datadog *Datadog) PostMetrics(points []float64, hostname string, metricType string, tags []string) string {
	m := &datadogClient.Metrics{
		Metric: "",
		Points: points,
		Type:   metricType,
		Host:   hostname,
		Tags:   tags,
	}
}

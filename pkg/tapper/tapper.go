package tapper

import (
	"github.com/gmendonca/tapper/pkg/datadog"
	"github.com/gmendonca/tapper/pkg/elasticsearch"
	"github.com/gmendonca/tapper/pkg/query"
)

// Tapper is your friendly Log and Metrics Analyzer
type Tapper struct {
	Elasticsearch *elasticsearch.Elasticsearch
	Datadog       *datadog.Datadog
	Dogstatsd     *datadog.Dogstatsd
}

// SendMetrics allow Tapper to get logs from ES and send them to Datadog or Dogstatsd
func (tapper *Tapper) SendMetrics(queryType string) {
	query.SendQueries(tapper.Datadog, tapper.Dogstatsd, tapper.Elasticsearch, queryType)
}

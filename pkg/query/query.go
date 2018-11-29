package query

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gmendonca/tapper/pkg/datadog"
	"github.com/gmendonca/tapper/pkg/elasticsearch"
	"github.com/gmendonca/tapper/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	queryFilter       = "timestamp"
	metricNamePattern = "%s.query.count"
	metricType        = "gauge"
)

//Point is the Queries Count Metrics struct
type Point struct {
	Hostname string
	Count    float64
	Source   string
}

// GetQueries go through the ES Index looking for indices with format hive-<hostname>-2018.11.21
// and then get some information of the results. Right now is getting the records of the last five minutes
func GetQueries(dogstatsd *datadog.Dogstatsd, elasticsearch *elasticsearch.Elasticsearch, queryType string) []Point {
	now := time.Now().Format(time.RFC3339)
	count := 5
	from := time.Now().Add(time.Duration(-count) * time.Minute).Format(time.RFC3339)

	query := elastic.NewRangeQuery(queryFilter)
	query.Gte(from)
	query.Lt(now)

	src, err := query.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	indices := elasticsearch.GetIndicesNames(queryType)

	var queryPoints []Point

	ctx := context.Background()
	client := elasticsearch.GetClient()
	for _, index := range indices {
		scroll := client.Scroll(index.IndexName).Query(query).Pretty(true).Size(100)

		count := 0

		for {
			searchResult, err := scroll.Do(ctx)
			if err == io.EOF {
				break
			}

			log.Debug("Found a total of %d queries\n", searchResult.TotalHits())

			if err != nil {
				panic(err)
			}

			log.Debug("Query took %d milliseconds\n", searchResult.TookInMillis)

			var tags []string

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				if queryType == "hive" {

					var hiveQuery model.HiveQuery
					err := json.Unmarshal(*hit.Source, &hiveQuery)
					if err != nil {
						// Deserialization failed
						panic(err)
					}

					tags = append(tags, "queue:"+hiveQuery.Queue)
					tags = append(tags, "user:"+hiveQuery.User)
					tags = append(tags, "host:"+index.Hostname)
					tags = append(tags, "application:"+queryType)

					duration := cast.ToDuration(hiveQuery.TimeTaken)
					dogstatsd.SendTiming(queryType, "query.duration", tags, duration)
					compileTime := cast.ToDuration(hiveQuery.CompileTime)
					dogstatsd.SendTiming(queryType, "query.compile_time", tags, compileTime)

				} else if queryType == "presto" {

					var prestoQuery model.PrestoQuery
					err := json.Unmarshal(*hit.Source, &prestoQuery)
					if err != nil {
						// Deserialization failed
						panic(err)
					}

					tags = append(tags, "memory_pool:"+prestoQuery.MemoryPool)
					tags = append(tags, "source:"+prestoQuery.Source)
					tags = append(tags, "state:"+prestoQuery.State)
					tags = append(tags, "user:"+prestoQuery.User)
					tags = append(tags, "user_agent:"+prestoQuery.UserAgent)
					tags = append(tags, "host:"+index.Hostname)
					tags = append(tags, "application:"+queryType)

					duration := cast.ToDuration(prestoQuery.TimeTaken)
					dogstatsd.SendTiming(queryType, "query.duration", tags, duration)
				}

				dogstatsd.SendGauge(queryType, "query.point", tags, float64(1))
				tags = []string{}
				count++
			}
		}

		queryPoint := Point{
			Hostname: index.Hostname,
			Count:    float64(count),
			Source:   queryType,
		}
		queryPoints = append(queryPoints, queryPoint)
	}

	return queryPoints
}

//SendQueries send the info from queries as Datadog metrics
func SendQueries(datadog *datadog.Datadog, dogstatsd *datadog.Dogstatsd, elasticsearch *elasticsearch.Elasticsearch, queryType string) {
	queryPoints := GetQueries(dogstatsd, elasticsearch, queryType)

	metricName := fmt.Sprintf(metricNamePattern, queryType)

	for _, queryPoint := range queryPoints {
		datadog.PostMetrics(metricName, queryPoint.Count, queryPoint.Hostname, metricType, nil)
	}
}

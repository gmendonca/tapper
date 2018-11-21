package elasticsearch

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

const (
	timestamp = "@timestamp"
)

func (elasticsearch *Elasticsearch) GetQueries() {
	now := time.Now().Format(time.RFC3339)
	count := 5
	from := time.Now().Add(time.Duration(-count) * time.Minute).Format(time.RFC3339)

	query := elastic.NewRangeQuery(timestamp)
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
}

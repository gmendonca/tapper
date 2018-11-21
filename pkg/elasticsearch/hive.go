package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

const (
	timestamp = "timestamp"
)

func inTimeSpan(start, end, check time.Time) bool {
	// Get Dates after the start or equal to end
	return check.After(start) || check == end
}

func (elasticsearch *Elasticsearch) getHiveIndicesNames() []string {
	client := elasticsearch.getClient()
	names, err := client.IndexNames()
	if err != nil {
		// Handle error
		panic(err)
	}

	r, _ := regexp.Compile("hive-[a-z0-9\\.\\-]+-(?P<year>\\d{4}).(?P<month>\\d{2}).(?P<day>\\d{2})")

	var indices []string

	d := 24 * time.Hour
	now := time.Now().UTC().Truncate(d)
	count := 5
	from := time.Now().UTC().Add(time.Duration(-count) * time.Minute).Truncate(d)

	for _, name := range names {
		if r.MatchString(name) {
			match := r.FindStringSubmatch(name)
			result := make(map[string]string)
			for i, name := range r.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			date := fmt.Sprintf("%s-%s-%s", result["year"], result["month"], result["day"])
			check, err := time.Parse("2006-01-02", date)
			if err != nil {
				panic(err)
			}
			if inTimeSpan(from, now, check) {
				indices = append(indices, name)
				fmt.Println(name)
			}
		}
	}

	return indices
}

// GetQueries go through the ES Index looking for indices with format hive-<hostname>-2018.11.21
// and then get some information of the results. Right now is getting the records of the last five minutes
func (elasticsearch *Elasticsearch) GetQueries() {
	now := time.Now().Format(time.RFC3339)
	count := 30
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

	indices := elasticsearch.getHiveIndicesNames()

	ctx := context.Background()
	client := elasticsearch.getClient()
	for _, index := range indices {
		searchResult, err := client.Search().
			Index(index).
			Query(query).
			Sort("timestamp", true).
			Pretty(true).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

		count := 0

		var hiveQuery HiveQuery
		for _, item := range searchResult.Each(reflect.TypeOf(hiveQuery)) {
			t := item.(HiveQuery)
			fmt.Printf("HiveQuery by %s on %s\n", t.User, t.Timestamp)
			count = count + 1
		}

		fmt.Printf("Count = %d\n", count)

		// TotalHits is another convenience function that works even when something goes wrong.
		fmt.Printf("Found a total of %d queries\n", searchResult.TotalHits())
	}
}

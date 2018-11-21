package elasticsearch

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

const (
	timestamp = "@timestamp"
	indexName = ""
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
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

	now := time.Now()
	count := 5
	from := time.Now().Add(time.Duration(-count) * time.Minute)

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
			check, _ := time.Parse(time.RFC3339, date)
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

	elasticsearch.getHiveIndicesNames()

	// ctx := context.Background()
	// client := elasticsearch.getClient()
	// searchResult, err := client.Search().
	// Index("twitter").         // search in index "twitter"
	// Query(query).             // specify the query
	// Sort("@timestamp", true). // sort by "user" field, ascending
	// Pretty(true).             // pretty print request and response JSON
	// Do(ctx)                   // execute
	//
	// if err != nil {
	// // Handle error
	// panic(err)
	// }
	//
	// // searchResult is of type SearchResult and returns hits, suggestions,
	// // and all kinds of other information from Elasticsearch.
	// fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	//
	// var ttyp Tweet
	// for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// t := item.(Tweet)
	// fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	// }
	//
	// // TotalHits is another convenience function that works even when something goes wrong.
	// fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())
}

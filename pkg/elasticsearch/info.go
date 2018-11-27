package elasticsearch

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gmendonca/tapper/pkg/datadog"
	log "github.com/sirupsen/logrus"
)

type IndexPair struct {
	Hostname  string
	IndexName string
}

func inTimeSpan(start, end, check time.Time) bool {
	// Get Dates after the start or equal to end
	return check.After(start) || check == end
}

func (elasticsearch *Elasticsearch) GetIndicesNames(prefix string) []IndexPair {
	client := elasticsearch.GetClient()
	names, err := client.IndexNames()
	if err != nil {
		// Handle error
		panic(err)
	}

	r, _ := regexp.Compile(fmt.Sprintf("%s-(?P<hostname>[a-z0-9\\.\\-]+)-(?P<year>\\d{4}).(?P<month>\\d{2}).(?P<day>\\d{2})", prefix))

	var indices []IndexPair

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
				indexPair := IndexPair{
					Hostname:  result["hostname"],
					IndexName: name,
				}
				indices = append(indices, indexPair)
				log.Debug(name)
			}
		}
	}

	return indices
}

func (elasticsearch *Elasticsearch) analyzeQueries(dogstatsd *datadog.Dogstatsd, query string) []string {
	parts := strings.Fields(query)

	return parts
}

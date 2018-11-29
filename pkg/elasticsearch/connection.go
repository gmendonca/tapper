package elasticsearch

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

//Elasticsearch is the information to access the Elasticsearch host
type Elasticsearch struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

//GetURL uses Elasticsearch struct to build up the url to access it
func (elasticsearch *Elasticsearch) GetURL() string {
	protocol := "http"
	if elasticsearch.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, elasticsearch.Host, strconv.Itoa(elasticsearch.Port))
}

// GetClient returns an Elasticsearch client to interact with
func (elasticsearch *Elasticsearch) GetClient() *elastic.Client {
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticsearch.GetURL()))
	if err != nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("Connected to Elasticsearch at %s", elasticsearch.GetURL()))
	return client
}

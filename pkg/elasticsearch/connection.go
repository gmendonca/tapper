package elasticsearch

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

type Elasticsearch struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

func (elasticsearch *Elasticsearch) getURL() string {
	protocol := "http"
	if elasticsearch.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, elasticsearch.Host, strconv.Itoa(elasticsearch.Port))
}

func (elasticsearch *Elasticsearch) GetClient() *elastic.Client {
	client, err := elastic.NewSimpleClient(elastic.SetURL(elasticsearch.getURL()))
	if err != nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("Connected to Elasticsearch at %s", elasticsearch.getURL()))
	return client
}

package elasticsearch

import (
	"fmt"
	"strconv"

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

func (elasticsearch *Elasticsearch) GetClient() error {
	e := elastic.NewClient()
	e.SetURL(elasticsearch.getURL())
}

package elasticsearch

import (
	"fmt"
	"strconv"
)

type Elasticsearch struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

func (elasticsearch *Elasticsearch) GetURL() string {
	protocol := "http"
	if elasticsearch.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, elasticsearch.Host, strconv.Itoa(elasticsearch.Port))
}

package cloudera

import (
	"fmt"
	"strconv"
)

type Cloudera struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

func (cloudera *Cloudera) GetURL() string {
	protocol := "http"
	if cloudera.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, cloudera.Host, strconv.Itoa(cloudera.Port))
}

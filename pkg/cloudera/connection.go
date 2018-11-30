package cloudera

import (
	"fmt"
	"strconv"
)

//Cloudera holds information to access the Cloudera Manager API
type Cloudera struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

//GetURL uses Cloudera struct to build up the url to access it
func (cloudera *Cloudera) GetURL() string {
	protocol := "http"
	if cloudera.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, cloudera.Host, strconv.Itoa(cloudera.Port))
}

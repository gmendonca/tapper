package datadog

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

// Dogstatsd Struct that holds Host and Port
type Dogstatsd struct {
	Host string
	Port int
}

// GetClient provides a clinet for dogstastd interface
func (dogstatsd *Dogstatsd) getClient() *statsd.Client {
	address := fmt.Sprintf("%s:%d", dogstatsd.Host, dogstatsd.Port)
	c, err := statsd.New(address)
	if err != nil {
		panic(err)
	}
	return c
}

func (dogstatsd *Dogstatsd) SendGauge(namespace string, name string, tags []string, value float64) bool {
	c := dogstatsd.getClient()

	c.Namespace = namespace
	c.Tags = tags

	err := c.Gauge(name, value, tags, 1)

	if err != nil {
		return false
	}
	return true
}

func (dogstatsd *Dogstatsd) SendCounter(namespace string, name string, tags []string) bool {
	c := dogstatsd.getClient()

	c.Namespace = namespace
	c.Tags = tags

	err := c.Incr(name, tags, 1)

	if err != nil {
		return false
	}
	return true
}

package datadog

import (
	"fmt"
	"time"

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

// SendGauge sends a Gauge Metric to Dogstatsd
func (dogstatsd *Dogstatsd) SendGauge(namespace string, name string, tags []string, value float64) bool {
	c := dogstatsd.getClient()

	c.Namespace = fmt.Sprintf("%s.", namespace)
	c.Tags = tags

	err := c.Gauge(name, value, tags, 1)

	if err != nil {
		return false
	}
	return true
}

// SendTiming sends a Timing Metric to Dogstatsd
func (dogstatsd *Dogstatsd) SendTiming(namespace string, name string, tags []string, duration time.Duration) bool {
	c := dogstatsd.getClient()

	c.Namespace = fmt.Sprintf("%s.", namespace)
	c.Tags = tags

	err := c.Timing(name, duration, tags, 1)

	if err != nil {
		return false
	}
	return true
}

// SendCounter sends a Counter Increment Metric to Dogstatsd
func (dogstatsd *Dogstatsd) SendCounter(namespace string, name string, tags []string) bool {
	c := dogstatsd.getClient()

	c.Namespace = fmt.Sprintf("%s.", namespace)
	c.Tags = tags

	err := c.Incr(name, tags, 1)

	if err != nil {
		return false
	}
	return true
}

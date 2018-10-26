package main

import (
	"fmt"

	"github.com/gmendonca/query-metrics-go/pkg/cloudera"
)

func main() {
	c := &cloudera.Cloudera{
		Host:     "cloudera.test.com",
		Port:     7180,
		Username: "username",
		Password: "password",
		SSL:      true,
	}

	fmt.Println(c.GetURL())
}

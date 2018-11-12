package main

import (
	"fmt"

	"github.com/gmendonca/query-metrics-go/pkg/cloudera"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	c := &cloudera.Cloudera{
		Host:     viper.GetString("CLOUDERA_HOST"),
		Port:     viper.GetInt("CLOUDERA_PORT"),
		Username: viper.GetString("CLOUDERA_USERNAME"),
		Password: viper.GetString("CLOUDERA_PASSWORD"),
		SSL:      viper.GetBool("CLOUDERA_SSL"),
	}

	fmt.Println(c.GetURL())
}

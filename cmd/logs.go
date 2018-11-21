package cmd

import (
	"fmt"

	"github.com/gmendonca/tapper/pkg/elasticsearch"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "logs command line interface",
	Long:  `logs, the command line interface`,
	Run: func(cmd *cobra.Command, args []string) {
		e := &elasticsearch.Elasticsearch{
			Host:     viper.GetString("elasticsearch.host"),
			Port:     viper.GetInt("elasticsearch.port"),
			Username: viper.GetString("elasticsearch.username"),
			Password: viper.GetString("elasticsearch.password"),
			SSL:      viper.GetBool("elasticsearch.ssl"),
		}

		log.Info(fmt.Sprintf("Connected to Elasticsearch at %s", e.GetURL()))

		e.GetQueries()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

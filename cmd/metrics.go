package cmd

import (
	"fmt"

	"github.com/gmendonca/tapper/pkg/cloudera"
	"github.com/gmendonca/tapper/pkg/datadog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "metrics command line interface",
	Long:  `metrics, the command line interface`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &cloudera.Cloudera{
			Host:     viper.GetString("cloudera.host"),
			Port:     viper.GetInt("cloudera.port"),
			Username: viper.GetString("cloudera.username"),
			Password: viper.GetString("cloudera.password"),
			SSL:      viper.GetBool("cloudera.ssl"),
		}

		log.Info(fmt.Sprintf("Connected to Cloudera at %s", c.GetURL()))

		d := &datadog.Datadog{
			APIKey:        viper.GetString("datadog.api_key"),
			ApplicationID: viper.GetString("datadog.application_id"),
		}

		c.SendHiveMetastoreOpenConnectionMetrics(d)

		c.SendHiveServerOpenConnectionMetrics(d)
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}

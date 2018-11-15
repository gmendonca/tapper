package cmd

import (
	"fmt"
	"os"

	"github.com/gmendonca/query-metrics-go/pkg/cloudera"
	"github.com/gmendonca/query-metrics-go/pkg/datadog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "something",
	Long:  `tl; dr;`,
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
			ApiKey:        viper.GetString("datadog.api_key"),
			ApplicationId: viper.GetString("datadog.application_id"),
		}

		c.SendHiveMetastoreOpenConnectionMetrics(d)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

package cloudera

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gmendonca/tapper/pkg/datadog"
	log "github.com/sirupsen/logrus"
)

type ClouderaTimeSeries struct {
	Items []struct {
		TimeSeries []struct {
			Metadata struct {
				MetricName string    `json:"metricName"`
				EntityName string    `json:"entityName"`
				StartTime  time.Time `json:"startTime"`
				EndTime    time.Time `json:"endTime"`
				Attributes struct {
					ClusterName        string `json:"clusterName"`
					RackID             string `json:"rackId"`
					RoleConfigGroup    string `json:"roleConfigGroup"`
					ClusterDisplayName string `json:"clusterDisplayName"`
					HostID             string `json:"hostId"`
					Hostname           string `json:"hostname"`
					RoleName           string `json:"roleName"`
					ServiceType        string `json:"serviceType"`
					EntityName         string `json:"entityName"`
					Version            string `json:"version"`
					ServiceName        string `json:"serviceName"`
					Category           string `json:"category"`
					RoleType           string `json:"roleType"`
					Active             string `json:"active"`
					ServiceDisplayName string `json:"serviceDisplayName"`
				} `json:"attributes"`
				UnitNumerators              []string      `json:"unitNumerators"`
				UnitDenominators            []interface{} `json:"unitDenominators"`
				Expression                  string        `json:"expression"`
				MetricCollectionFrequencyMs int           `json:"metricCollectionFrequencyMs"`
				RollupUsed                  string        `json:"rollupUsed"`
			} `json:"metadata"`
			Data []interface{} `json:"data"`
		} `json:"timeSeries"`
		Warnings        []interface{} `json:"warnings"`
		TimeSeriesQuery string        `json:"timeSeriesQuery"`
	} `json:"items"`
}

type ClouderaPoint struct {
	Point       float64
	Hostname    string
	ClusterName string
}

const HiveMetastore = "HIVEMETASTORE"
const HiveServer = "HIVESERVER2"

func (cloudera *Cloudera) GetHiveOpenConnectionMetrics(roleType string) []ClouderaPoint {
	now := time.Now().Format(time.RFC3339)
	count := 5
	from := time.Now().Add(time.Duration(-count) * time.Minute).Format(time.RFC3339)

	endpoint := "api/v18/timeseries"

	url := fmt.Sprintf("%s/%s?query=select+hive_open_connections+where+roleType%%3D%s&contentType=application%%2Fjson&from=%s&to=%s", cloudera.GetURL(), endpoint, roleType, from, now)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cloudera.Username, cloudera.Password)
	resp, err := client.Do(req)

	if err != nil {
		return []ClouderaPoint{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var clouderaTimeSeries ClouderaTimeSeries
	jsonErr := json.Unmarshal(body, &clouderaTimeSeries)

	if jsonErr != nil {
		return []ClouderaPoint{}
	}

	var points []float64
	var hostname string
	var clusterName string

	var clouderaPoints []ClouderaPoint
	for _, item := range clouderaTimeSeries.Items {
		for _, timeserie := range item.TimeSeries {
			for _, datai := range timeserie.Data {
				data, _ := datai.(map[string]interface{})
				points = append(points, data["value"].(float64))
			}
			hostname = timeserie.Metadata.Attributes.Hostname
			clusterName = timeserie.Metadata.Attributes.ClusterName

			if len(points) == 0 {
				// No data points
				continue
			}

			sum := float64(0)

			for _, point := range points {
				sum = sum + float64(point)
			}

			clouderaPoint := ClouderaPoint{
				Point:       sum / float64(len(points)),
				Hostname:    hostname,
				ClusterName: clusterName,
			}

			clouderaPoints = append(clouderaPoints, clouderaPoint)
			points = make([]float64, 0)
		}
	}

	return clouderaPoints
}

func (cloudera *Cloudera) SendHiveMetastoreOpenConnectionMetrics(datadog *datadog.Datadog) {
	clouderaPoints := cloudera.GetHiveOpenConnectionMetrics(HiveMetastore)
	metricName := "cloudera.hive.metastore.openconnections"
	metricType := "gauge"

	for _, clouderaPoint := range clouderaPoints {
		tags := []string{fmt.Sprintf("cluster:%s", clouderaPoint.ClusterName)}

		run, err := datadog.PostMetrics(metricName, clouderaPoint.Point, clouderaPoint.Hostname, metricType, tags)

		if run {
			log.Info(fmt.Sprintf("Metric %s %f posted for cluster %s", metricName, clouderaPoint.Point, clouderaPoint.ClusterName))
		} else {
			log.Error(fmt.Sprintf("Metric %s %f not posted for cluster %s", metricName, clouderaPoint.Point, clouderaPoint.ClusterName))
			log.Error(err)
		}
	}
}

func (cloudera *Cloudera) SendHiveServerOpenConnectionMetrics(datadog *datadog.Datadog) {
	clouderaPoints := cloudera.GetHiveOpenConnectionMetrics(HiveServer)
	metricName := "cloudera.hive.server.openconnections"
	metricType := "gauge"

	for _, clouderaPoint := range clouderaPoints {
		tags := []string{fmt.Sprintf("cluster:%s", clouderaPoint.ClusterName)}

		run, err := datadog.PostMetrics(metricName, clouderaPoint.Point, clouderaPoint.Hostname, metricType, tags)

		if run {
			log.Info(fmt.Sprintf("Metric %s %f posted for cluster %s", metricName, clouderaPoint.Point, clouderaPoint.ClusterName))
		} else {
			log.Error(fmt.Sprintf("Metric %s %f not posted for cluster %s", metricName, clouderaPoint.Point, clouderaPoint.ClusterName))
			log.Error(err)
		}
	}
}

package cloudera

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

func (cloudera *Cloudera) GetHiveMetastoreOpenConnectionMetrics() string {
	now := time.Now().Format(time.RFC3339)
	count := 15
	from := time.Now().Add(time.Duration(-count) * time.Minute).Format(time.RFC3339)

	endpoint := "api/v18/timeseries"

	url := fmt.Sprintf("%s/%s?query=select+hive_open_connections+where+roleType%%3DHIVESERVER2&contentType=application%%2Fjson&from=%s&to=%s", cloudera.GetURL(), endpoint, from, now)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cloudera.Username, cloudera.Password)
	resp, err := client.Do(req)

	if err != nil {
		return "err"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var clouderaTimeSeries ClouderaTimeSeries
	jsonErr := json.Unmarshal(body, &clouderaTimeSeries)

	if jsonErr != nil {
		return "err"
	}
	for _, item := range clouderaTimeSeries.Items {
		for _, timeserie := range item.TimeSeries {
			for _, datai := range timeserie.Data {
				data, _ := datai.(map[string]interface{})
				fmt.Println(data["value"])
			}
		}
	}
	return "ok"
}

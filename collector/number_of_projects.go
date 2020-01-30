package collector

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

const nopMetricName = "number_of_projects"

type Project struct {
	ProjectID         int64        `json:"project_id"`
	OwnerID           int64        `json:"owner_id"`
	Name              string       `json:"name"`
	CreationTime      string       `json:"creation_time"`
	UpdateTime        string       `json:"update_time"`
	Deleted           bool         `json:"deleted"`
	OwnerName         string       `json:"owner_name"`
	Togglable         bool         `json:"togglable"`
	CurrentUserRoleID int64        `json:"current_user_role_id"`
	RepoCount         int64        `json:"repo_count"`
	ChartCount        int64        `json:"chart_count"`
	Metadata          Metadata     `json:"metadata"`
	CveWhitelist      CveWhitelist `json:"cve_whitelist"`
}

type CveWhitelist struct {
	ID           int64       `json:"id"`
	ProjectID    int64       `json:"project_id"`
	Items        interface{} `json:"items"`
	CreationTime string      `json:"creation_time"`
	UpdateTime   string      `json:"update_time"`
}

type Metadata struct {
	Public string `json:"public"`
}

var numOfProjectsMetric = prometheus.NewDesc(
	nopMetricName,
	"Number of projects in Harbor Registry",
	nil,
	nil,
)

type NumOfProjects struct{
	Client http.Client
}

func NewNumOfProjectsScraper() NumOfProjects{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return NumOfProjects{Client: http.Client{
		Transport: tr,
	}}
}

func (h NumOfProjects) MetricName() string {
	return nopMetricName
}

func (h NumOfProjects) Update(ch chan<- prometheus.Metric) error {
	decodedToken := os.Getenv("HARBOR_TOKEN")
	domain := os.Getenv("REGISTRY_SERVER") + "/api/projects"
	encodedToken := base64.StdEncoding.EncodeToString([]byte(decodedToken))

	req, err := http.NewRequest("GET", domain, nil)
	req.Header.Add("authorization", fmt.Sprintf("Basic %s", encodedToken))
	req.Header.Add("accept", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		fmt.Println(fmt.Errorf("error making request: %s", domain))
		fmt.Println(err)
	}
	var projects []Project
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println(fmt.Errorf("error unmarshalling response"))
		fmt.Println(err)
	}
	ch <- prometheus.MustNewConstMetric(numOfProjectsMetric,
		prometheus.GaugeValue, float64(len(projects)))
	return nil
}

// Assert Interface
var _ Scraper = NumOfProjects{}

package collector

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

type HarborHealthDashboard struct {
	Client http.Client
}

func NewHarborHealthDashboardScraper() HarborHealthDashboard {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return HarborHealthDashboard{Client: http.Client{
		Transport: tr,
	}}
}

func (h HarborHealthDashboard) MetricName() string {
	return HealthMetricName
}

func (h HarborHealthDashboard) Update(ch chan<- prometheus.Metric) error {
	// TODO: Logging lager
	var isUp int
	domain := os.Getenv("REGISTRY_SERVER")
	if domain == "" {
		return errors.New("missing environment variable REGISTRY_SERVER")
	}

	resp, err := h.Client.Get(fmt.Sprintf("%s/api/health", domain))
	if err != nil || resp.StatusCode != 200 {
		isUp = 0
	} else {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return err
		}
		healthStatuses, err := unmarshalWelcome(body)
		if err != nil {
			return err
		}
		isUp = 1
		for _, component := range healthStatuses.Components {
			if component.Status != Healthy {
				isUp = 0
				break
			}
		}
	}
	ch <- prometheus.MustNewConstMetric(
		HarborHealthDashboardMetric,
		prometheus.GaugeValue,
		float64(isUp),
	)
	return nil
}

func unmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

type Welcome struct {
	Status     string      `json:"status"`
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Assert Interface
var _ Scraper = HarborHealthDashboard{}

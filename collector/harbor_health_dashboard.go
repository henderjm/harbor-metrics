package collector

import (
	"crypto/tls"
	"errors"
	"fmt"
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
	domain := os.Getenv("REGISTRY_SERVER")
	if domain == "" {
		return errors.New("missing environment variable REGISTRY_SERVER")
	}

	var isUp int
	resp, err := h.Client.Get(fmt.Sprintf(domain))
	if err != nil {
		// TODO: Logging lager
		isUp = 0
	} else if resp.StatusCode == 200 {
		isUp = 1
	}

	ch <- prometheus.MustNewConstMetric(HarborHealthDashboardMetric,
		prometheus.GaugeValue, float64(isUp))
	return nil
}

// Assert Interface
var _ Scraper = HarborHealthDashboard{}

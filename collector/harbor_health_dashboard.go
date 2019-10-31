package collector

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const metricName = "harbor_health_collector"

var harborHealthDashboardMetric = prometheus.NewDesc(
	"harbor_health_collector",
	"Indicates the health of the harbor frontend",
	nil,
	nil,
)

type HarborHealthDashboard struct{}

func (h HarborHealthDashboard) MetricName() string {
	return metricName
}

func (h HarborHealthDashboard) Update(ch chan<- prometheus.Metric) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
	}

	domain := "https://192.168.64.2:30003"
	var isUp int
	resp, err := client.Get(fmt.Sprintf(domain))
	if err != nil {
		fmt.Println("error")
		fmt.Printf("message: %s\n", err.Error())
		isUp = 0
	} else if resp.StatusCode == 200 {
		isUp = 1
	}
	fmt.Printf("status code: %d\n", resp.StatusCode)

	fmt.Println(fmt.Sprintf("**MARK**isUp: %d **", isUp))

	ch <- prometheus.MustNewConstMetric(harborHealthDashboardMetric,
		prometheus.GaugeValue, float64(isUp))
	return nil
}

// Assert Interface
var _ Scraper = HarborHealthDashboard{}

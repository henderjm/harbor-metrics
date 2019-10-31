package collector

import "github.com/prometheus/client_golang/prometheus"

type Scraper interface {
	Update(ch chan<- prometheus.Metric) error
	MetricName() string
}

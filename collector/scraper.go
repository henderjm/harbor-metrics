package collector

import "github.com/prometheus/client_golang/prometheus"

type Scraper interface {
	Update(ch chan<- prometheus.Metric) error
	MetricName() string
}

// Const
const Healthy = "healthy"
const HealthMetricName = "harbor_health_collector"

// prometheus.Desc variables
var HarborHealthDashboardMetric = prometheus.NewDesc(
	HealthMetricName,
	"Indicates the health of the harbor frontend",
	nil,
	nil,
)

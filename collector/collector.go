package collector

import (
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"sync"
	"time"
)

type Value interface {
	Type() string
}

// Gauge represents a gauge metric value, such as a temperature.
// This is Go's equivalent to the C type "gauge_t".
type Gauge float64

// Type returns "gauge".
func (v Gauge) Type() string { return "gauge" }

type metrics map[string]*prometheus.Desc

type ValueList struct {
	Time     time.Time
	Interval time.Duration
	Values   []Value
	DSNames  []string
}

type harborCollector struct {
	valueLists map[string]*prometheus.Desc
	exitChan   <-chan struct{}
	mux        sync.Mutex
}

var (
	harborMetrics = metrics{
		"harbor_alive": newHarborMetric("harbor_health_collector", "Indicates the health of the harbor frontend", nil),
		"project_count": newHarborMetric("num_of_projects", "Counts the total number of projects in the Harbor Registry", nil),
	}
)

func newHarborMetric(metricName, helper string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(metricName, helper, nil, nil)
}

func NewHarborCollector() *harborCollector {
	fmt.Println("**Initialising Harbor Collector**")
	c := &harborCollector{
		exitChan:   make(chan struct{}),
		valueLists: harborMetrics,
	}

	return c
}

// Collect implements prometheus.Collector.
func (c harborCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("**COLLECTING**")

	c.mux.Lock() // To protect metrics from concurrent collects
	defer c.mux.Unlock()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport:     tr,
	}

	domain := "https://192.168.99.100:30003"
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

	ch <- prometheus.MustNewConstMetric(c.valueLists["harbor_alive"],
		prometheus.GaugeValue, float64(isUp))
}

// Describe implements prometheus.Collector.
func (c harborCollector) Describe(ch chan<- *prometheus.Desc) {
	fmt.Println("**DESCRIBING**")
	for _, m := range harborMetrics{
		ch <- m
	}
}

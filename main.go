package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Value interface {
	Type() string
}

// Gauge represents a gauge metric value, such as a temperature.
// This is Go's equivalent to the C type "gauge_t".
type Gauge float64

// Type returns "gauge".
func (v Gauge) Type() string { return "gauge" }

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

var metricDesc *prometheus.Desc

func init() {
	metricDesc = prometheus.NewDesc("harbor_health_collector", "Indicates the health of the harbor frontend", nil, nil)
}

func newHarborCollector() *harborCollector {
	c := &harborCollector{
		exitChan:   make(chan struct{}),
		valueLists: map[string]*prometheus.Desc{"ha": prometheus.NewDesc("harbor_health_collector", "Indicates the health of the harbor frontend", nil, nil)},
	}

	return c
}

// Collect implements prometheus.Collector.
func (c harborCollector) Collect(ch chan<- prometheus.Metric) {

	c.mux.Lock() // To protect metrics from concurrent collects
	defer c.mux.Unlock()

	domain := "reg.mydomain.io"
	client := http.DefaultClient
	var isUp int
	resp, err := client.Get(fmt.Sprintf("http://%s", domain))
	if err != nil {
		isUp = 0
	}
	if resp.StatusCode == 200 {
		isUp = 1
	}

	ch <- prometheus.MustNewConstMetric(c.valueLists["ha"],
		prometheus.GaugeValue, float64(isUp))
}

// Describe implements prometheus.Collector.
func (c harborCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricDesc
}

func main() {
	fmt.Println("starting harbor metrics collector")

	registry := prometheus.NewRegistry()
	registry.MustRegister(newHarborCollector())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		log.Printf("SIGTERM received: %v. Exiting...", <-signalChan)
		os.Exit(0)
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>Harbor Exporter</title></head>
			<body>
			<h1>Harbor Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Printf("Error while sending a response for the '/' path: %v", err)
		}
	})

	log.Printf("Harbor Prometheus Exporter has successfully started")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

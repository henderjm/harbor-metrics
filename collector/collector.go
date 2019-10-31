package collector

import (
	"fmt"
	"sync"
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

//type metrics map[int]*prometheus.Desc

type ValueList struct {
	Time     time.Time
	Interval time.Duration
	Values   []Value
	DSNames  []string
}

type harborCollector struct {
	valueLists map[int]*prometheus.Desc
	scrapers   []Scraper
	exitChan   <-chan struct{}
	mux        sync.Mutex
}

func NewHarborCollector(scrapers []Scraper) *harborCollector {
	c := &harborCollector{
		exitChan: make(chan struct{}),
		scrapers: scrapers,
	}

	return c
}

func (c harborCollector) scrape(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, scraper := range c.scrapers {
		wg.Add(1)
		go func(scraper Scraper) {
			defer wg.Done()
			err := scraper.Update(ch)
			if err != nil {
				fmt.Println(fmt.Errorf("error scraping: %s\n", scraper.MetricName()))
				fmt.Println(err)
			} else {
				fmt.Printf("successfully scraped: %s\n", scraper.MetricName())
			}
		}(scraper)
	}
	wg.Wait()
}

// Collect implements prometheus.Collector.
func (c harborCollector) Collect(ch chan<- prometheus.Metric) {
	c.scrape(ch)
}

// Describe implements prometheus.Collector.
func (c harborCollector) Describe(ch chan<- *prometheus.Desc) {
}

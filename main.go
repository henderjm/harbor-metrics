package main

import (
	"fmt"
	"henderjm/harbor-metrics/collector"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
)

var scrapers = []collector.Scraper{
	collector.HarborHealthDashboard{},
	collector.NumOfProjects{},
}

func scraperHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.NewHarborCollector(scrapers))
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}

		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func main() {
	fmt.Println("starting harbor metrics collector")

	handler := scraperHandler()
	http.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handler))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		log.Printf("SIGTERM received: %v. Exiting...", <-signalChan)
		os.Exit(0)
	}()

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

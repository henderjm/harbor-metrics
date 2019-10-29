package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"henderjm/harbor-metrics/collector"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
)


func main() {
	fmt.Println("starting harbor metrics collector")

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewHarborCollector())

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

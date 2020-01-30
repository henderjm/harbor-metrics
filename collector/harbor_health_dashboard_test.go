package collector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collector "henderjm/harbor-metrics/collector"
)

var _ = Describe("HarborHealthDashboard", func() {

	var healthScraper collector.Scraper = collector.HarborHealthDashboard{}

	It("Should have a metric name", func() {
		Expect(healthScraper.MetricName()).To(Equal("harbor_health_collector"))
	})
})

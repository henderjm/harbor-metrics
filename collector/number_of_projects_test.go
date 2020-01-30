package collector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"henderjm/harbor-metrics/collector"
)

var _ = Describe("NumberOfProjects", func() {

	var nopScraper collector.Scraper = collector.NewNumOfProjectsScraper()

	It("Should have a metric name", func() {
		Expect(nopScraper.MetricName()).To(Equal("number_of_projects"))
	})
})

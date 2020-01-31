package collector_test

import (
	"os"

	collector "henderjm/harbor-metrics/collector"

	"github.com/dankinder/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("HarborHealthDashboard", func() {

	var (
		healthScraper collector.Scraper = collector.NewHarborHealthDashboardScraper()
		s             *httpmock.Server
	)

	It("Should have a metric name", func() {
		Expect(healthScraper.MetricName()).To(Equal("harbor_health_collector"))
	})

	Context("When updating metrics", func() {
		Context("When Harbor is up", func() {
			BeforeEach(func() {
				downstream := &httpmock.MockHandler{}
				downstream.On("Handle", "GET", "/", mock.Anything).Return(httpmock.Response{
					Body:   []byte(`{"status": "healthy"}`),
					Status: 200,
				})

				s = httpmock.NewServer(downstream)
				err := os.Setenv("REGISTRY_SERVER", s.URL())
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should send a metric of value 1", func() {
				ch := make(chan prometheus.Metric, 1)
				err := healthScraper.Update(ch)
				Expect(err).ToNot(HaveOccurred())
				expectedResult := prometheus.MustNewConstMetric(
					collector.HarborHealthDashboardMetric,
					prometheus.GaugeValue,
					float64(1),
				)
				Eventually(ch).Should(Receive(Equal(expectedResult)))
				s.Close()
			})
		})

		Context("When Harbor is down", func() {
			BeforeEach(func() {
				downstream := &httpmock.MockHandler{}
				downstream.On("Handle", "GET", "/", mock.Anything).Return(httpmock.Response{
					Body:   []byte(`{"status": "healthy"}`),
					Status: 500,
				})

				s = httpmock.NewServer(downstream)
				err := os.Setenv("REGISTRY_SERVER", s.URL())
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should send a metric of value 0", func() {
				ch := make(chan prometheus.Metric, 1)
				err := healthScraper.Update(ch)
				Expect(err).ToNot(HaveOccurred())
				expectedResult := prometheus.MustNewConstMetric(
					collector.HarborHealthDashboardMetric,
					prometheus.GaugeValue,
					float64(0),
				)
				Eventually(ch).Should(Receive(Equal(expectedResult)))
				s.Close()
			})
		})
	})
})

package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// LinksCounter is a metrics to count the number of crawled links.
	LinksCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "crawler_links_count_total",
		Help: "Count of links returned",
	})
	// LinksErrorCounter is a metrics to count the number of error to craw links.
	LinksErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "crawler_links_error_count_total",
		Help: "Count of links returned in error",
	})
)

func init() {
	prometheus.MustRegister(LinksCounter)
	prometheus.MustRegister(LinksErrorCounter)
}

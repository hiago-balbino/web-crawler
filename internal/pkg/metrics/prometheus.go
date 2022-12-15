package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	linearBucketStart = 0.0001
	linearBucketWidth = 0.05
	linearBucketCount = 20
)

var (
	// LinksCounter is a metric to count the number of crawled links.
	LinksCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "crawler_links_count_total",
		Help: "Count of links returned",
	})
	// LinksErrorCounter is a metric to count the number of error to craw links.
	LinksErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "crawler_links_error_count_total",
		Help: "Count of links returned in error",
	})
	// DeltaTimeToProcessLinks is a metric used to measure the delta time to process the links
	DeltaTimeToProcessLinks = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "crawler_delta_time_to_process_links",
		Help:    "Delta time to process links",
		Buckets: prometheus.LinearBuckets(linearBucketStart, linearBucketWidth, linearBucketCount),
	})
)

func init() {
	prometheus.MustRegister(LinksCounter)
	prometheus.MustRegister(LinksErrorCounter)
	prometheus.MustRegister(DeltaTimeToProcessLinks)
}

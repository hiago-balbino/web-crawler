package crawler

// CrawlerService is a abstraction to handle with web crawler
type CrawlerService interface {
	// Craw execute the call to craw pages concurrently and will respect depth param
	Craw(uri string, depth int) ([]string, error)
}

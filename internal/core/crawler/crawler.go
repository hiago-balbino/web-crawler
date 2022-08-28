package crawler

// Crawler is a abstraction to handle with web crawler
type Crawler interface {
	// Craw execute the call to craw pages concurrently and will respect depth param
	Craw(uri string, depth int) ([]string, error)
}

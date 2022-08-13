package crawler

// Crawler is a abstraction to handle with web crawler
type Crawler interface {
	// Craw execute the call to craw pages
	Craw(uri string)
}

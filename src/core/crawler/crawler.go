package crawler

// Crawler is a abstraction to handle with web crawler
type Crawler interface {
	// Run execute the call to craw pages
	Run(uri string)
}

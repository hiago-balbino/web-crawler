package crawler

import "github.com/stretchr/testify/mock"

// CrawlerMock is a mock to abstraction that handle with web crawler
type CrawlerMock struct {
	mock.Mock
}

// Run execute the call to craw pages
func (c *CrawlerMock) Run(uri string) {
	c.Called(uri)
}

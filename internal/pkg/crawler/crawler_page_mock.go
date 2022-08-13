package crawler

import "github.com/stretchr/testify/mock"

// CrawlerPageMock is a mock to abstraction that handle with web crawler
type CrawlerPageMock struct {
	mock.Mock
}

// Craw execute the call to craw pages
func (c *CrawlerPageMock) Craw(uri string) {
	c.Called(uri)
}

package crawler

import "github.com/stretchr/testify/mock"

// CrawlerServiceMock is a mock to abstraction that handle with web crawler.
type CrawlerServiceMock struct {
	mock.Mock
}

// Craw execute the call to craw pages concurrently and will respect depth param.
func (c *CrawlerServiceMock) Craw(uri string, depth int) ([]string, error) {
	args := c.Called(uri, depth)

	return args.Get(0).([]string), args.Error(1)
}

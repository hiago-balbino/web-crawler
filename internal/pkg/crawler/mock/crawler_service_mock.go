package crawler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// CrawlerServiceMock is a mock to abstraction that handle with web crawler.
type CrawlerServiceMock struct {
	mock.Mock
}

// Craw execute the call to craw pages concurrently and will respect depth param.
func (c *CrawlerServiceMock) Craw(ctx context.Context, uri string, depth uint) ([]string, error) {
	args := c.Called(ctx, uri, depth)

	return args.Get(0).([]string), args.Error(1)
}

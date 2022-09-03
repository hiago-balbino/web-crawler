package crawler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// CrawlerDatabaseMock is a mock to repository that handle with persistence layer.
type CrawlerDatabaseMock struct {
	mock.Mock
}

// Find is a method to fetch links crawled from database.
func (c *CrawlerDatabaseMock) Find(ctx context.Context, uri string) ([]string, error) {
	args := c.Called(ctx, uri)

	return args.Get(0).([]string), args.Error(1)
}

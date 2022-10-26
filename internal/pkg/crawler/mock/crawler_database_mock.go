package crawler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// CrawlerDatabaseMock is a mock to repository that handle with persistence layer.
type CrawlerDatabaseMock struct {
	mock.Mock
}

// Insert is a method to insert new page crawled on database.
func (c *CrawlerDatabaseMock) Insert(ctx context.Context, uri string, depth uint, uris []string) error {
	args := c.Called(ctx, uri, depth, uris)

	return args.Error(0)
}

// Find is a method to fetch links crawled from database.
func (c *CrawlerDatabaseMock) Find(ctx context.Context, uri string, depth uint) ([]string, error) {
	args := c.Called(ctx, uri, depth)

	return args.Get(0).([]string), args.Error(1)
}

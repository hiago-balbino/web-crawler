package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CrawlerDatabaseMock struct {
	mock.Mock
}

func (c *CrawlerDatabaseMock) Insert(ctx context.Context, uri string, depth uint, uris []string) error {
	args := c.Called(ctx, uri, depth, uris)

	return args.Error(0)
}

func (c *CrawlerDatabaseMock) Find(ctx context.Context, uri string, depth uint) ([]string, error) {
	args := c.Called(ctx, uri, depth)

	return args.Get(0).([]string), args.Error(1)
}

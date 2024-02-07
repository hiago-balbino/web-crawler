package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CrawlerUsecaseMock struct {
	mock.Mock
}

func (c *CrawlerUsecaseMock) Craw(ctx context.Context, uri string, depth uint) ([]string, error) {
	args := c.Called(ctx, uri, depth)

	return args.Get(0).([]string), args.Error(1)
}

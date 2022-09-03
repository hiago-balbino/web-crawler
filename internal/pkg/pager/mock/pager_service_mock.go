package pager

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
)

// PagerServiceMock is a mock to abstraction that handle with pages.
type PagerServiceMock struct {
	mock.Mock
}

// GetNode fetch and parse response body to return HTML Node.
func (p *PagerServiceMock) GetNode(uri string) (*html.Node, error) {
	args := p.Called(uri)

	return args.Get(0).(*html.Node), args.Error(1)
}

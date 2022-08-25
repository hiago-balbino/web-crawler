package pager

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
)

// PageProviderMock is a mock to abstraction that handle with pages
type PageProviderMock struct {
	mock.Mock
}

// GetNode fetch and parse response body to return HTML Node
func (p *PageProviderMock) GetNode(uri string) (*html.Node, error) {
	args := p.Called(uri)
	return args.Get(0).(*html.Node), args.Error(1)
}

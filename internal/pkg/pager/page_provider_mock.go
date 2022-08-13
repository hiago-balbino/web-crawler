package pager

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// PageProviderMock is a mock to abstraction that handle with pages
type PageProviderMock struct {
	mock.Mock
}

// Get fetch the response body from the URI
func (p *PageProviderMock) Get(uri string) (*http.Response, error) {
	args := p.Called(uri)
	return args.Get(0).(*http.Response), args.Error(1)
}

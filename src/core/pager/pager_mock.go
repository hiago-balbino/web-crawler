package pager

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// PagerMock is a mock to abstraction that handle with pages
type PagerMock struct {
	mock.Mock
}

// Get fetch the response body from the URI
func (p *PagerMock) Get(uri string) (*http.Response, error) {
	args := p.Called(uri)
	return args.Get(0).(*http.Response), args.Error(1)
}

package pager

import (
	"net/http"

	"github.com/hiago-balbino/web-crawler/internal/core/pager"
)

// PageProvider is a structure with HTTP Client that will be instantiated once using package sync
type PageProvider struct {
	httpClient *http.Client
}

// NewPageProvider is a constructor to create a new instance of PageProvider
func NewPageProvider(httpClient *http.Client) pager.Pager {
	return PageProvider{httpClient: httpClient}
}

// Get fetch the response body from the URI
func (c PageProvider) Get(uri string) (*http.Response, error) {
	response, err := c.httpClient.Get(uri)
	defer func() {
		if err == nil {
			_ = response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	return response, nil
}

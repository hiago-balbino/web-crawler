package pager

import (
	"net/http"
	"sync"
)

var once sync.Once
var httpClient *http.Client

// ContentPage is a structure with HTTP Client that will be instantiated once using package sync
type ContentPage struct {
	httpClient *http.Client
}

// NewContentPage is a constructor to create a new instance of ContentPage
func NewContentPage() Pager {
	once.Do(func() {
		httpClient = &http.Client{}
	})
	return ContentPage{httpClient: httpClient}
}

// Get fetch the response body from the URI
func (c ContentPage) Get(uri string) (*http.Response, error) {
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

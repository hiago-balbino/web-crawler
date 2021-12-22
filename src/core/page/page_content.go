package page

import (
	"net/http"
	"sync"
)

var once sync.Once

// ContentPage is a structure with HTTP Client that will be instantiated once using package sync
type ContentPage struct {
	httpClient *http.Client
}

// NewContentPage is a constructor to create a new instance of ContentPage
func NewContentPage() Pager {
	contentPage := ContentPage{}
	once.Do(func() {
		contentPage.httpClient = &http.Client{}
	})
	return contentPage
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

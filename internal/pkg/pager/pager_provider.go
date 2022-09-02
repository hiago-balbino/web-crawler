package pager

import (
	"net/http"

	"golang.org/x/net/html"
)

// PagerProvider is a structure with HTTP Client to handle with pages
type PagerProvider struct {
	httpClient *http.Client
}

// NewPagerProvider is a constructor to create a new instance of PagerProvider
func NewPagerProvider(httpClient *http.Client) PagerProvider {
	return PagerProvider{httpClient: httpClient}
}

// GetNode fetch and parse response body to return HTML Node
func (c PagerProvider) GetNode(uri string) (*html.Node, error) {
	response, err := c.httpClient.Get(uri)
	defer func() {
		if err == nil {
			_ = response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	return node, nil
}

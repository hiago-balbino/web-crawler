package pager

import (
	"net/http"

	"github.com/hiago-balbino/web-crawler/internal/core/pager"
	"golang.org/x/net/html"
)

// PageProvider is a structure with HTTP Client that will be instantiated once using package sync
type PageProvider struct {
	httpClient *http.Client
}

// NewPageProvider is a constructor to create a new instance of PageProvider
func NewPageProvider(httpClient *http.Client) pager.Pager {
	return PageProvider{httpClient: httpClient}
}

// GetNode fetch and parse response body to return HTML Node
func (c PageProvider) GetNode(uri string) (*html.Node, error) {
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

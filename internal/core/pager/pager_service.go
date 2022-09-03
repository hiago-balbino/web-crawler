package pager

import (
	"golang.org/x/net/html"
)

// PagerService is a abstraction to handle with pages.
type PagerService interface {
	// GetNode fetch and parse response body to return HTML Node.
	GetNode(uri string) (*html.Node, error)
}

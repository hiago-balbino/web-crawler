package pager

import (
	"golang.org/x/net/html"
)

// Pager is a abstraction to handle with pages
type Pager interface {
	// GetNode fetch and parse response body to return HTML Node
	GetNode(uri string) (*html.Node, error)
}

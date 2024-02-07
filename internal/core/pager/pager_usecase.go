package pager

import (
	"golang.org/x/net/html"
)

type PagerUsecase interface {
	GetNode(uri string) (*html.Node, error)
}

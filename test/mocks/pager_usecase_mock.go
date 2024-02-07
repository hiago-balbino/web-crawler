package mocks

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
)

type PagerUsecaseMock struct {
	mock.Mock
}

func (p *PagerUsecaseMock) GetNode(uri string) (*html.Node, error) {
	args := p.Called(uri)

	return args.Get(0).(*html.Node), args.Error(1)
}

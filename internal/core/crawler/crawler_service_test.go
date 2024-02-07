package crawler

import (
	"context"
	"errors"
	"testing"

	"github.com/hiago-balbino/web-crawler/test/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestCrawlerService_Craw(t *testing.T) {
	ctx := context.Background()
	URI := "https://anyurl.com"
	internalURI := "https://internal-anyurl.com"
	randomInternalURI := "https://random-internal-anyurl.com"
	subInternalURI := "https://sub-internal-anyurl.com"
	lastInternalURI := "https://last-internal-anyurl.com"
	unexpectedErr := errors.New("unexpected error")

	testCases := map[string]func(*testing.T, *mocks.PagerUsecaseMock, *mocks.CrawlerDatabaseMock){
		"should return error to GetNode from pager provider": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{}
			pagerMock.On("GetNode", URI).Return(node, unexpectedErr)
			databaseMock.On("Insert", ctx, URI, depth, []string{}).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.EqualError(t, err, unexpectedErr.Error())
			assert.Empty(t, links)
		},
		"should return empty when node is nil": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			var node *html.Node
			pagerMock.On("GetNode", URI).Return(node, nil)
			databaseMock.On("Insert", ctx, URI, depth, []string{}).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return empty when not found link tag attribute": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{Type: html.ElementNode}
			pagerMock.On("GetNode", URI).Return(node, nil)
			databaseMock.On("Insert", ctx, URI, depth, []string{}).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return link fetched from provider when database returns error": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
			}
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(unexpectedErr)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			databaseMock.AssertCalled(t, "Find", ctx, URI, depth)
			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return link from database": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{internalURI}, nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			uris := []string{internalURI}
			databaseMock.AssertNotCalled(t, "Insert", ctx, URI, depth, uris)
			databaseMock.AssertCalled(t, "Find", ctx, URI, depth)
			pagerMock.AssertNotCalled(t, "GetNode", URI)
			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return link when have only one attribute": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return only one link when have two attribute but the last item has invalid key property": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}, {Key: "class", Val: "name"}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return only one link when have two attribute but the last item has invalid val link property": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}, {Key: hrefProp, Val: "index.html"}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return links when have two valid attributes": func(t *testing.T, pagerMock *mocks.PagerUsecaseMock, databaseMock *mocks.CrawlerDatabaseMock) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{Key: hrefProp, Val: internalURI},
					{Key: hrefProp, Val: randomInternalURI},
				},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", randomInternalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI, randomInternalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return links from parent and child node when first child also have next sibling": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(1)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{{Key: hrefProp, Val: randomInternalURI}},
					NextSibling: &html.Node{
						Type: html.ElementNode,
						Data: linkTag,
						Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalURI}},
					},
				},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", randomInternalURI).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", lastInternalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI, randomInternalURI, lastInternalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return links from parent and child node but will break when empty URIs": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(2)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{{Key: hrefProp, Val: randomInternalURI}},
				},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", randomInternalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI, randomInternalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return links from first and second node and need to ignore the third node to respect depth": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(2)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			firstNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
			}
			secondNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: randomInternalURI}},
			}
			thirdNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalURI}},
			}
			pagerMock.On("GetNode", URI).Return(firstNode, nil)
			pagerMock.On("GetNode", internalURI).Return(secondNode, nil)
			pagerMock.On("GetNode", randomInternalURI).Return(thirdNode, nil)
			uris := []string{internalURI, randomInternalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
		"should return links from first and second node considering when node has more than one attributes and need to respect depth": func(
			t *testing.T,
			pagerMock *mocks.PagerUsecaseMock,
			databaseMock *mocks.CrawlerDatabaseMock,
		) {
			depth := uint(2)
			databaseMock.On("Find", ctx, URI, depth).Return([]string{}, unexpectedErr)
			firstNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
			}
			secondNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: randomInternalURI}, {Key: hrefProp, Val: subInternalURI}},
			}
			thirdNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalURI}},
			}
			pagerMock.On("GetNode", URI).Return(firstNode, nil)
			pagerMock.On("GetNode", internalURI).Return(secondNode, nil)
			pagerMock.On("GetNode", randomInternalURI).Return(thirdNode, nil)
			pagerMock.On("GetNode", subInternalURI).Return(&html.Node{}, nil)
			uris := []string{internalURI, randomInternalURI, subInternalURI}
			databaseMock.On("Insert", ctx, URI, depth, uris).Return(nil)

			crawler := NewCrawlerService(pagerMock, databaseMock)
			links, err := crawler.Craw(ctx, URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, uris, links)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			pagerMock := new(mocks.PagerUsecaseMock)
			databaseMock := new(mocks.CrawlerDatabaseMock)

			run(t, pagerMock, databaseMock)
		})
	}
}

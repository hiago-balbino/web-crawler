package crawler

import (
	"errors"
	"testing"

	crawler "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/mock"
	pager "github.com/hiago-balbino/web-crawler/internal/pkg/pager/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestCrawlerPage_Craw(t *testing.T) {
	URI := "https://anyurl.com"
	internalURI := "https://internal-anyurl.com"
	randomInternalURI := "https://random-internal-anyurl.com"
	subInternalURI := "https://sub-internal-anyurl.com"
	lastInternalURI := "https://last-internal-anyurl.com"

	testCases := map[string]func(*testing.T, *pager.PagerServiceMock, *crawler.CrawlerDatabaseMock){
		"should return error to GetNode from pager provider": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			depth := 1
			node := &html.Node{}
			unknownErr := errors.New("unknown error")
			pagerMock.On("GetNode", URI).Return(node, unknownErr)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.EqualError(t, err, unknownErr.Error())
			assert.Empty(t, links)
		},
		"should return empty when node is nil": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			depth := 1
			var node *html.Node
			pagerMock.On("GetNode", URI).Return(node, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return empty when not found link tag attribute": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			depth := 1
			node := &html.Node{Type: html.ElementNode}
			pagerMock.On("GetNode", URI).Return(node, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return link when have only one attribute": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI}, links)
		},
		"should return only one link when have two attribute but the last item has invalid key property": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}, {Key: "class", Val: "name"}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI}, links)
		},
		"should return only one link when have two attribute but the last item has invalid val link property": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalURI}, {Key: hrefProp, Val: "index.html"}},
			}
			pagerMock.On("GetNode", URI).Return(node, nil)
			pagerMock.On("GetNode", internalURI).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI}, links)
		},
		"should return links when have two valid attributes": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			depth := 1
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

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI, randomInternalURI}, links)
		},
		"should return links from parent and child node when first child also have next sibling": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 1
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

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI, randomInternalURI, lastInternalURI}, links)
		},
		"should return links from parent and child node but will break when empty URIs": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 2
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

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI, randomInternalURI}, links)
		},
		"should return links from first and second node and need to ignore the third node to respect depth": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 2
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

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI, randomInternalURI}, links)
		},
		"should return links from first and second node considering when node has more than one attributes and need to respect depth": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			depth := 2
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

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(URI, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalURI, randomInternalURI, subInternalURI}, links)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			pagerMock := new(pager.PagerServiceMock)
			databaseMock := new(crawler.CrawlerDatabaseMock)

			run(t, pagerMock, databaseMock)
		})
	}
}

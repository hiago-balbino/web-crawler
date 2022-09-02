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
	testCases := map[string]func(*testing.T, *pager.PagerServiceMock, *crawler.CrawlerDatabaseMock){
		"should return error to GetNode from pager provider": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			depth := 1
			node := &html.Node{}
			unknownErr := errors.New("unknown error")
			pagerMock.On("GetNode", uri).Return(node, unknownErr)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.EqualError(t, err, unknownErr.Error())
			assert.Empty(t, links)
		},
		"should return empty when node is nil": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			depth := 1
			var node *html.Node
			pagerMock.On("GetNode", uri).Return(node, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return empty when not found link tag attribute": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			depth := 1
			node := &html.Node{Type: html.ElementNode}
			pagerMock.On("GetNode", uri).Return(node, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.Empty(t, links)
		},
		"should return link when have only one attribute": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri}, links)
		},
		"should return only one link when have two attribute but the last item has invalid key property": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}, {Key: "class", Val: "name"}},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri}, links)
		},
		"should return only one link when have two attribute but the last item has invalid val link property": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}, {Key: hrefProp, Val: "index.html"}},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri}, links)
		},
		"should return links when have two valid attributes": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{Key: hrefProp, Val: internalUri},
					{Key: hrefProp, Val: anotherInternalUri},
				},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri}, links)
		},
		"should return links from parent and child node": func(t *testing.T, pagerMock *pager.PagerServiceMock, databaseMock *crawler.CrawlerDatabaseMock) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{{Key: hrefProp, Val: anotherInternalUri}},
				},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri}, links)
		},
		"should return links from parent and child node when first child also have next sibling": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			lastInternalUri := "https://last-internal-anyurl.com"
			depth := 1
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{{Key: hrefProp, Val: anotherInternalUri}},
					NextSibling: &html.Node{
						Type: html.ElementNode,
						Data: linkTag,
						Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalUri}},
					},
				},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", lastInternalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri, lastInternalUri}, links)
		},
		"should return links from parent and child node but will break when empty URIs": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			depth := 2
			node := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{{Key: hrefProp, Val: anotherInternalUri}},
				},
			}
			pagerMock.On("GetNode", uri).Return(node, nil)
			pagerMock.On("GetNode", internalUri).Return(&html.Node{}, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri}, links)
		},
		"should return links from first and second node and need to ignore the third node to respect depth": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {
			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			lastInternalUri := "https://last-internal-anyurl.com"
			depth := 2
			firstNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
			}
			secondNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: anotherInternalUri}},
			}
			thirdNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalUri}},
			}
			pagerMock.On("GetNode", uri).Return(firstNode, nil)
			pagerMock.On("GetNode", internalUri).Return(secondNode, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(thirdNode, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri}, links)
		},
		"should return links from first and second node considering when node has more than one attributes and need to respect depth": func(
			t *testing.T,
			pagerMock *pager.PagerServiceMock,
			databaseMock *crawler.CrawlerDatabaseMock,
		) {

			uri := "https://anyurl.com"
			internalUri := "https://internal-anyurl.com"
			anotherInternalUri := "https://another-internal-anyurl.com"
			subInternalUri := "https://sub-internal-anyurl.com"
			lastInternalUri := "https://last-internal-anyurl.com"
			depth := 2
			firstNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: internalUri}},
			}
			secondNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: anotherInternalUri}, {Key: hrefProp, Val: subInternalUri}},
			}
			thirdNode := &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{Key: hrefProp, Val: lastInternalUri}},
			}
			pagerMock.On("GetNode", uri).Return(firstNode, nil)
			pagerMock.On("GetNode", internalUri).Return(secondNode, nil)
			pagerMock.On("GetNode", anotherInternalUri).Return(thirdNode, nil)
			pagerMock.On("GetNode", subInternalUri).Return(&html.Node{}, nil)

			crawler := NewCrawlerPage(pagerMock, databaseMock)
			links, err := crawler.Craw(uri, depth)

			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{internalUri, anotherInternalUri, subInternalUri}, links)
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

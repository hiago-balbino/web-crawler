package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestCrawlerPage_Craw(t *testing.T) {

}

func TestExtractAddresses(t *testing.T) {
	testCases := []struct {
		name          string
		node          *html.Node
		expectedLinks []string
	}{
		{
			name: "should return link when have only one attribute",
			node: &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{{
					Key: hrefProp,
					Val: "https://google.com",
				}},
			},
			expectedLinks: []string{"https://google.com"},
		},
		{
			name: "should return link when have more than one attribute",
			node: &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{
						Key: hrefProp,
						Val: "https://google.com",
					},
					{
						Key: hrefProp,
						Val: "https://twitter.com",
					},
				},
			},
			expectedLinks: []string{"https://google.com", "https://twitter.com"},
		},
		{
			name: "should return only one link when have two attribute but the last is invalid",
			node: &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{
						Key: hrefProp,
						Val: "https://google.com",
					},
					{
						Key: "class",
						Val: "name",
					},
				},
			},
			expectedLinks: []string{"https://google.com"},
		},
		{
			name: "should return links from parent and child node",
			node: &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{
						Key: hrefProp,
						Val: "https://google.com",
					},
				},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{
						{
							Key: hrefProp,
							Val: "https://twitter.com",
						},
					},
				},
			},
			expectedLinks: []string{"https://google.com", "https://twitter.com"},
		},
		{
			name: "should return links from parent and child node when first child also have next sibling",
			node: &html.Node{
				Type: html.ElementNode,
				Data: linkTag,
				Attr: []html.Attribute{
					{
						Key: hrefProp,
						Val: "https://google.com",
					},
				},
				FirstChild: &html.Node{
					Type: html.ElementNode,
					Data: linkTag,
					Attr: []html.Attribute{
						{
							Key: hrefProp,
							Val: "https://twitter.com",
						},
					},
					NextSibling: &html.Node{
						Type: html.ElementNode,
						Data: linkTag,
						Attr: []html.Attribute{
							{
								Key: hrefProp,
								Val: "https://twitch.com",
							},
						},
					},
				},
			},
			expectedLinks: []string{"https://google.com", "https://twitter.com", "https://twitch.com"},
		},
	}

	for _, test := range testCases {
		links := extractAddresses(nil, test.node)

		assert.ElementsMatch(t, test.expectedLinks, links)
	}
}

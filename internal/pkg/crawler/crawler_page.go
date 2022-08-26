package crawler

import (
	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/internal/core/pager"
	"golang.org/x/net/html"
)

const (
	linkTag  = "a"
	hrefProp = "href"
)

// CrawlerPage is a implementation to handle with web crawler
type CrawlerPage struct {
	provider pager.Pager
}

// NewCrawlerPage is a constructor to create a new instance of CrawlerPage
func NewCrawlerPage(pager pager.Pager) crawler.Crawler {
	return CrawlerPage{provider: pager}
}

// Craw execute the call to craw pages
func (p CrawlerPage) Craw(uri string) {
	panic("implement me")
}

// extractAddresses recursively extracts the addresses of the HTML node
func extractAddresses(links []string, node *html.Node) []string {
	if node.Type == html.ElementNode && node.Data == linkTag {
		for _, attr := range node.Attr {
			if attr.Key == hrefProp {
				links = append(links, attr.Val)
			}
		}
	}

	for next := node.FirstChild; next != nil; next = next.NextSibling {
		links = extractAddresses(links, next)
	}

	return links
}

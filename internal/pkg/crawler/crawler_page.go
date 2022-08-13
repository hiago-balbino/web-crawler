package crawler

import (
	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/internal/core/pager"
)

// CrawlerPage is a implementation to handle with web crawler
type CrawlerPage struct {
	pager.Pager
}

// NewCrawlerPage is a constructor to create a new instance of CrawlerPage
func NewCrawlerPage(pager pager.Pager) crawler.Crawler {
	return CrawlerPage{pager}
}

// Craw execute the call to craw pages
func (p CrawlerPage) Craw(uri string) {
	panic("implement me")
}

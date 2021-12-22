package crawler

import "github.com/hiago-balbino/web-crawler/src/core/page"

// PageCrawler is a implementation to handle with web crawler
type PageCrawler struct {
	page.Pager
}

// NewPageCrawler is a constructor to create a new instance of PageCrawler
func NewPageCrawler(pager page.Pager) Crawler {
	return PageCrawler{pager}
}

// Run execute the call to craw pages
func (p PageCrawler) Run(uri string) {
	panic("implement me")
}

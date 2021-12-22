package crawler

import "github.com/hiago-balbino/web-crawler/src/core/pager"

// PageCrawler is a implementation to handle with web crawler
type PageCrawler struct {
	pager.Pager
}

// NewPageCrawler is a constructor to create a new instance of PageCrawler
func NewPageCrawler(pager pager.Pager) Crawler {
	return PageCrawler{pager}
}

// Run execute the call to craw pages
func (p PageCrawler) Run(uri string) {
	panic("implement me")
}

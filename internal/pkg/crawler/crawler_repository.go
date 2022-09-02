package crawler

import (
	"context"

	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
)

// CrawlerRepository is a repository to handle with persistence layer
type CrawlerRepository struct {
	crawler.CrawlerDatabase
}

// NewCrawlerRepository is a constructor to create a new instance of CrawlerRepository
func NewCrawlerRepository(repository crawler.CrawlerDatabase) CrawlerRepository {
	return CrawlerRepository{repository}
}

// Find is a method to fetch links crawled from database
func (c CrawlerRepository) Find(ctx context.Context, uri string) ([]string, error) {
	panic("not implemented")
}

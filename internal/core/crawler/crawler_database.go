package crawler

import "context"

// CrawlerDatabase is a repository to handle with persistence layer.
type CrawlerDatabase interface {
	// Find is a method to fetch links crawled from database.
	Find(ctx context.Context, uri string) ([]string, error)
}

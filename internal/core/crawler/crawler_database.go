package crawler

import "context"

// CrawlerDatabase is a repository to handle with persistence layer.
type CrawlerDatabase interface {
	// Insert is a method to insert new page crawled on database.
	Insert(ctx context.Context, uri string, depth int, uris []string) error

	// Find is a method to fetch links crawled from database.
	Find(ctx context.Context, uri string, depth int) ([]string, error)
}

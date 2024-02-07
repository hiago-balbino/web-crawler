package crawler

import "context"

type CrawlerDatabase interface {
	Insert(ctx context.Context, uri string, depth uint, uris []string) error
	Find(ctx context.Context, uri string, depth uint) ([]string, error)
}

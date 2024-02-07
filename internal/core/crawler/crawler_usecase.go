package crawler

import "context"

type CrawlerUsecase interface {
	Craw(ctx context.Context, uri string, depth uint) ([]string, error)
}

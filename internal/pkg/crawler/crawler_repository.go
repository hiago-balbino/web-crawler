package crawler

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CrawlerRepository is a repository to handle with persistence layer
type CrawlerRepository struct {
	client *mongo.Client
}

// NewCrawlerRepository is a constructor to create a new instance of CrawlerRepository
func NewCrawlerRepository(ctx context.Context) CrawlerRepository {
	opts := options.Client().ApplyURI("")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	return CrawlerRepository{client}
}

// Find is a method to fetch links crawled from database
func (c CrawlerRepository) Find(ctx context.Context, uri string) ([]string, error) {
	panic("not implemented")
}

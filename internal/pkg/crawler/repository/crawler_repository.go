package crawler

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CrawlerRepository is a repository to handle with persistence layer.
type CrawlerRepository struct {
	client *mongo.Client
}

// NewCrawlerRepository is a constructor to create a new instance of CrawlerRepository.
func NewCrawlerRepository(ctx context.Context) CrawlerRepository {
	username := viper.GetString("MONGODB_USERNAME")
	password := viper.GetString("MONGODB_PASSWORD")
	host := viper.GetString("MONGODB_HOST")
	port := viper.GetString("MONGODB_PORT")
	endpoint := fmt.Sprintf("mongodb://%s:%s@%s", username, password, net.JoinHostPort(host, port))

	opts := options.Client().ApplyURI(endpoint)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err.Error())
	}

	return CrawlerRepository{client}
}

// Find is a method to fetch links crawled from database.
func (c CrawlerRepository) Find(ctx context.Context, uri string, depth int) ([]string, error) {
	databaseName := viper.GetString("MONGODB_DATABASE")
	collectionName := viper.GetString("MONGODB_COLLECTION")
	collection := c.client.Database(databaseName).Collection(collectionName)
	filter := bson.D{{Key: "uri", Value: uri}, {Key: "depth", Value: depth}}

	dataPage := dataPage{}
	err := collection.FindOne(ctx, filter).Decode(&dataPage)
	if err != nil {
		return nil, err
	}

	return dataPage.URIs, nil
}

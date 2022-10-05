package crawler

import (
	"context"
	"fmt"
	"net"

	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.GetLogger()

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
	if noUserInformation(username, password) {
		endpoint = fmt.Sprintf("mongodb://%s", net.JoinHostPort(host, port))
	}

	opts := options.Client().ApplyURI(endpoint)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Error("error connecting to mongodb")
	}

	return CrawlerRepository{client}
}

func noUserInformation(username, password string) bool {
	return username == "" && password == ""
}

// Insert is a method to insert new page crawled on database.
func (c CrawlerRepository) Insert(ctx context.Context, uri string, depth int, uris []string) error {
	dataPage := dataPage{
		URI:   uri,
		Depth: depth,
		URIs:  uris,
	}
	_, err := c.getCollection().InsertOne(ctx, dataPage)
	if err != nil {
		log.Error("error while inserting new data into collection")

		return err
	}

	return nil
}

// Find is a method to fetch links crawled from database.
func (c CrawlerRepository) Find(ctx context.Context, uri string, depth int) ([]string, error) {
	filter := bson.D{{Key: "uri", Value: uri}, {Key: "depth", Value: depth}}
	dataPage := dataPage{}
	err := c.getCollection().FindOne(ctx, filter).Decode(&dataPage)
	if err != nil {
		log.Error("error while fetching data from collection")

		return nil, err
	}

	return dataPage.URIs, nil
}

func (c CrawlerRepository) getCollection() *mongo.Collection {
	databaseName := viper.GetString("MONGODB_DATABASE")
	collectionName := viper.GetString("MONGODB_COLLECTION")

	return c.client.Database(databaseName).Collection(collectionName)
}
package storage

import (
	"context"
	"fmt"
	"net"

	"github.com/hiago-balbino/web-crawler/v2/internal/pkg/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.GetLogger()

type CrawlerMongodbRepository struct {
	client *mongo.Client
}

func NewCrawlerMongodbRepository(ctx context.Context) CrawlerMongodbRepository {
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
		log.Error("error connecting to mongodb", logger.FieldError(err))
	}

	return CrawlerMongodbRepository{client}
}

func noUserInformation(username, password string) bool {
	return username == "" && password == ""
}

func (c CrawlerMongodbRepository) Insert(ctx context.Context, uri string, depth uint, uris []string) error {
	pageDataInfo := pageDataInfo{
		URI:   uri,
		Depth: depth,
		URIs:  uris,
	}
	_, err := c.getCollection().InsertOne(ctx, pageDataInfo)
	if err != nil {
		log.Error("error while inserting new data into collection", logger.FieldError(err))

		return err
	}

	return nil
}

func (c CrawlerMongodbRepository) Find(ctx context.Context, uri string, depth uint) ([]string, error) {
	filter := bson.D{{Key: "uri", Value: uri}, {Key: "depth", Value: depth}}
	pageDataInfo := pageDataInfo{}
	err := c.getCollection().FindOne(ctx, filter).Decode(&pageDataInfo)
	if err != nil {
		log.Error("error while fetching data from collection", logger.FieldError(err))

		return nil, err
	}

	return pageDataInfo.URIs, nil
}

func (c CrawlerMongodbRepository) getCollection() *mongo.Collection {
	databaseName := viper.GetString("MONGODB_DATABASE")
	collectionName := viper.GetString("MONGODB_COLLECTION")

	return c.client.Database(databaseName).Collection(collectionName)
}

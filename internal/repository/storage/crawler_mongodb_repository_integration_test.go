package storage

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongodbRepositoryIntegrationTestSuite struct {
	suite.Suite
	testcontainers.Container
	repository CrawlerMongodbRepository
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(MongodbRepositoryIntegrationTestSuite))
}

func (suite *MongodbRepositoryIntegrationTestSuite) SetupSuite() {
	suite.defaultDBEnviroments()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	genericRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}
	container, err := testcontainers.GenericContainer(ctx, genericRequest)
	assert.NoError(suite.T(), err)

	host, err := container.Host(ctx)
	assert.NoError(suite.T(), err)
	viper.Set("MONGODB_HOST", host)

	port, err := container.MappedPort(ctx, "27017")
	assert.NoError(suite.T(), err)
	viper.Set("MONGODB_PORT", string(port))

	suite.Container = container
	suite.repository = NewCrawlerMongodbRepository(ctx)
}

func (suite *MongodbRepositoryIntegrationTestSuite) TearDownSuite() {
	err := suite.Terminate(context.Background())
	assert.NoError(suite.T(), err)

	viper.Reset()
}

func (suite *MongodbRepositoryIntegrationTestSuite) TestInsert() {
	ctx := context.Background()
	uri := "http://crawler.com"
	depth := uint(1)
	uris := []string{"http://subcrawler.com"}

	suite.Suite.T().Run("should return error to insert when invalid database name", func(t *testing.T) {
		defer func() {
			suite.defaultDBEnviroments()
		}()
		viper.Set("MONGODB_DATABASE", "")

		repository := NewCrawlerMongodbRepository(ctx)
		err := repository.Insert(ctx, uri, depth, uris)

		assert.NotNil(suite.T(), err)
	})

	suite.Suite.T().Run("should insert data page with success", func(t *testing.T) {
		err := suite.repository.Insert(ctx, uri, depth, uris)

		assert.NoError(suite.T(), err)
	})
}

func (suite *MongodbRepositoryIntegrationTestSuite) TestFind() {
	ctx := context.Background()
	uri := "http://crawler.com"
	depth := uint(1)
	uris := []string{"http://subcrawler.com"}

	suite.Suite.T().Run("should return error to find URIs stored", func(t *testing.T) {
		defer func() {
			suite.defaultDBEnviroments()
		}()
		viper.Set("MONGODB_DATABASE", "")

		repository := NewCrawlerMongodbRepository(ctx)
		storedURIs, err := repository.Find(ctx, uri, depth)

		assert.NotNil(suite.T(), err)
		assert.Empty(suite.T(), storedURIs)
	})

	suite.Suite.T().Run("should return empty slice when try to find URIs stored", func(t *testing.T) {
		storedURIs, err := suite.repository.Find(ctx, uri, depth)

		assert.EqualError(suite.T(), err, mongo.ErrNoDocuments.Error())
		assert.Empty(suite.T(), storedURIs)
	})

	suite.Suite.T().Run("should return stored URIs with success", func(t *testing.T) {
		dataPage := pageDataInfo{URI: uri, Depth: depth, URIs: uris}
		_, err := suite.repository.getCollection().InsertOne(ctx, dataPage)
		assert.NoError(suite.T(), err)

		storedURIs, err := suite.repository.Find(ctx, uri, depth)

		assert.NoError(suite.T(), err)
		assert.ElementsMatch(suite.T(), uris, storedURIs)
	})
}

func (suite *MongodbRepositoryIntegrationTestSuite) defaultDBEnviroments() {
	viper.Set("MONGODB_DATABASE", "database_test")
	viper.Set("MONGODB_COLLECTION", "collection_test")
}

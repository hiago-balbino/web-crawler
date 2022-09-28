package crawler

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	testcontainers.Container
	repository CrawlerRepository
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(RepositoryIntegrationTestSuite))
}

func (suite *RepositoryIntegrationTestSuite) SetupSuite() {
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
	suite.repository = NewCrawlerRepository(ctx)
}

func (suite *RepositoryIntegrationTestSuite) TearDownSuite() {
	err := suite.Terminate(context.Background())
	assert.NoError(suite.T(), err)

	viper.Reset()
}

func (suite *RepositoryIntegrationTestSuite) TestInsert() {
	err := suite.repository.client.Ping(context.Background(), nil)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestFind() {
	err := suite.repository.client.Ping(context.Background(), nil)
	assert.NoError(suite.T(), err)
}

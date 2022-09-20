package crawler

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(RepositoryIntegrationTestSuite))
}

func (suite *RepositoryIntegrationTestSuite) SetupTest() {}

func (suite *RepositoryIntegrationTestSuite) TearDownTest() {}

func (suite *RepositoryIntegrationTestSuite) TestInsert() {}

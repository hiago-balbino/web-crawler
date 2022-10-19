package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiago-balbino/web-crawler/configuration"
	core "github.com/hiago-balbino/web-crawler/internal/core/crawler"
	repo "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/repository"
	crawler "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/service"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	pager "github.com/hiago-balbino/web-crawler/internal/pkg/pager/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = logger.GetLogger()

// Server is an structure to support the API.
type Server struct {
	crawlerService core.CrawlerService
}

// NewServer create a new instante of Server structure.
func NewServer() Server {
	configuration.InitConfigurations()

	ctx := context.Background()
	httpClient := new(http.Client)
	provider := pager.NewPagerProvider(httpClient)
	database := repo.NewCrawlerRepository(ctx)
	service := crawler.NewCrawlerPage(provider, database)

	return Server{crawlerService: service}
}

// Start initialize the API.
func (s Server) Start() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WORKING")
	})

	err := router.Run()
	if err != nil {
		log.Fatal("error while server starting", zap.Field{Type: zapcore.StringType, String: err.Error()})
	}
}

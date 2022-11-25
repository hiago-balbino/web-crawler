package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiago-balbino/web-crawler/configuration"
	repo "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/repository"
	crawler "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/service"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	pager "github.com/hiago-balbino/web-crawler/internal/pkg/pager/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = logger.GetLogger()

// Server is an structure to support the API.
type Server struct {
	handler Handler
}

// NewServer create a new instante of Server structure.
func NewServer() Server {
	configuration.InitConfigurations()
	provider := pager.NewPagerProvider(new(http.Client))
	database := repo.NewCrawlerRepository(context.Background())
	service := crawler.NewCrawlerPage(provider, database)
	handler := NewHandler(service)

	return Server{handler: handler}
}

// Start initialize the API.
func (s Server) Start() {
	router := s.setupRoutes("templates/*")

	if err := router.Run(fmt.Sprintf(":%s", viper.GetString("PORT"))); err != nil {
		log.Fatal("error while server starting", zap.Field{Type: zapcore.StringType, String: err.Error()})
	}
}

func (s Server) setupRoutes(templatePath string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob(templatePath)

	router.GET("/index", s.handler.index)
	router.GET("/crawler", s.handler.getCrawledPage)

	return router
}

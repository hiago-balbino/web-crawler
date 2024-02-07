package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiago-balbino/web-crawler/config"
	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/internal/core/pager"
	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	"github.com/hiago-balbino/web-crawler/internal/repository/storage"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = logger.GetLogger()

type Server struct {
	handler Handler
}

func NewServer() Server {
	config.InitConfigurations()
	pagerService := pager.NewPagerService(new(http.Client))
	crawlerDatabase := storage.NewCrawlerMongodbRepository(context.Background())
	crawlerService := crawler.NewCrawlerService(pagerService, crawlerDatabase)
	handler := NewHandler(crawlerService)

	return Server{handler: handler}
}

func (s Server) Start() {
	router := s.setupRoutes("web/templates/*")

	monitor := ginmetrics.GetMonitor()
	monitor.SetMetricPath("/metrics")
	monitor.Use(router)

	if err := router.Run(fmt.Sprintf(":%s", viper.GetString("API_PORT"))); err != nil {
		log.Fatal("error while server starting", zap.Field{Type: zapcore.StringType, String: err.Error()})
	}
}

func (s Server) setupRoutes(templatePath string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob(templatePath)

	router.GET("/index", s.handler.index)
	router.GET("/crawler", s.handler.getPageCrawled)

	return router
}

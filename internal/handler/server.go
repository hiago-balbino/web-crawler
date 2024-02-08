package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiago-balbino/web-crawler/v2/config"
	"github.com/hiago-balbino/web-crawler/v2/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/v2/internal/core/pager"
	"github.com/hiago-balbino/web-crawler/v2/internal/pkg/logger"
	"github.com/hiago-balbino/web-crawler/v2/internal/repository/storage"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/spf13/viper"
)

var log = logger.GetLogger()

type Server struct {
	handler Handler
}

func NewServer() Server {
	config.InitConfigurations()
	pagerService := pager.NewPagerService(&http.Client{Timeout: viper.GetDuration("API_REQUEST_TIMEOUT")})
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
		log.Fatal("error while server starting", logger.FieldError(err))
	}
}

func (s Server) setupRoutes(templatePath string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob(templatePath)

	router.GET("/index", s.handler.index)
	router.GET("/crawler", s.handler.getPageCrawled)

	return router
}

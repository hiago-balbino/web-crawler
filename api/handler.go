package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	core "github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Handler is an struct that contains service and functions to handle in API.
type Handler struct {
	service core.CrawlerService
}

// NewHandler is a constructor to create a new instance of Handler.
func NewHandler(service core.CrawlerService) Handler {
	return Handler{service: service}
}

// getCrawledPage is a function to handle with crawled page.
func (h Handler) getCrawledPage(c *gin.Context) {
	var schema requestSchema
	if err := c.BindQuery(&schema); err != nil {
		log.Error("error binding query params", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.JSON(http.StatusBadRequest, errMessage{err.Error()})

		return
	}

	if err := schema.validate(); err != nil {
		log.Error("error validating parameters", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.JSON(http.StatusBadRequest, errMessage{err.Error()})
	}
}

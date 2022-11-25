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
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})

		return
	}

	if err := schema.validate(); err != nil {
		log.Error("error validating parameters", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})

		return
	}

	links, err := h.service.Craw(c.Request.Context(), schema.URI, schema.Depth)
	if err != nil {
		log.Error("error crawling page", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})

		return
	}

	if len(links) == 0 {
		c.HTML(http.StatusOK, "empty_result.html", gin.H{"message": "The process did not return any valid results"})

		return
	}

	c.HTML(http.StatusOK, "links.html", gin.H{"links": links})
}

// index is a function to return index HTML page.
func (h Handler) index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

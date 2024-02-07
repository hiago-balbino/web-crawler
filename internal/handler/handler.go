package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	core "github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Handler struct {
	service core.CrawlerUsecase
}

func NewHandler(service core.CrawlerUsecase) Handler {
	return Handler{service: service}
}

func (h Handler) getPageCrawled(c *gin.Context) {
	var crawPageInfo crawPageInfo
	if err := c.BindQuery(&crawPageInfo); err != nil {
		log.Error("error binding query params", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})

		return
	}

	if err := crawPageInfo.validate(); err != nil {
		log.Error("error validating parameters", zap.Field{Type: zapcore.StringType, String: err.Error()})
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})

		return
	}

	links, err := h.service.Craw(c.Request.Context(), crawPageInfo.URI, crawPageInfo.Depth)
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

func (h Handler) index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

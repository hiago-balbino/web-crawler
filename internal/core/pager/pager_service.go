package pager

import (
	"net/http"

	"github.com/hiago-balbino/web-crawler/internal/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/html"
)

var log = logger.GetLogger()

type PagerService struct {
	httpClient *http.Client
}

func NewPagerService(httpClient *http.Client) PagerService {
	return PagerService{httpClient: httpClient}
}

func (c PagerService) GetNode(uri string) (*html.Node, error) {
	response, err := c.httpClient.Get(uri)
	defer func() {
		if err == nil {
			_ = response.Body.Close()
		}
	}()
	if err != nil {
		log.Error("error to perform get request in provider", zap.Field{Type: zapcore.StringType, String: err.Error()})

		return nil, err
	}

	node, err := html.Parse(response.Body)
	if err != nil {
		log.Error("error to parse response body to html", zap.Field{Type: zapcore.StringType, String: err.Error()})

		return nil, err
	}

	return node, nil
}

package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/hiago-balbino/web-crawler/test/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGetPageCrawled(t *testing.T) {
	givenDepth := uint(2)
	givenURI := "https://anyuritest.com"

	t.Run("should return 4xx error", func(t *testing.T) {
		handler := setupHandler(nil)
		server := httptest.NewServer(handler)
		defer server.Close()

		e := httpexpect.Default(t, server.URL)

		t.Run("when empty URI query param", func(t *testing.T) {
			e.GET("/crawler").
				WithQuery("depth", givenDepth).
				Expect().
				Status(http.StatusBadRequest).
				Body().Contains(errEmptyURI.Error())
		})
		t.Run("when empty depth query param", func(t *testing.T) {
			e.GET("/crawler").
				WithQuery("uri", givenURI).
				Expect().
				Status(http.StatusBadRequest).
				Body().Contains(errEmptyDepth.Error())
		})
		t.Run("when negative depth query param", func(t *testing.T) {
			e.GET("/crawler").
				WithQuery("uri", givenURI).
				WithQuery("depth", -1).
				Expect().
				Status(http.StatusBadRequest).
				Body().Contains(`strconv.ParseUint: parsing &#34;-1&#34;: invalid syntax`)
		})
	})

	t.Run("should return 5xx error when fail to perform HTTP request to fetch page", func(t *testing.T) {
		unexpectedErr := errors.New("unexpected error")
		crawlerService := new(mocks.CrawlerUsecaseMock)
		crawlerService.On("Craw", mock.Anything, givenURI, givenDepth).Return(make([]string, 0), unexpectedErr)

		handler := setupHandler(crawlerService)
		server := httptest.NewServer(handler)
		defer server.Close()

		e := httpexpect.Default(t, server.URL)

		e.GET("/crawler").
			WithQuery("uri", givenURI).
			WithQuery("depth", givenDepth).
			Expect().
			Status(http.StatusInternalServerError).
			Body().Contains(unexpectedErr.Error())
	})

	t.Run("should return 2xx", func(t *testing.T) {
		t.Run("when process did not return any results", func(t *testing.T) {
			links := make([]string, 0)
			crawlerService := new(mocks.CrawlerUsecaseMock)
			crawlerService.On("Craw", mock.Anything, givenURI, givenDepth).Return(links, nil)

			handler := setupHandler(crawlerService)
			server := httptest.NewServer(handler)
			defer server.Close()

			e := httpexpect.Default(t, server.URL)

			e.GET("/crawler").
				WithQuery("uri", givenURI).
				WithQuery("depth", givenDepth).
				Expect().
				Status(http.StatusOK).
				Body().Contains("The process did not return any valid results")
		})
		t.Run("when page is successfully crawled", func(t *testing.T) {
			links := []string{"https://firstlink.com", "https://secondlink.com", "https://thirdlink.com"}
			crawlerService := new(mocks.CrawlerUsecaseMock)
			crawlerService.On("Craw", mock.Anything, givenURI, givenDepth).Return(links, nil)

			handler := setupHandler(crawlerService)
			server := httptest.NewServer(handler)
			defer server.Close()

			e := httpexpect.Default(t, server.URL)

			e.GET("/crawler").
				WithQuery("uri", givenURI).
				WithQuery("depth", givenDepth).
				Expect().
				Status(http.StatusOK).
				Body().
				Contains(links[0]).
				Contains(links[1]).
				Contains(links[2])
		})
	})
}

func TestIndex(t *testing.T) {
	t.Run("should return 2xx when load index page", func(t *testing.T) {
		handler := setupHandler(nil)
		server := httptest.NewServer(handler)
		defer server.Close()

		e := httpexpect.Default(t, server.URL)

		e.GET("/index").
			Expect().
			Status(http.StatusOK)
	})
}

func setupHandler(service crawler.CrawlerUsecase) *gin.Engine {
	handler := NewHandler(service)
	server := Server{handler: handler}
	router := server.setupRoutes("../../web/templates/*")

	return router
}

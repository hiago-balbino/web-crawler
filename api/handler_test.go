package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	mock "github.com/hiago-balbino/web-crawler/internal/pkg/crawler/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetCrawledPage(t *testing.T) {
	givenDepth := uint(2)
	givenURI := "https://anyuritest.com"

	t.Run("should return 4xx error", func(t *testing.T) {
		t.Run("when empty URI query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?depth=%d", givenDepth)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Contains(t, response.Body.String(), errEmptyURI.Error())
		})
		t.Run("when empty depth query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?uri=%s", givenURI)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Contains(t, response.Body.String(), errEmptyDepth.Error())
		})
		t.Run("when negative depth query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?uri=%s&depth=%d", givenURI, -1)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Contains(t, response.Body.String(), `strconv.ParseUint: parsing &#34;-1&#34;: invalid syntax`)
		})
	})

	t.Run("should return 5xx error when fail to perform HTTP request to fetch page", func(t *testing.T) {
		path := fmt.Sprintf("/crawler?uri=%s&depth=%d", givenURI, givenDepth)
		request, err := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		unexpectedErr := errors.New("unexpected error")
		crawlerService := new(mock.CrawlerServiceMock)
		crawlerService.On("Craw", context.Background(), givenURI, givenDepth).Return(make([]string, 0), unexpectedErr)
		runServerHTTPTest(request, response, crawlerService)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Contains(t, response.Body.String(), unexpectedErr.Error())
	})

	t.Run("should return 2xx", func(t *testing.T) {
		path := fmt.Sprintf("/crawler?uri=%s&depth=%d", givenURI, givenDepth)
		request, err := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		t.Run("when process did not return any results", func(t *testing.T) {
			links := make([]string, 0)
			crawlerService := new(mock.CrawlerServiceMock)
			crawlerService.On("Craw", context.Background(), givenURI, givenDepth).Return(links, nil)
			runServerHTTPTest(request, response, crawlerService)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Contains(t, response.Body.String(), "The process did not return any valid results")
		})
		t.Run("when page is successfully crawled", func(t *testing.T) {
			links := []string{"https://firstlink.com", "https://secondlink.com", "https://thirdlink.com"}
			crawlerService := new(mock.CrawlerServiceMock)
			crawlerService.On("Craw", context.Background(), givenURI, givenDepth).Return(links, nil)
			runServerHTTPTest(request, response, crawlerService)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Contains(t, response.Body.String(), links[0])
			assert.Contains(t, response.Body.String(), links[1])
			assert.Contains(t, response.Body.String(), links[2])
		})
	})
}

func TestIndex(t *testing.T) {
	t.Run("should return 2xx when load index page", func(t *testing.T) {
		path := "/index"
		request, err := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		runServerHTTPTest(request, response, nil)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func runServerHTTPTest(request *http.Request, response *httptest.ResponseRecorder, service crawler.CrawlerService) {
	handler := NewHandler(service)
	server := Server{handler: handler}
	router := server.setupRoutes("../templates/*")
	router.ServeHTTP(response, request)
}

package api

import (
	"context"
	"encoding/json"
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

			expectedMessage, errMarshal := json.Marshal(errMessage{errEmptyURI.Error()})
			assert.NoError(t, errMarshal)
			assert.Equal(t, string(expectedMessage), response.Body.String())
			assert.Equal(t, http.StatusBadRequest, response.Code)
		})
		t.Run("when empty depth query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?uri=%s", givenURI)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			expectedMessage, errMarshal := json.Marshal(errMessage{errEmptyDepth.Error()})
			assert.NoError(t, errMarshal)
			assert.Equal(t, string(expectedMessage), response.Body.String())
			assert.Equal(t, http.StatusBadRequest, response.Code)
		})
		t.Run("when negative depth query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?uri=%s&depth=%d", givenURI, -1)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			assert.Contains(t, response.Body.String(), `{"message":"strconv.ParseUint: parsing \"-1\": invalid syntax"}`)
			assert.Equal(t, http.StatusBadRequest, response.Code)
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

		expectedMessage, errMarshal := json.Marshal(errMessage{unexpectedErr.Error()})
		assert.NoError(t, errMarshal)
		assert.Equal(t, string(expectedMessage), response.Body.String())
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("should return 2xx and links when page is successfully crawled", func(t *testing.T) {
		path := fmt.Sprintf("/crawler?uri=%s&depth=%d", givenURI, givenDepth)
		request, err := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		assert.NoError(t, err)

		links := []string{"https://firstlink.com", "https://secondlink.com", "https://thirdlink.com"}
		crawlerService := new(mock.CrawlerServiceMock)
		crawlerService.On("Craw", context.Background(), givenURI, givenDepth).Return(links, nil)
		runServerHTTPTest(request, response, crawlerService)

		expectedLinks, errMarshal := json.Marshal(responseSchema{links})
		assert.NoError(t, errMarshal)
		assert.Equal(t, string(expectedLinks), response.Body.String())
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func runServerHTTPTest(request *http.Request, response *httptest.ResponseRecorder, service crawler.CrawlerService) {
	handler := NewHandler(service)
	server := Server{handler: handler}
	router := server.setupRoutes()
	router.ServeHTTP(response, request)
}

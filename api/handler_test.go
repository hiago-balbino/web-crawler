package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiago-balbino/web-crawler/internal/core/crawler"
	"github.com/stretchr/testify/assert"
)

func TestGetCrawledPage(t *testing.T) {
	givenDepth := 2
	givenURI := "https://anyuritest.com"

	t.Run("should return 4xx error", func(t *testing.T) {
		t.Run("when empty URI query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?depth=%d", givenDepth)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			expectedMessage, err := json.Marshal(errMessage{errEmptyURI.Error()})
			assert.NoError(t, err)
			assert.Equal(t, string(expectedMessage), response.Body.String())
			assert.Equal(t, http.StatusBadRequest, response.Code)
		})
		t.Run("when empty depth query param", func(t *testing.T) {
			path := fmt.Sprintf("/crawler?uri=%s", givenURI)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			assert.NoError(t, err)

			runServerHTTPTest(request, response, nil)

			expectedMessage, err := json.Marshal(errMessage{errEmptyDepth.Error()})
			assert.NoError(t, err)
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
}

func runServerHTTPTest(request *http.Request, response *httptest.ResponseRecorder, service crawler.CrawlerService) {
	handler := NewHandler(service)
	server := Server{handler: handler}
	router := server.setupRoutes()
	router.ServeHTTP(response, request)
}

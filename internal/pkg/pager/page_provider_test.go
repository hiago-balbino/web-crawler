package pager

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestPageProvider_Get(t *testing.T) {
	testCases := []struct {
		name          string
		uri           string
		statusCode    int
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:          "should return error when http connection fails",
			uri:           "http://invalid.com",
			statusCode:    http.StatusInternalServerError,
			isExpectedErr: true,
			expectedErr:   errors.New("invalid host"),
		},
		{
			name:          "should return 5xx status code when http connection has valid host but fails",
			uri:           "http://google.com",
			statusCode:    http.StatusInternalServerError,
			isExpectedErr: true,
			expectedErr:   errors.New("fail to process"),
		},
		{
			name:          "should return 2xx status code when http connection has success",
			uri:           "http://google.com",
			statusCode:    http.StatusOK,
			isExpectedErr: false,
			expectedErr:   nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			defer gock.Off()

			requestIntercept := gock.New(test.uri)
			requestIntercept.Reply(test.statusCode)
			if test.isExpectedErr {
				requestIntercept.ReplyError(test.expectedErr)
			}
			httpClient := &http.Client{Transport: &http.Transport{}}
			gock.InterceptClient(httpClient)

			response, err := NewPageProvider(httpClient).Get(test.uri)

			if test.isExpectedErr {
				assert.ErrorIs(t, err, test.expectedErr)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, test.statusCode, response.StatusCode)
			}
		})
	}
}

package pager

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestPagerProvider_GetNode(t *testing.T) {
	testCases := []struct {
		name          string
		uri           string
		statusCode    int
		body          string
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:          "should return error when http connection fails",
			uri:           "http://invalid.com",
			statusCode:    http.StatusInternalServerError,
			body:          "",
			isExpectedErr: true,
			expectedErr:   errors.New("invalid host"),
		},
		{
			name:          "should return error when http connection has valid host but fails",
			uri:           "http://google.com",
			statusCode:    http.StatusInternalServerError,
			body:          "",
			isExpectedErr: true,
			expectedErr:   errors.New("fail to process"),
		},
		{
			name:          "should not return error when http connection and html parse execute successfully",
			uri:           "http://google.com",
			statusCode:    http.StatusOK,
			body:          `<a href="http://google.com">link</a>`,
			isExpectedErr: false,
			expectedErr:   nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			defer gock.Off()
			httpClient := httpClientMock(test)

			node, err := NewPagerProvider(httpClient).GetNode(test.uri)

			if test.isExpectedErr {
				assert.ErrorIs(t, err, test.expectedErr)
				assert.Nil(t, node)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, node)
			}
		})
	}
}

func httpClientMock(
	test struct {
		name          string
		uri           string
		statusCode    int
		body          string
		isExpectedErr bool
		expectedErr   error
	},
) *http.Client {

	requestIntercept := gock.New(test.uri)
	requestIntercept.Reply(test.statusCode)
	requestIntercept.ReplyFunc(func(r *gock.Response) {
		r.BodyString(test.body)
	})
	if test.isExpectedErr {
		requestIntercept.ReplyError(test.expectedErr)
	}
	httpClient := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(httpClient)

	return httpClient
}

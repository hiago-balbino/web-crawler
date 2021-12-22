package pager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentPage_Get(t *testing.T) {
	testCases := []struct {
		name          string
		uri           string
		isExpectedErr bool
	}{
		{
			name:          "should return error when http connection fails",
			uri:           "invalid",
			isExpectedErr: true,
		},
		{
			name:          "should return response when http connection has success",
			uri:           "https://google.com",
			isExpectedErr: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			response, err := NewContentPage().Get(test.uri)

			if test.isExpectedErr {
				assert.NotNil(t, err)
				assert.Nil(t, response)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response)
			}
		})
	}
}

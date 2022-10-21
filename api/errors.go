package api

import "errors"

type errMessage struct {
	Message string `json:"message"`
}

var (
	errEmptyURI   = errors.New("URI param cannot be empty")
	errEmptyDepth = errors.New("depth param cannot be empty")
)

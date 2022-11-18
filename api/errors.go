package api

import "errors"

var (
	errEmptyURI   = errors.New("URI param cannot be empty")
	errEmptyDepth = errors.New("depth param cannot be empty")
)

package page

import "net/http"

// Pager is a abstraction to handle with pages
type Pager interface {
	// Get fetch the response body from the URI
	Get(uri string) (*http.Response, error)
}

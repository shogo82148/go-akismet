package akismet

import "net/http"

// HTTPClient is a interface for http client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	// HTTPClient is underlying http client.
	// If it is nil, http.DefaultClient is used.
	HTTPClient HTTPClient

	// UserAgent is User-Agent header of requests.
	UserAgent string

	// APIKey is an API Key for using Akismet API.
	APIKey string
}

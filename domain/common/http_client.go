package common

import (
	"context"
	"net/http"
)

// HTTPClient defines an interface for sending HTTP requests and receiving HTTP responses.
// It includes methods for executing GET requests and generic request handling.
type HTTPClient interface {
	// Get sends a GET request to the specified URL and returns the response or an error.
	Get(url string) (*http.Response, error)
	// Post sends a POST request to the specified URL with the given content type and body,
	Post(url string, contentType string, body []byte) (*http.Response, error)
	// Do execute the provided HTTP request and returns the response or an error.
	Do(req *http.Request) (*http.Response, error)
	// DoWithContext executes the provided HTTP request with a context and returns the response or an error.
	DoWithContext(req *http.Request, ctx context.Context) (*http.Response, error)
}

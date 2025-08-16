package http_client

import (
	"bytes"
	"context"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"net/http"
	"time"
)

// StandardHttpClient implements the common.HTTPClient interface using the standard http.Client from the net/http package.
type StandardHttpClient struct {
	client *http.Client
}

// NewStandardHttpClient creates a new instance of StandardHttpClient with the specified timeout in seconds.
func NewStandardHttpClient(timeOut int) common.HTTPClient {
	return &StandardHttpClient{
		client: &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
	}
}

// Get performs a GET request to the specified URL and returns the response.
func (c *StandardHttpClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Post performs a POST request to the specified URL with the provided body and returns the response.
func (c *StandardHttpClient) Post(url string, contentType string, body []byte) (*http.Response, error) {
	return c.client.Post(url, contentType, bytes.NewReader(body))
}

// Do perform a request using the provided http.Request and returns the response.
func (c *StandardHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// DoWithContext performs a request using the provided http.Request with a context and returns the response.
func (c *StandardHttpClient) DoWithContext(req *http.Request, ctx context.Context) (*http.Response, error) {
	// The standard http.Client already supports context via the Request.WithContext method.
	return c.client.Do(req.WithContext(ctx))
}

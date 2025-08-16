package http_client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStandardHttpClient_Get(t *testing.T) {
	tests := []struct {
		name       string
		mockServer func() *httptest.Server
		wantErr    bool
	}{
		{
			name: "successful get",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			wantErr: false,
		},
		{
			name: "server error",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.mockServer()
			defer server.Close()

			client := &StandardHttpClient{client: &http.Client{}}
			_, err := client.Get(server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStandardHttpClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		mockServer func() *httptest.Server
		body       []byte
		wantErr    bool
	}{
		{
			name: "successful post",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			body:    []byte("test"),
			wantErr: false,
		},
		{
			name: "server error",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			body:    []byte("test"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.mockServer()
			defer server.Close()

			client := &StandardHttpClient{client: &http.Client{}}
			_, err := client.Post(server.URL, "application/json", tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStandardHttpClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		mockServer func() *httptest.Server
		request    func(url string) *http.Request
		wantErr    bool
	}{
		{
			name: "successful do",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			request: func(url string) *http.Request {
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantErr: false,
		},
		{
			name: "server error",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			request: func(url string) *http.Request {
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.mockServer()
			defer server.Close()

			client := &StandardHttpClient{client: &http.Client{}}
			_, err := client.Do(tt.request(server.URL))
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStandardHttpClient_DoWithContext(t *testing.T) {
	tests := []struct {
		name       string
		mockServer func() *httptest.Server
		request    func(url string) *http.Request
		wantErr    bool
		cancel     bool
	}{
		{
			name: "successful do with context",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			request: func(url string) *http.Request {
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantErr: false,
			cancel:  false,
		},
		{
			name: "context cancelled",
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					w.WriteHeader(http.StatusOK)
				}))
			},
			request: func(url string) *http.Request {
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantErr: true,
			cancel:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.mockServer()
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			if tt.cancel {
				cancel()
			} else {
				defer cancel()
			}

			client := &StandardHttpClient{client: &http.Client{}}
			_, err := client.DoWithContext(tt.request(server.URL), ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

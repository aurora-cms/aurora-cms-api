package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HttpClient struct {
	mock.Mock
}

func (m *HttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *HttpClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

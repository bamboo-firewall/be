package http

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, opts ...clientOption) *Client {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) NewRequest() *Request {
	return NewRequest(c)
}

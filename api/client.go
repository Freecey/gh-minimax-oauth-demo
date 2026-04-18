package api

import (
	"net/http"
)

// Client represents a basic HTTP client for API calls
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

// RESTClient represents a client for REST API calls
type RESTClient struct {
	*Client
}

// NewRESTClient creates a new REST client
func NewRESTClient() *RESTClient {
	return &RESTClient{
		Client: NewClient(),
	}
}
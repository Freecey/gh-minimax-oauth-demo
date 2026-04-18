package api

import (
	"net/http"
)

// Client represents an API client
type Client struct {
	httpClient *http.Client
	host       string
}

// NewClient creates a new API client
func NewClient(host string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
	}
}

// GQLClient represents a GraphQL client
type GQLClient struct {
	client *Client
}

// NewGQLClient creates a new GraphQL client
func NewGQLClient(host string) *GQLClient {
	return &GQLClient{
		client: NewClient(host),
	}
}

// Query represents a GraphQL query
type Query struct {
	Query     string
	Variables map[string]interface{}
}

// NewQuery creates a new GraphQL query
func NewQuery(query string) *Query {
	return &Query{
		Query:     query,
		Variables: make(map[string]interface{}),
	}
}

// Response represents a GraphQL response
type Response struct {
	Data   interface{}
	Errors []interface{}
}

// Do executes a GraphQL query
func (g *GQLClient) Do(q *Query) (*Response, error) {
	return &Response{
		Data:   map[string]interface{}{"result": "success"},
		Errors: []interface{}{},
	}, nil
}

// RESTClient represents a REST API client
type RESTClient struct {
	client *Client
}

// NewRESTClient creates a new REST client
func NewRESTClient(host string) *RESTClient {
	return &RESTClient{
		client: NewClient(host),
	}
}

// Do executes a REST request
func (r *RESTClient) Do(method, url string, body interface{}) (interface{}, error) {
	return map[string]interface{}{
		"method": method,
		"url":    url,
		"body":   body,
		"status": "ok",
	}, nil
}

// Request represents an HTTP request
type Request struct {
	Method string
	URL    string
	Body   interface{}
}

// NewRequest creates a new HTTP request
func NewRequest(method, url string, body interface{}) *Request {
	return &Request{
		Method: method,
		URL:    url,
		Body:   body,
	}
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// NewResponse creates a new HTTP response
func NewResponse(statusCode int, body interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
		Headers:    make(map[string]string),
	}
}
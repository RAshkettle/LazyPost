// Package models provides data structures for the LazyPost application.
// It contains models that represent HTTP requests and their components.
package models

// Request represents an HTTP request with all its components.
// It encapsulates all the data needed to make an HTTP request,
// including method, URL, headers and body.
type Request struct {
	Method  string            // HTTP method (GET, POST, PUT, DELETE, etc.)
	URL     string            // Target URL for the HTTP request
	Headers map[string]string // HTTP headers as key-value pairs
	Body    string            // Request body content (for POST, PUT, etc.)
}

// NewRequest creates a new HTTP request model with default values.
// By default, the request uses the GET method with no headers or body.
func NewRequest() Request {
	return Request{
		Method:  "GET",
		URL:     "",
		Headers: make(map[string]string),
		Body:    "",
	}
}

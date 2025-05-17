package ui

// RequestCompleteMsg is sent when an HTTP request has completed.
// It contains the response data from the request.
type RequestCompleteMsg struct {
	Headers string // Formatted headers string
	Body    string // Response body text
	Error   error  // Any error that occurred during the request
}

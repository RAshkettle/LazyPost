package ui

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// handleSubmit processes the form submission.
// It validates the URL, shows the loading spinner, and executes the request asynchronously.
// Returns a tea.Cmd if any needs to be executed.
func (a *App) handleSubmit() tea.Cmd {
	// Validate URL
	rawURL := a.urlInput.GetText()
	isValid := validateURL(rawURL)
	if !isValid {
		// Show a toast notification for invalid URL
		a.toast.Show("Invalid URL: The Provided URL is not valid.")
		// Move focus back to URL input if invalid
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(true)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(false)
		return nil
	}

	// Prepare for request - don't change focus yet
	a.methodSelector.SetActive(false)
	a.urlInput.SetActive(false)
	a.submitButton.SetActive(false)

	// Show the loading spinner directly over the URL input
	spinnerCmd := a.spinner.Show("Sending request...")

	// Get selected HTTP method
	method := a.methodSelector.GetSelectedMethod()

	// Get parameters from ParamsContainer via QueryTab
	// The GetQueryTab() method is now available on TabsContainer
	queryParams := a.tabContainer.GetQueryTab().ParamsInput.GetParams()
	finalURL, err := buildURLWithParams(rawURL, queryParams)
	if err != nil {
		// This error would typically be from parsing the rawURL, which should be caught by validateURL
		// but as a safeguard:
		a.toast.Show(fmt.Sprintf("Error building URL: %v", err))
		a.spinner.Hide()           // Hide spinner as we are not proceeding
		a.urlInput.SetActive(true) // Allow user to correct URL
		return nil
	}
	// Return a command that will execute the HTTP request asynchronously
	return tea.Batch(
		spinnerCmd,
		func() tea.Msg {
			// Create HTTP client
			client := &http.Client{}

			// Create request with the selected method and potentially modified URL
			req, err := http.NewRequest(method, finalURL, nil)
			if err != nil {
				return RequestCompleteMsg{
					Error: err,
				}
			}

			// Execute the HTTP request
			resp, err := client.Do(req)
			if err != nil {
				return RequestCompleteMsg{
					Error: err,
				}
			}

			defer resp.Body.Close()

			// Process response headers
			var headersContent strings.Builder

			// Add yellow and bold formatting for the "Status:" label
			headersContent.WriteString(fmt.Sprintf("\033[1;33mStatus:\033[0m %s\n\n", resp.Status))

			// Format each header with yellow and bold for the header name and colon
			for key, values := range resp.Header {
				for _, value := range values {
					headersContent.WriteString(fmt.Sprintf("\033[1;33m%s:\033[0m %s\n", key, value))
				}
			}

			// Process response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return RequestCompleteMsg{
					Error:   err,
					Headers: headersContent.String(),
				}
			}

			// Return the response data
			return RequestCompleteMsg{
				Headers: headersContent.String(),
				Body:    string(body),
			}
		},
	)
}

// buildURLWithParams takes a raw URL string and a map of query parameters,
// appends the parameters to the URL, and returns the modified URL string.
// It handles URL encoding for parameter names and values.
func buildURLWithParams(rawURL string, params map[string]string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	for name, value := range params {
		if strings.TrimSpace(name) != "" {
			query.Add(name, value) // url.Values.Add handles encoding internally for Add
		}
	}
	parsedURL.RawQuery = query.Encode() // Encode ensures correct formatting & escaping

	return parsedURL.String(), nil
}

package ui

import (
	"encoding/json" // Added import
	"regexp"
)

// validateURL checks if the provided string is a valid URL.
// It uses a regular expression to validate the URL format,
// ensuring it has the proper scheme, domain and optional path.
func validateURL(url string) bool {
	if url == "" {
		return false
	}

	// Parse the URL to reject URLs with unencoded spaces
	// While also allowing valid components like:
	// - HTTP and HTTPS protocols only
	// - Domain names with hyphens (including consecutive hyphens)
	// - Valid TLDs (2 or more characters)
	// - Optional port numbers (1-5 digits, limited to 0-65535)
	// - Optional path (no unencoded spaces)
	// - Optional query parameters
	// - Optional fragments

	// First, check for spaces in the URL (except in percent-encoded form)
	if regexp.MustCompile(`\s`).MatchString(url) {
		return false
	}

	// Basic URL regex pattern without space validation
	pattern := `^(http|https)://[a-zA-Z0-9]+([-\.][a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}(:[0-9]{1,5})?(\/[^?#]*)?(\?[^#]*)?(#.*)?$`
	matched, _ := regexp.MatchString(pattern, url)
	if !matched {
		return false
	}

	// Additional validation for port numbers (should be 0-65535)
	portPattern := `:([0-9]+)`
	portRegex := regexp.MustCompile(portPattern)
	portMatches := portRegex.FindStringSubmatch(url)

	if len(portMatches) > 1 {
		// We found a port number, check if it's valid
		port := portMatches[1]
		if len(port) > 5 {
			return false
		}

		// Convert port to integer for proper comparison
		portNum := 0
		for _, digit := range port {
			portNum = portNum*10 + int(digit-'0')
		}

		if portNum > 65535 {
			return false
		}
	}

	return true
}

// IsValidJSON checks if the provided string is valid JSON.
func IsValidJSON(s string) bool {
	// If the string is empty, it can be considered valid JSON (e.g., an empty object or array, or just empty).
	// However, json.Valid considers an empty string as invalid.
	// For the purpose of a request body, an empty string is often a valid case (no body).
	// If specific JSON (like {} or []) is required for an empty body, this logic might change.
	// For now, let's assume an empty string means "no body" and is valid in that context.
	// If the user intends to send JSON, it should not be an empty string unless it's "null", "{}", "[]", etc.
	// Let's stick to strict JSON validation: an empty string is not valid JSON.
	if s == "" {
		return true // Assuming an empty body is acceptable if no JSON content is provided.
		             // If strict JSON is always required, this should be false or handled upstream.
	}
	return json.Valid([]byte(s))
}

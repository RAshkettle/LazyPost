package ui

import (
	"testing"
)

// TestValidateURL tests the validateURL function with a variety of common valid and invalid URL formats.
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		// Valid URLs
		{
			name:     "Simple HTTP URL",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "Simple HTTPS URL",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "URL with subdomain",
			url:      "https://blog.example.com",
			expected: true,
		},
		{
			name:     "URL with path",
			url:      "https://example.com/path",
			expected: true,
		},
		{
			name:     "URL with path segments",
			url:      "https://example.com/path/to/resource",
			expected: true,
		},
		{
			name:     "URL with port",
			url:      "https://example.com:8080",
			expected: true,
		},
		{
			name:     "URL with port and path",
			url:      "https://example.com:8080/path",
			expected: true,
		},
		{
			name:     "URL with numbers in domain",
			url:      "https://example123.com",
			expected: true,
		},
		{
			name:     "URL with query parameters",
			url:      "https://example.com?param=value",
			expected: true,
		},
		{
			name:     "URL with multiple query parameters",
			url:      "https://example.com?param1=value1&param2=value2",
			expected: true,
		},
		{
			name:     "URL with fragment",
			url:      "https://example.com#section",
			expected: true,
		},
		{
			name:     "Complex URL with all components",
			url:      "https://sub.example123.com:8443/path/to/resource?param=value#section",
			expected: true,
		},

		// Invalid URLs
		{
			name:     "Empty string",
			url:      "",
			expected: false,
		},
		{
			name:     "Missing protocol",
			url:      "example.com",
			expected: false,
		},
		{
			name:     "Invalid protocol",
			url:      "ftp://example.com",
			expected: false,
		},
		{
			name:     "IP address instead of domain",
			url:      "http://192.168.1.1",
			expected: false, // IP addresses not supported by the current regex
		},
		{
			name:     "Missing domain",
			url:      "http://",
			expected: false,
		},
		{
			name:     "Invalid TLD (too short)",
			url:      "http://example.c",
			expected: false,
		},
		{
			name:     "Local domain without TLD",
			url:      "http://localhost",
			expected: false, // Local domains not supported by the current regex
		},
		{
			name:     "Invalid characters in domain",
			url:      "http://ex@mple.com",
			expected: false,
		},
		{
			name:     "URL with spaces",
			url:      "http://example.com/path with spaces",
			expected: false, // Spaces not encoded
		},
		{
			name:     "Malformed URL",
			url:      "http:/example.com",
			expected: false,
		},
		{
			name:     "Protocol only",
			url:      "http://",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validateURL(test.url)
			if result != test.expected {
				t.Errorf("validateURL(%q) = %v, expected %v", test.url, result, test.expected)
			}
		})
	}
}

// TestURLEdgeCases tests the validateURL function with edge cases and less common URL formats.
func TestURLEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "URL with valid percent encoding",
			url:      "https://example.com/path%20with%20spaces",
			expected: true,
		},
		{
			name:     "URL with hyphenated domain",
			url:      "https://my-website-domain.com",
			expected: true,
		},
		{
			name:     "URL with very long domain",
			url:      "https://this-is-a-very-long-domain-name-that-is-valid-and-should-be-accepted-by-the-regex.com",
			expected: true,
		},
		{
			name:     "URL with unusually long TLD",
			url:      "https://example.international",
			expected: true,
		},
		{
			name:     "URL with double hyphens in domain",
			url:      "https://example--domain.com",
			expected: true,
		},
		{
			name:     "URL with underscores",
			url:      "https://example_domain.com",
			expected: false, // Underscores not allowed in hostname
		},
		{
			name:     "URL with invalid port (too large)",
			url:      "https://example.com:99999",
			expected: false, // Ports must be 1-5 digits (up to 65535)
		},
		{
			name:     "URL with invalid port (non-numeric)",
			url:      "https://example.com:port",
			expected: false,
		},
		{
			name:     "URL with username and password",
			url:      "https://user:password@example.com",
			expected: false, // Not supported by the current regex
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validateURL(test.url)
			if result != test.expected {
				t.Errorf("validateURL(%q) = %v, expected %v", test.url, result, test.expected)
			}
		})
	}
}

// TestURLPerformance runs a basic performance check on the validateURL function.
// It can be skipped in short mode.
func TestURLPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	url := "https://www.example.com/path/to/resource?param1=value1&param2=value2#section"
	iterations := 10000

	t.Run("Performance", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			validateURL(url)
		}
		// No assertions here, this just validates that the function can handle
		// a large number of calls without issues
	})
}

// TestIsValidJSON tests the IsValidJSON function with various JSON strings.
func TestIsValidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected bool
	}{
		// Valid JSON
		{
			name:     "Empty string (considered valid as per current function logic)",
			jsonStr:  "",
			expected: true,
		},
		{
			name:     "Valid empty object",
			jsonStr:  "{}",
			expected: true,
		},
		{
			name:     "Valid empty array",
			jsonStr:  "[]",
			expected: true,
		},
		{
			name:     "Valid simple object",
			jsonStr:  `{"key": "value"}`,
			expected: true,
		},
		{
			name:     "Valid object with multiple keys",
			jsonStr:  `{"key1": "value1", "key2": 123, "key3": true}`,
			expected: true,
		},
		{
			name:     "Valid array of strings",
			jsonStr:  `["apple", "banana", "cherry"]`,
			expected: true,
		},
		{
			name:     "Valid array of mixed types",
			jsonStr:  `["string", 100, false, null, {}, []]`,
			expected: true,
		},
		{
			name:     "Valid JSON string",
			jsonStr:  `"hello world"`,
			expected: true,
		},
		{
			name:     "Valid JSON number (integer)",
			jsonStr:  "12345",
			expected: true,
		},
		{
			name:     "Valid JSON number (float)",
			jsonStr:  "123.45",
			expected: true,
		},
		{
			name:     "Valid JSON boolean true",
			jsonStr:  "true",
			expected: true,
		},
		{
			name:     "Valid JSON boolean false",
			jsonStr:  "false",
			expected: true,
		},
		{
			name:     "Valid JSON null",
			jsonStr:  "null",
			expected: true,
		},
		{
			name:     "Valid nested object",
			jsonStr:  `{"name": "John Doe", "age": 30, "address": {"street": "123 Main St", "city": "Anytown"}}`,
			expected: true,
		},
		{
			name:     "Valid object with array value",
			jsonStr:  `{"name": "Fruits", "items": ["apple", "orange"]}`,
			expected: true,
		},

		// Invalid JSON
		{
			name:     "Malformed object (missing closing brace)",
			jsonStr:  `{"key": "value"`,
			expected: false,
		},
		{
			name:     "Malformed object (trailing comma)",
			jsonStr:  `{"key1": "value1",}`,
			expected: false,
		},
		{
			name:     "Malformed array (missing closing bracket)",
			jsonStr:  `["apple", "banana"`,
			expected: false,
		},
		{
			name:     "Malformed array (trailing comma)",
			jsonStr:  `["apple", "banana",]`,
			expected: false,
		},
		{
			name:     "Invalid: Single quotes for strings",
			jsonStr:  `{'key': 'value'}`,
			expected: false,
		},
		{
			name:     "Invalid: Unquoted key",
			jsonStr:  `{key: "value"}`,
			expected: false,
		},
		{
			name:     "Invalid: Plain string not in quotes (unless it's true, false, null, or a number)",
			jsonStr:  "not a valid top-level json string without quotes unless keyword",
			expected: false,
		},
		{
			name:     "Invalid: Number with leading zero (not single zero)",
			jsonStr:  "0123",
			expected: false,
		},
		{
			name:     "Valid: Single zero",
			jsonStr:  "0",
			expected: true,
		},
		{
			name:     "Invalid: Hex number",
			jsonStr:  "0x10",
			expected: false,
		},
		{
			name:     "Invalid: NaN",
			jsonStr:  "NaN",
			expected: false,
		},
		{
			name:     "Invalid: Infinity",
			jsonStr:  "Infinity",
			expected: false,
		},
		{
			name:     "Invalid: Just a comma",
			jsonStr:  ",",
			expected: false,
		},
		{
			name:     "Invalid: Just a colon",
			jsonStr:  ":",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidJSON(tt.jsonStr); got != tt.expected {
				t.Errorf("IsValidJSON(%q) = %v, want %v", tt.jsonStr, got, tt.expected)
			}
		})
	}
}

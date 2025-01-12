package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TODO: improve this (make more like APIClient)

// Fetch makes an HTTP request and unmarshals the response body into the provided result interface.
// Parameters:
// - method: HTTP method (e.g., "GET", "POST").
// - url: API endpoint.
// - token: Optional Bearer token for authentication (pass "" if not needed).
// - body: Optional request body (pass nil for no body).
// - result: Pointer to a variable where the JSON response will be unmarshalled.
func Fetch(method, url, token string, body interface{}, result interface{}, headers map[string]string) error {
	// Serialize the body to JSON if provided
	var requestBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonData)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Add custom headers if provided
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Make the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Check for non-OK HTTP status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read and unmarshal the response body
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		err = json.Unmarshal(responseData, result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}

func AddQueryParams(baseURL string, params map[string]string) (string, error) {
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Create URL query parameters from the map
	query := parsedURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	// Attach the query parameters to the base URL and return it
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

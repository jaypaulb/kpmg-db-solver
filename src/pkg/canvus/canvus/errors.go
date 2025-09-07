// Package canvus provides error types for the Canvus SDK.
package canvus

import "fmt"

// APIError represents an error returned by the Canvus API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

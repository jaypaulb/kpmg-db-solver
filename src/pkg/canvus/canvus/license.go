package canvus

import (
	"context"
	"fmt"
)

// LicenseInfo represents the license information for the server.
type LicenseInfo struct {
	Key       string   `json:"key,omitempty"`
	Valid     bool     `json:"valid,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	Type      string   `json:"type,omitempty"`
	Seats     int      `json:"seats,omitempty"`
	IssuedTo  string   `json:"issued_to,omitempty"`
	IssuedBy  string   `json:"issued_by,omitempty"`
	Features  []string `json:"features,omitempty"`
	// Add other fields as needed based on the API response
}

// GetLicenseInfo retrieves the current license information from the Canvus API.
func (s *Session) GetLicenseInfo(ctx context.Context) (*LicenseInfo, error) {
	var info LicenseInfo
	err := s.doRequest(ctx, "GET", "license", nil, &info, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetLicenseInfo: %w", err)
	}
	return &info, nil
}

package canvus

import (
	"context"
	"fmt"
)

// ServerInfo represents information about the Canvus Server instance.
type ServerInfo struct {
	API      []string `json:"api"`
	Go       string   `json:"go"`
	ServerID string   `json:"server_id"`
	Version  string   `json:"version"`
}

// GetServerInfo retrieves information about the Canvus Server instance from the Canvus API.
func (s *Session) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	var info ServerInfo
	err := s.doRequest(ctx, "GET", "server-info", nil, &info, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetServerInfo: %w", err)
	}
	return &info, nil
}

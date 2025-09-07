package canvus

import (
	"context"
	"fmt"
)

// VideoInputSource represents a video input source for a client device.
type VideoInputSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListVideoInputs retrieves all video input widgets for a given canvas.
func (s *Session) ListVideoInputs(ctx context.Context, canvasID string) ([]VideoInput, error) {
	var inputs []VideoInput
	path := fmt.Sprintf("canvases/%s/video-inputs", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &inputs, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListVideoInputs: %w", err)
	}
	return inputs, nil
}

// CreateVideoInput creates a new video input widget on a canvas. The payload must include 'source' and 'host-id'.
func (s *Session) CreateVideoInput(ctx context.Context, canvasID string, req interface{}) (*VideoInput, error) {
	var input VideoInput
	path := fmt.Sprintf("canvases/%s/video-inputs", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateVideoInput: %w", err)
	}
	return &input, nil
}

// DeleteVideoInput deletes a video input widget by ID for a given canvas.
func (s *Session) DeleteVideoInput(ctx context.Context, canvasID, inputID string) error {
	path := fmt.Sprintf("canvases/%s/video-inputs/%s", canvasID, inputID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// ListClientVideoInputs retrieves all video input sources for a given client device.
func (s *Session) ListClientVideoInputs(ctx context.Context, clientID string) ([]VideoInputSource, error) {
	var sources []VideoInputSource
	path := fmt.Sprintf("clients/%s/video-inputs", clientID)
	err := s.doRequest(ctx, "GET", path, nil, &sources, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListClientVideoInputs: %w", err)
	}
	return sources, nil
}

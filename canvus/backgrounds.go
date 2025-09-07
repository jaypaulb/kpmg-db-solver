package canvus

import (
	"context"
	"fmt"
)

// GetCanvasBackground retrieves the background settings for a specified canvas.
func (s *Session) GetCanvasBackground(ctx context.Context, canvasID string) (*CanvasBackground, error) {
	var bg CanvasBackground
	path := fmt.Sprintf("canvases/%s/background", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &bg, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetCanvasBackground: %w", err)
	}
	return &bg, nil
}

// PatchCanvasBackground sets the background settings for a specified canvas (solid color or haze).
func (s *Session) PatchCanvasBackground(ctx context.Context, canvasID string, req interface{}) error {
	path := fmt.Sprintf("canvases/%s/background", canvasID)
	return s.doRequest(ctx, "PATCH", path, req, nil, nil, false)
}

// PostCanvasBackground sets the background to an image for a specified canvas. The request must be a multipart POST with a 'data' part and optional 'json' part.
func (s *Session) PostCanvasBackground(ctx context.Context, canvasID string, multipartBody interface{}) error {
	path := fmt.Sprintf("canvases/%s/background", canvasID)
	return s.doRequest(ctx, "POST", path, multipartBody, nil, nil, false)
}

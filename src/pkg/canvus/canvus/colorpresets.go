package canvus

import (
	"context"
	"fmt"
)

// ListColorPresets retrieves all color presets for a given canvas.
func (s *Session) ListColorPresets(ctx context.Context, canvasID string) ([]ColorPreset, error) {
	var presets []ColorPreset
	path := fmt.Sprintf("canvases/%s/colorpresets", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &presets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListColorPresets: %w", err)
	}
	return presets, nil
}

// GetColorPreset retrieves a color preset by name for a given canvas.
func (s *Session) GetColorPreset(ctx context.Context, canvasID, name string) (*ColorPreset, error) {
	var preset ColorPreset
	path := fmt.Sprintf("canvases/%s/colorpresets/%s", canvasID, name)
	err := s.doRequest(ctx, "GET", path, nil, &preset, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetColorPreset: %w", err)
	}
	return &preset, nil
}

// CreateColorPreset creates a new color preset on a canvas.
func (s *Session) CreateColorPreset(ctx context.Context, canvasID string, req interface{}) (*ColorPreset, error) {
	var preset ColorPreset
	path := fmt.Sprintf("canvases/%s/colorpresets", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &preset, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateColorPreset: %w", err)
	}
	return &preset, nil
}

// UpdateColorPreset updates a color preset by name.
func (s *Session) UpdateColorPreset(ctx context.Context, canvasID, name string, req interface{}) (*ColorPreset, error) {
	var preset ColorPreset
	path := fmt.Sprintf("canvases/%s/colorpresets/%s", canvasID, name)
	err := s.doRequest(ctx, "PATCH", path, req, &preset, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateColorPreset: %w", err)
	}
	return &preset, nil
}

// DeleteColorPreset deletes a color preset by name.
func (s *Session) DeleteColorPreset(ctx context.Context, canvasID, name string) error {
	path := fmt.Sprintf("canvases/%s/colorpresets/%s", canvasID, name)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// GetColorPresets retrieves the color presets for the specified canvas.
func (s *Session) GetColorPresets(ctx context.Context, canvasID string) (*ColorPresets, error) {
	var presets ColorPresets
	path := fmt.Sprintf("api/v1/canvases/%s/color-presets", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &presets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetColorPresets: %w", err)
	}
	return &presets, nil
}

// PatchColorPresets updates the color presets for the specified canvas.
func (s *Session) PatchColorPresets(ctx context.Context, canvasID string, req *ColorPresets) (*ColorPresets, error) {
	var updated ColorPresets
	path := fmt.Sprintf("api/v1/canvases/%s/color-presets", canvasID)
	err := s.doRequest(ctx, "PATCH", path, req, &updated, nil, false)
	if err != nil {
		return nil, fmt.Errorf("PatchColorPresets: %w", err)
	}
	return &updated, nil
}

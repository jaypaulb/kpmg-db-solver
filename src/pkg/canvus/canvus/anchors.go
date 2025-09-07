package canvus

import (
	"context"
	"fmt"
)

// ListAnchors retrieves all anchors for a given canvas.
func (s *Session) ListAnchors(ctx context.Context, canvasID string) ([]Anchor, error) {
	var anchors []Anchor
	path := fmt.Sprintf("canvases/%s/anchors", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &anchors, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListAnchors: %w", err)
	}
	return anchors, nil
}

// GetAnchor retrieves an anchor by ID for a given canvas.
func (s *Session) GetAnchor(ctx context.Context, canvasID, anchorID string) (*Anchor, error) {
	var anchor Anchor
	path := fmt.Sprintf("canvases/%s/anchors/%s", canvasID, anchorID)
	err := s.doRequest(ctx, "GET", path, nil, &anchor, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetAnchor: %w", err)
	}
	return &anchor, nil
}

// CreateAnchor creates a new anchor on a canvas.
func (s *Session) CreateAnchor(ctx context.Context, canvasID string, req interface{}) (*Anchor, error) {
	var anchor Anchor
	path := fmt.Sprintf("canvases/%s/anchors", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &anchor, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateAnchor: %w", err)
	}
	return &anchor, nil
}

// UpdateAnchor updates an anchor by ID for a given canvas.
func (s *Session) UpdateAnchor(ctx context.Context, canvasID, anchorID string, req interface{}) (*Anchor, error) {
	var anchor Anchor
	path := fmt.Sprintf("canvases/%s/anchors/%s", canvasID, anchorID)
	err := s.doRequest(ctx, "PATCH", path, req, &anchor, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateAnchor: %w", err)
	}
	return &anchor, nil
}

// DeleteAnchor deletes an anchor by ID for a given canvas.
func (s *Session) DeleteAnchor(ctx context.Context, canvasID, anchorID string) error {
	path := fmt.Sprintf("canvases/%s/anchors/%s", canvasID, anchorID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

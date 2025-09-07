package canvus

import (
	"context"
	"fmt"
)

// ListConnectors retrieves all connectors for a given canvas.
func (s *Session) ListConnectors(ctx context.Context, canvasID string) ([]Connector, error) {
	var connectors []Connector
	path := fmt.Sprintf("canvases/%s/connectors", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &connectors, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListConnectors: %w", err)
	}
	return connectors, nil
}

// GetConnector retrieves a connector by ID for a given canvas.
func (s *Session) GetConnector(ctx context.Context, canvasID, connectorID string) (*Connector, error) {
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	err := s.doRequest(ctx, "GET", path, nil, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetConnector: %w", err)
	}
	return &connector, nil
}

// CreateConnector creates a new connector on a canvas.
// If req["src"] or req["dst"] is a map (widget JSON), the widget is created first and its ID is used.
func (s *Session) CreateConnector(ctx context.Context, canvasID string, req interface{}) (*Connector, error) {
	m, ok := req.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("CreateConnector: req must be a map[string]interface{}")
	}
	// Helper to create widget if needed
	resolveEnd := func(key string) (string, error) {
		v, ok := m[key]
		if !ok {
			return "", fmt.Errorf("CreateConnector: missing %s", key)
		}
		// If already a string, treat as ID
		if id, ok := v.(string); ok {
			return id, nil
		}
		// If map, create widget
		if widgetData, ok := v.(map[string]interface{}); ok {
			widget, err := s.CreateWidget(ctx, canvasID, widgetData)
			if err != nil {
				return "", fmt.Errorf("CreateConnector: failed to create widget for %s: %w", key, err)
			}
			return widget.ID, nil
		}
		return "", fmt.Errorf("CreateConnector: %s must be string or widget JSON", key)
	}
	// Resolve src
	srcID, err := resolveEnd("src")
	if err != nil {
		return nil, err
	}
	// Resolve dst
	dstID, err := resolveEnd("dst")
	if err != nil {
		return nil, err
	}
	// Build connector request
	m["src"] = map[string]interface{}{"id": srcID}
	m["dst"] = map[string]interface{}{"id": dstID}
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors", canvasID)
	err = s.doRequest(ctx, "POST", path, m, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateConnector: %w", err)
	}
	return &connector, nil
}

// UpdateConnector updates a connector by ID for a given canvas.
func (s *Session) UpdateConnector(ctx context.Context, canvasID, connectorID string, req interface{}) (*Connector, error) {
	var connector Connector
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	err := s.doRequest(ctx, "PATCH", path, req, &connector, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateConnector: %w", err)
	}
	return &connector, nil
}

// DeleteConnector deletes a connector by ID for a given canvas.
func (s *Session) DeleteConnector(ctx context.Context, canvasID, connectorID string) error {
	path := fmt.Sprintf("canvases/%s/connectors/%s", canvasID, connectorID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

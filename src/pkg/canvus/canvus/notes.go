package canvus

import (
	"context"
	"fmt"
)

// ListNotes retrieves all notes for a given canvas.
func (s *Session) ListNotes(ctx context.Context, canvasID string) ([]Note, error) {
	var notes []Note
	path := fmt.Sprintf("canvases/%s/notes", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &notes, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListNotes: %w", err)
	}
	return notes, nil
}

// GetNote retrieves a note by ID for a given canvas.
func (s *Session) GetNote(ctx context.Context, canvasID, noteID string) (*Note, error) {
	var note Note
	path := fmt.Sprintf("canvases/%s/notes/%s", canvasID, noteID)
	err := s.doRequest(ctx, "GET", path, nil, &note, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetNote: %w", err)
	}
	return &note, nil
}

// CreateNote creates a new note on a canvas.
func (s *Session) CreateNote(ctx context.Context, canvasID string, req interface{}) (*Note, error) {
	var note Note
	path := fmt.Sprintf("canvases/%s/notes", canvasID)
	err := s.doRequest(ctx, "POST", path, req, &note, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateNote: %w", err)
	}
	return &note, nil
}

// UpdateNote updates a note by ID for a given canvas.
func (s *Session) UpdateNote(ctx context.Context, canvasID, noteID string, req interface{}) (*Note, error) {
	var note Note
	path := fmt.Sprintf("canvases/%s/notes/%s", canvasID, noteID)
	err := s.doRequest(ctx, "PATCH", path, req, &note, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateNote: %w", err)
	}
	return &note, nil
}

// DeleteNote deletes a note by ID for a given canvas.
func (s *Session) DeleteNote(ctx context.Context, canvasID, noteID string) error {
	path := fmt.Sprintf("canvases/%s/notes/%s", canvasID, noteID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

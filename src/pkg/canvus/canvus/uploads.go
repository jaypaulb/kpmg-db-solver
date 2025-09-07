package canvus

import (
	"context"
	"fmt"
)

// UploadNote uploads a note to the uploads folder of a canvas. The request must be a multipart POST with a 'json' part.
func (s *Session) UploadNote(ctx context.Context, canvasID string, multipartBody interface{}) (*Note, error) {
	var note Note
	path := fmt.Sprintf("canvases/%s/uploads-folder", canvasID)
	err := s.doRequest(ctx, "POST", path, multipartBody, &note, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UploadNote: %w", err)
	}
	return &note, nil
}

// UploadAsset uploads a file asset to the uploads folder of a canvas. The request must be a multipart POST with a 'data' part and optional 'json' part.
func (s *Session) UploadAsset(ctx context.Context, canvasID string, multipartBody interface{}) (*Asset, error) {
	var asset Asset
	path := fmt.Sprintf("canvases/%s/uploads-folder", canvasID)
	err := s.doRequest(ctx, "POST", path, multipartBody, &asset, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UploadAsset: %w", err)
	}
	return &asset, nil
}

// ListUploadsFolder attempts to get the contents of the uploads folder for a canvas.
// This endpoint may not be documented but could exist.
func (s *Session) ListUploadsFolder(ctx context.Context, canvasID string) ([]Widget, error) {
	var widgets []Widget
	path := fmt.Sprintf("canvases/%s/uploads-folder", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &widgets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListUploadsFolder: %w", err)
	}
	return widgets, nil
}

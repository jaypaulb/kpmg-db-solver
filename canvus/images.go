package canvus

import (
	"context"
	"fmt"
	"io"
)

// ListImages retrieves all images for a given canvas.
func (s *Session) ListImages(ctx context.Context, canvasID string) ([]Image, error) {
	var images []Image
	path := fmt.Sprintf("canvases/%s/images", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &images, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListImages: %w", err)
	}
	return images, nil
}

// GetImage retrieves an image by ID for a given canvas.
func (s *Session) GetImage(ctx context.Context, canvasID, imageID string) (*Image, error) {
	var image Image
	path := fmt.Sprintf("canvases/%s/images/%s", canvasID, imageID)
	err := s.doRequest(ctx, "GET", path, nil, &image, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetImage: %w", err)
	}
	return &image, nil
}

// CreateImage creates a new image on a canvas. This must be a multipart POST with a 'json' and 'data' part.
func (s *Session) CreateImage(ctx context.Context, canvasID string, multipartBody io.Reader, contentType string) (*Image, error) {
	var image Image
	path := fmt.Sprintf("canvases/%s/images", canvasID)
	err := s.doRequest(ctx, "POST", path, multipartBody, &image, nil, false, contentType)
	if err != nil {
		return nil, fmt.Errorf("CreateImage: %w", err)
	}
	return &image, nil
}

// UpdateImage updates an image by ID for a given canvas.
func (s *Session) UpdateImage(ctx context.Context, canvasID, imageID string, req interface{}) (*Image, error) {
	var image Image
	path := fmt.Sprintf("canvases/%s/images/%s", canvasID, imageID)
	err := s.doRequest(ctx, "PATCH", path, req, &image, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateImage: %w", err)
	}
	return &image, nil
}

// DeleteImage deletes an image by ID for a given canvas.
func (s *Session) DeleteImage(ctx context.Context, canvasID, imageID string) error {
	path := fmt.Sprintf("canvases/%s/images/%s", canvasID, imageID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// DownloadImage downloads an image file by ID for a given canvas.
func (s *Session) DownloadImage(ctx context.Context, canvasID, imageID string) ([]byte, error) {
	path := fmt.Sprintf("canvases/%s/images/%s/download", canvasID, imageID)
	var data []byte
	err := s.doRequest(ctx, "GET", path, nil, &data, nil, true)
	if err != nil {
		return nil, fmt.Errorf("DownloadImage: %w", err)
	}
	return data, nil
}

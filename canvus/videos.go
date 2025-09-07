package canvus

import (
	"context"
	"fmt"
)

// ListVideos retrieves all videos for a given canvas.
func (s *Session) ListVideos(ctx context.Context, canvasID string) ([]Video, error) {
	var videos []Video
	path := fmt.Sprintf("canvases/%s/videos", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &videos, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListVideos: %w", err)
	}
	return videos, nil
}

// GetVideo retrieves a single video by ID for a given canvas.
func (s *Session) GetVideo(ctx context.Context, canvasID, videoID string) (*Video, error) {
	var video Video
	path := fmt.Sprintf("canvases/%s/videos/%s", canvasID, videoID)
	err := s.doRequest(ctx, "GET", path, nil, &video, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetVideo: %w", err)
	}
	return &video, nil
}

// DownloadVideo downloads a video file by ID for a given canvas.
func (s *Session) DownloadVideo(ctx context.Context, canvasID, videoID string) ([]byte, error) {
	path := fmt.Sprintf("canvases/%s/videos/%s/download", canvasID, videoID)
	var data []byte
	err := s.doRequest(ctx, "GET", path, nil, &data, nil, true)
	if err != nil {
		return nil, fmt.Errorf("DownloadVideo: %w", err)
	}
	return data, nil
}

// CreateVideo creates a new video on a canvas. This must be a multipart POST with a 'json' and 'data' part.
func (s *Session) CreateVideo(ctx context.Context, canvasID string, multipartBody interface{}, contentType string) (*Video, error) {
	var video Video
	path := fmt.Sprintf("canvases/%s/videos", canvasID)
	err := s.doRequest(ctx, "POST", path, multipartBody, &video, nil, false, contentType)
	if err != nil {
		return nil, fmt.Errorf("CreateVideo: %w", err)
	}
	return &video, nil
}

// UpdateVideo updates a video by ID for a given canvas.
func (s *Session) UpdateVideo(ctx context.Context, canvasID, videoID string, req interface{}) (*Video, error) {
	var video Video
	path := fmt.Sprintf("canvases/%s/videos/%s", canvasID, videoID)
	err := s.doRequest(ctx, "PATCH", path, req, &video, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateVideo: %w", err)
	}
	return &video, nil
}

// DeleteVideo deletes a video by ID for a given canvas.
func (s *Session) DeleteVideo(ctx context.Context, canvasID, videoID string) error {
	path := fmt.Sprintf("canvases/%s/videos/%s", canvasID, videoID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

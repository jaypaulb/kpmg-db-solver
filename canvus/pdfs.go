package canvus

import (
	"context"
	"fmt"
)

// ListPDFs retrieves all PDFs for a given canvas.
func (s *Session) ListPDFs(ctx context.Context, canvasID string) ([]PDF, error) {
	var pdfs []PDF
	path := fmt.Sprintf("canvases/%s/pdfs", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &pdfs, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListPDFs: %w", err)
	}
	return pdfs, nil
}

// GetPDF retrieves a single PDF by ID for a given canvas.
func (s *Session) GetPDF(ctx context.Context, canvasID, pdfID string) (*PDF, error) {
	var pdf PDF
	path := fmt.Sprintf("canvases/%s/pdfs/%s", canvasID, pdfID)
	err := s.doRequest(ctx, "GET", path, nil, &pdf, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetPDF: %w", err)
	}
	return &pdf, nil
}

// DownloadPDF downloads a PDF file by ID for a given canvas.
func (s *Session) DownloadPDF(ctx context.Context, canvasID, pdfID string) ([]byte, error) {
	path := fmt.Sprintf("canvases/%s/pdfs/%s/download", canvasID, pdfID)
	var data []byte
	err := s.doRequest(ctx, "GET", path, nil, &data, nil, true)
	if err != nil {
		return nil, fmt.Errorf("DownloadPDF: %w", err)
	}
	return data, nil
}

// CreatePDF creates a new PDF on a canvas. This must be a multipart POST with a 'json' and 'data' part.
func (s *Session) CreatePDF(ctx context.Context, canvasID string, multipartBody interface{}, contentType string) (*PDF, error) {
	var pdf PDF
	path := fmt.Sprintf("canvases/%s/pdfs", canvasID)
	err := s.doRequest(ctx, "POST", path, multipartBody, &pdf, nil, false, contentType)
	if err != nil {
		return nil, fmt.Errorf("CreatePDF: %w", err)
	}
	return &pdf, nil
}

// UpdatePDF updates a PDF by ID for a given canvas.
func (s *Session) UpdatePDF(ctx context.Context, canvasID, pdfID string, req interface{}) (*PDF, error) {
	var pdf PDF
	path := fmt.Sprintf("canvases/%s/pdfs/%s", canvasID, pdfID)
	err := s.doRequest(ctx, "PATCH", path, req, &pdf, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdatePDF: %w", err)
	}
	return &pdf, nil
}

// DeletePDF deletes a PDF by ID for a given canvas.
func (s *Session) DeletePDF(ctx context.Context, canvasID, pdfID string) error {
	path := fmt.Sprintf("canvases/%s/pdfs/%s", canvasID, pdfID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

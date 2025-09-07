package canvus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportedWidgetSet represents a set of widgets and assets exported from a canvas.
type ExportedWidgetSet struct {
	Widgets []Widget          // All exported widgets
	Assets  map[string]string // Asset filename, keyed by widget ID
	Region  *Rectangle        // The region used for export (nil for ExportByID)
}

// ExportWidgetsToFolder exports the specified widgets (and their assets) to a folder. Returns the export folder path.
// Accepts sharedCanvasID to blank parent_id for widgets whose parent is the shared canvas.
func (s *Session) ExportWidgetsToFolder(ctx context.Context, canvasID string, widgetIDs []string, region Rectangle, sharedCanvasID string, baseFolder string) (string, error) {
	if baseFolder == "" {
		baseFolder = filepath.Join("export", time.Now().Format("20060102_150405"))
	}
	exportFolder := baseFolder
	if err := os.MkdirAll(exportFolder, 0755); err != nil {
		return "", err
	}
	var selected []Widget
	assets := make(map[string]string)
	for _, id := range widgetIDs {
		w, err := s.GetWidget(ctx, canvasID, id)
		if err != nil {
			fmt.Printf("[EXPORT] Failed to get widget %s: %v\n", id, err)
			return "", fmt.Errorf("ExportWidgetsToFolder: failed to get widget %s: %w", id, err)
		}
		fmt.Printf("[EXPORT] Processing widget ID=%s, Type=%s\n", w.ID, w.WidgetType)
		// Blank parent_id if it matches sharedCanvasID
		if sharedCanvasID != "" && w.ParentID == sharedCanvasID {
			w.ParentID = ""
		}
		selected = append(selected, *w)
		widgetTypeLower := w.WidgetType
		if widgetTypeLower != "" {
			widgetTypeLower = strings.ToLower(widgetTypeLower)
		}
		switch widgetTypeLower {
		case "image":
			fmt.Printf("[EXPORT] Attempting to export image asset for widget %s\n", w.ID)
			img, err := s.GetImage(ctx, canvasID, w.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to get image for widget %s: %v\n", w.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to get image %s: %w", w.ID, err)
			}
			data, err := s.DownloadImage(ctx, canvasID, img.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to download image asset %s: %v\n", img.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to download image asset %s: %w", img.ID, err)
			}
			filename := "image_" + w.ID + ".jpg"
			filePath := filepath.Join(exportFolder, filename)
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				fmt.Printf("[EXPORT] Failed to write image file %s: %v\n", filePath, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to write image file: %w", err)
			}
			fmt.Printf("[EXPORT] Wrote image file: %s\n", filePath)
			assets[w.ID] = filename
		case "pdf":
			fmt.Printf("[EXPORT] Attempting to export pdf asset for widget %s\n", w.ID)
			pdf, err := s.GetPDF(ctx, canvasID, w.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to get pdf for widget %s: %v\n", w.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to get pdf %s: %w", w.ID, err)
			}
			data, err := s.DownloadPDF(ctx, canvasID, pdf.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to download pdf asset %s: %v\n", pdf.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to download pdf asset %s: %w", pdf.ID, err)
			}
			filename := "pdf_" + w.ID + ".pdf"
			filePath := filepath.Join(exportFolder, filename)
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				fmt.Printf("[EXPORT] Failed to write pdf file %s: %v\n", filePath, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to write pdf file: %w", err)
			}
			fmt.Printf("[EXPORT] Wrote pdf file: %s\n", filePath)
			assets[w.ID] = filename
		case "video":
			fmt.Printf("[EXPORT] Attempting to export video asset for widget %s\n", w.ID)
			video, err := s.GetVideo(ctx, canvasID, w.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to get video for widget %s: %v\n", w.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to get video %s: %w", w.ID, err)
			}
			data, err := s.DownloadVideo(ctx, canvasID, video.ID)
			if err != nil {
				fmt.Printf("[EXPORT] Failed to download video asset %s: %v\n", video.ID, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to download video asset %s: %w", video.ID, err)
			}
			filename := "video_" + w.ID + ".mp4"
			filePath := filepath.Join(exportFolder, filename)
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				fmt.Printf("[EXPORT] Failed to write video file %s: %v\n", filePath, err)
				return "", fmt.Errorf("ExportWidgetsToFolder: failed to write video file: %w", err)
			}
			fmt.Printf("[EXPORT] Wrote video file: %s\n", filePath)
			assets[w.ID] = filename
		}
	}
	exportJSON := struct {
		Widgets []Widget          `json:"widgets"`
		Assets  map[string]string `json:"assets"`
		Region  *Rectangle        `json:"region"`
	}{
		Widgets: selected,
		Assets:  assets,
		Region:  &region,
	}
	jsonPath := filepath.Join(exportFolder, "export.json")
	jsonBytes, err := json.MarshalIndent(exportJSON, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ExportWidgetsToFolder: failed to marshal export JSON: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonBytes, 0644); err != nil {
		return "", fmt.Errorf("ExportWidgetsToFolder: failed to write export JSON: %w", err)
	}
	return exportFolder, nil
}

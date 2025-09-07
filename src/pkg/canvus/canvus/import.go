package canvus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

// ImportWidgetsToRegion imports widgets and assets from an ExportedWidgetSet into a specified region of a canvas.
// Widgets are scaled and translated to fit the target region. Assets (images, pdfs, videos) are uploaded as needed.
// Returns a slice of new widget IDs and any errors encountered.
func (s *Session) ImportWidgetsToRegion(ctx context.Context, canvasID string, exported *ExportedWidgetSet, targetRegion Rectangle) ([]string, error) {
	if exported == nil || len(exported.Widgets) == 0 {
		return nil, nil
	}
	var newIDs []string
	orig := exported.Region
	if orig == nil {
		return nil, fmt.Errorf("ImportWidgetsToRegion: original region is nil")
	}
	scaleX := targetRegion.Width / orig.Width
	scaleY := targetRegion.Height / orig.Height
	dx := targetRegion.X - orig.X*scaleX
	dy := targetRegion.Y - orig.Y*scaleY
	for _, w := range exported.Widgets {
		// Scale and translate location and size
		w.Location.X = w.Location.X*scaleX + dx
		w.Location.Y = w.Location.Y*scaleY + dy
		w.Size.Width *= scaleX
		w.Size.Height *= scaleY
		var createdID string
		widgetTypeLower := strings.ToLower(w.WidgetType)
		switch widgetTypeLower {
		case "image":
			imgData, ok := exported.Assets[w.ID]
			if !ok {
				return nil, fmt.Errorf("ImportWidgetsToRegion: missing image data for widget %s", w.ID)
			}
			meta := map[string]interface{}{
				"title":       w.ID,
				"widget_type": widgetTypeLower,
				"location":    map[string]interface{}{"x": w.Location.X, "y": w.Location.Y},
				"size":        map[string]interface{}{"width": w.Size.Width, "height": w.Size.Height},
			}
			body, contentType, buildErr := buildMultipartBody(meta, "data", "imported_image.jpg", []byte(imgData))
			if buildErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to build multipart body for image: %w", buildErr)
			}
			img, createErr := s.CreateImage(ctx, canvasID, body, contentType)
			if createErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to create image: %w", createErr)
			}
			createdID = img.ID
		case "pdf":
			pdfData, ok := exported.Assets[w.ID]
			if !ok {
				return nil, fmt.Errorf("ImportWidgetsToRegion: missing PDF data for widget %s", w.ID)
			}
			meta := map[string]interface{}{
				"title":       w.ID,
				"widget_type": widgetTypeLower,
				"location":    map[string]interface{}{"x": w.Location.X, "y": w.Location.Y},
				"size":        map[string]interface{}{"width": w.Size.Width, "height": w.Size.Height},
			}
			body, contentType, buildErr := buildMultipartBody(meta, "data", "imported.pdf", []byte(pdfData))
			if buildErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to build multipart body for pdf: %w", buildErr)
			}
			pdf, createErr := s.CreatePDF(ctx, canvasID, body, contentType)
			if createErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to create pdf: %w", createErr)
			}
			createdID = pdf.ID
		case "video":
			videoData, ok := exported.Assets[w.ID]
			if !ok {
				return nil, fmt.Errorf("ImportWidgetsToRegion: missing video data for widget %s", w.ID)
			}
			meta := map[string]interface{}{
				"title":       w.ID,
				"widget_type": widgetTypeLower,
				"location":    map[string]interface{}{"x": w.Location.X, "y": w.Location.Y},
				"size":        map[string]interface{}{"width": w.Size.Width, "height": w.Size.Height},
			}
			body, contentType, buildErr := buildMultipartBody(meta, "data", "imported_video.mp4", []byte(videoData))
			if buildErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to build multipart body for video: %w", buildErr)
			}
			video, createErr := s.CreateVideo(ctx, canvasID, body, contentType)
			if createErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to create video: %w", createErr)
			}
			createdID = video.ID
		default:
			created, createErr := s.CreateWidget(ctx, canvasID, widgetToMap(w))
			if createErr != nil {
				return nil, fmt.Errorf("ImportWidgetsToRegion: failed to create widget: %w", createErr)
			}
			createdID = created.ID
		}
		newIDs = append(newIDs, createdID)
	}
	return newIDs, nil
}

// widgetToMap converts a Widget struct to a map[string]interface{} for CreateWidget.
func widgetToMap(w Widget) map[string]interface{} {
	m := map[string]interface{}{
		"widget_type": strings.ToLower(w.WidgetType),
		"parent_id":   w.ParentID,
		"pinned":      w.Pinned,
		"scale":       w.Scale,
		"state":       w.State,
		"depth":       w.Depth,
	}
	if w.Location != nil {
		m["location"] = map[string]interface{}{"x": w.Location.X, "y": w.Location.Y}
	}
	if w.Size != nil {
		m["size"] = map[string]interface{}{"width": w.Size.Width, "height": w.Size.Height}
	}
	return m
}

// buildMultipartBody creates a multipart body with a JSON part and a file part.
func buildMultipartBody(meta map[string]interface{}, fieldName, filename string, fileData []byte) (io.Reader, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	// Add JSON part
	if meta != nil {
		jsonBytes, err := json.Marshal(meta)
		if err != nil {
			return nil, "", err
		}
		jsonPart, err := w.CreateFormField("json")
		if err != nil {
			return nil, "", err
		}
		_, err = jsonPart.Write(jsonBytes)
		if err != nil {
			return nil, "", err
		}
	}
	// Add file part
	if fileData != nil && fieldName != "" && filename != "" {
		filePart, err := w.CreateFormFile(fieldName, filename)
		if err != nil {
			return nil, "", err
		}
		_, err = filePart.Write(fileData)
		if err != nil {
			return nil, "", err
		}
	}
	w.Close()
	return &buf, w.FormDataContentType(), nil
}

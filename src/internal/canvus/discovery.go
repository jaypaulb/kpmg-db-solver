package canvus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	canvussdk "canvus-go-api/canvus"
)

// AssetInfo represents information about a discovered asset
type AssetInfo struct {
	Hash             string `json:"hash"`
	WidgetType       string `json:"widget_type"`
	OriginalFilename string `json:"original_filename"`
	CanvasID         string `json:"canvas_id"`
	CanvasName       string `json:"canvas_name"`
	WidgetID         string `json:"widget_id"`
	WidgetName       string `json:"widget_name"`
}

// DiscoveryResult represents the result of asset discovery
type DiscoveryResult struct {
	Assets    []AssetInfo `json:"assets"`
	Canvases  []canvussdk.Canvas `json:"canvases"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Errors    []string    `json:"errors"`
}

// RateLimiter controls the rate of API requests
type RateLimiter struct {
	requests chan struct{}
	rate     time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rate := time.Second / time.Duration(requestsPerSecond)
	rl := &RateLimiter{
		requests: make(chan struct{}, requestsPerSecond),
		rate:     rate,
	}

	// Start the rate limiter goroutine
	go rl.run()

	return rl
}

// run manages the rate limiting
func (rl *RateLimiter) run() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.requests <- struct{}{}:
		default:
			// Channel is full, skip this tick
		}
	}
}

// Wait blocks until a request slot is available
func (rl *RateLimiter) Wait() {
	<-rl.requests
}

// DiscoverAllAssets discovers all media assets across all canvases using the existing SDK
func DiscoverAllAssets(session *canvussdk.Session, requestsPerSecond int) (*DiscoveryResult, error) {
	startTime := time.Now()
	result := &DiscoveryResult{
		StartTime: startTime,
		Assets:    make([]AssetInfo, 0),
		Canvases:  make([]canvussdk.Canvas, 0),
		Errors:    make([]string, 0),
	}

	ctx := context.Background()

	// Get all canvases using the existing SDK
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get canvases: %w", err)
	}

	result.Canvases = canvases

	// Create rate limiter
	rateLimiter := NewRateLimiter(requestsPerSecond)

	// Process canvases in parallel with rate limiting
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10) // Limit concurrent requests

	for _, canvas := range canvases {
		wg.Add(1)
		go func(canvas canvussdk.Canvas) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			rateLimiter.Wait() // Rate limit

			// Extract media assets from widgets
			assets := extractMediaAssets(ctx, session, canvas)

			mu.Lock()
			result.Assets = append(result.Assets, assets...)
			mu.Unlock()
		}(canvas)
	}

	wg.Wait()

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// extractMediaAssets extracts media assets from widgets by calling the generic ListWidgets endpoint
func extractMediaAssets(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas) []AssetInfo {
	var assets []AssetInfo

	// Get all widgets for this canvas
	widgets, err := session.ListWidgets(ctx, fmt.Sprintf("%d", canvas.ID), nil)
	if err != nil {
		return assets // Return empty slice if we can't get widgets
	}

	// Process each widget and extract media assets
	for _, widget := range widgets {
		// Only process media widget types
		if !isMediaWidget(widget.WidgetType) {
			continue
		}

		// Get the specific widget details based on type
		asset := extractAssetFromWidget(ctx, session, canvas, widget)
		if asset != nil {
			assets = append(assets, *asset)
		}
	}

	return assets
}

// extractAssetFromWidget extracts asset information from a specific widget
func extractAssetFromWidget(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas, widget canvussdk.Widget) *AssetInfo {
	switch widget.WidgetType {
	case "Image":
		return extractImageAsset(ctx, session, canvas, widget)
	case "Pdf":
		return extractPDFAsset(ctx, session, canvas, widget)
	case "Video":
		return extractVideoAsset(ctx, session, canvas, widget)
	default:
		return nil
	}
}

// extractImageAsset extracts asset information from an Image widget
func extractImageAsset(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas, widget canvussdk.Widget) *AssetInfo {
	// Get the specific image details
	image, err := session.GetImage(ctx, fmt.Sprintf("%d", canvas.ID), widget.ID)
	if err != nil {
		return nil
	}

	if image.Hash == "" {
		return nil
	}

	return &AssetInfo{
		Hash:             image.Hash,
		WidgetType:       "Image",
		OriginalFilename: image.OriginalFilename,
		CanvasID:         fmt.Sprintf("%d", canvas.ID),
		CanvasName:       canvas.Name,
		WidgetID:         image.ID,
		WidgetName:       image.Title,
	}
}

// extractPDFAsset extracts asset information from a PDF widget
func extractPDFAsset(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas, widget canvussdk.Widget) *AssetInfo {
	// Get the specific PDF details
	pdf, err := session.GetPDF(ctx, fmt.Sprintf("%d", canvas.ID), widget.ID)
	if err != nil {
		return nil
	}

	// Since PDF struct might be incomplete, we'll use a generic approach
	// We'll need to check if the PDF has a hash field
	// For now, let's assume it follows the same pattern as Image
	// We can use reflection or type assertion to get the hash if it exists
	
	// Try to get hash from the PDF struct - this might need adjustment based on actual struct
	hash := ""
	filename := ""
	name := ""
	
	// Use reflection to get hash field if it exists
	if pdfValue := reflect.ValueOf(pdf); pdfValue.IsValid() && !pdfValue.IsNil() {
		if hashField := pdfValue.Elem().FieldByName("Hash"); hashField.IsValid() && hashField.CanInterface() {
			if hashStr, ok := hashField.Interface().(string); ok && hashStr != "" {
				hash = hashStr
			}
		}
		if filenameField := pdfValue.Elem().FieldByName("OriginalFilename"); filenameField.IsValid() && filenameField.CanInterface() {
			if filenameStr, ok := filenameField.Interface().(string); ok {
				filename = filenameStr
			}
		}
		if nameField := pdfValue.Elem().FieldByName("Name"); nameField.IsValid() && nameField.CanInterface() {
			if nameStr, ok := nameField.Interface().(string); ok {
				name = nameStr
			}
		}
	}

	if hash == "" {
		return nil
	}

	return &AssetInfo{
		Hash:             hash,
		WidgetType:       "Pdf",
		OriginalFilename: filename,
		CanvasID:         fmt.Sprintf("%d", canvas.ID),
		CanvasName:       canvas.Name,
		WidgetID:         widget.ID,
		WidgetName:       name,
	}
}

// extractVideoAsset extracts asset information from a Video widget
func extractVideoAsset(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas, widget canvussdk.Widget) *AssetInfo {
	// Get the specific video details
	video, err := session.GetVideo(ctx, fmt.Sprintf("%d", canvas.ID), widget.ID)
	if err != nil {
		return nil
	}

	// Since Video struct might be incomplete, we'll use a generic approach
	// We'll need to check if the Video has a hash field
	// For now, let's assume it follows the same pattern as Image
	// We can use reflection or type assertion to get the hash if it exists
	
	// Try to get hash from the Video struct - this might need adjustment based on actual struct
	hash := ""
	filename := ""
	name := ""
	
	// Use reflection to get hash field if it exists
	if videoValue := reflect.ValueOf(video); videoValue.IsValid() && !videoValue.IsNil() {
		if hashField := videoValue.Elem().FieldByName("Hash"); hashField.IsValid() && hashField.CanInterface() {
			if hashStr, ok := hashField.Interface().(string); ok && hashStr != "" {
				hash = hashStr
			}
		}
		if filenameField := videoValue.Elem().FieldByName("OriginalFilename"); filenameField.IsValid() && filenameField.CanInterface() {
			if filenameStr, ok := filenameField.Interface().(string); ok {
				filename = filenameStr
			}
		}
		if nameField := videoValue.Elem().FieldByName("Name"); nameField.IsValid() && nameField.CanInterface() {
			if nameStr, ok := nameField.Interface().(string); ok {
				name = nameStr
			}
		}
	}

	if hash == "" {
		return nil
	}

	return &AssetInfo{
		Hash:             hash,
		WidgetType:       "Video",
		OriginalFilename: filename,
		CanvasID:         fmt.Sprintf("%d", canvas.ID),
		CanvasName:       canvas.Name,
		WidgetID:         widget.ID,
		WidgetName:       name,
	}
}

// isMediaWidget checks if a widget type is a media asset
func isMediaWidget(widgetType string) bool {
	switch widgetType {
	case "Pdf", "Image", "Video":
		return true
	default:
		return false
	}
}

// GetUniqueAssets returns unique assets (deduplicated by hash)
func (result *DiscoveryResult) GetUniqueAssets() []AssetInfo {
	hashMap := make(map[string]AssetInfo)

	for _, asset := range result.Assets {
		if _, exists := hashMap[asset.Hash]; !exists {
			hashMap[asset.Hash] = asset
		}
	}

	uniqueAssets := make([]AssetInfo, 0, len(hashMap))
	for _, asset := range hashMap {
		uniqueAssets = append(uniqueAssets, asset)
	}

	return uniqueAssets
}

// GetAssetsByCanvas groups assets by canvas
func (result *DiscoveryResult) GetAssetsByCanvas() map[string][]AssetInfo {
	canvasMap := make(map[string][]AssetInfo)

	for _, asset := range result.Assets {
		canvasMap[asset.CanvasName] = append(canvasMap[asset.CanvasName], asset)
	}

	return canvasMap
}

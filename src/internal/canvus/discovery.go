package canvus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jaypaulb/kpmg-db-solver/internal/logging"
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
	logger := logging.GetLogger()

	// Get all widgets for this canvas
	logger.Verbose("Getting widgets for canvas '%s' (ID: %s)", canvas.Name, canvas.ID)
	widgets, err := session.ListWidgets(ctx, canvas.ID, nil)
	if err != nil {
		logger.Error("Failed to get widgets for canvas '%s' (ID: %s): %v", canvas.Name, canvas.ID, err)
		return assets // Return empty slice if we can't get widgets
	}

	logger.Verbose("Found %d widgets in canvas '%s' (ID: %s)", len(widgets), canvas.Name, canvas.ID)

	// Log the raw widget response in verbose mode
	if len(widgets) > 0 {
		logger.Verbose("Widget response for canvas '%s':", canvas.Name)
		for i, widget := range widgets {
			logger.Verbose("  Widget %d: ID=%s, Type=%s", i+1, widget.ID, widget.WidgetType)
		}
	}

	// Process each widget and extract media assets
	mediaCount := 0
	for _, widget := range widgets {
		// Get the specific widget details to check for hash field
		asset := extractAssetFromWidget(ctx, session, canvas, widget)
		if asset != nil {
			assets = append(assets, *asset)
			mediaCount++
			logger.Verbose("Found media asset: %s (%s) - Hash: %s", asset.WidgetName, asset.WidgetType, asset.Hash)
		}
	}

	logger.Verbose("Extracted %d media assets from canvas '%s' (ID: %s)", mediaCount, canvas.Name, canvas.ID)
	return assets
}

// extractAssetFromWidget extracts asset information from a widget if it has a hash field
func extractAssetFromWidget(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas, widget canvussdk.Widget) *AssetInfo {
	logger := logging.GetLogger()

	// Get the specific widget details based on type
	var widgetDetails interface{}
	var err error

	logger.Verbose("Getting details for widget ID=%s, Type=%s in canvas '%s'", widget.ID, widget.WidgetType, canvas.Name)

	switch widget.WidgetType {
	case "Image":
		widgetDetails, err = session.GetImage(ctx, canvas.ID, widget.ID)
	case "Pdf":
		widgetDetails, err = session.GetPDF(ctx, canvas.ID, widget.ID)
	case "Video":
		widgetDetails, err = session.GetVideo(ctx, canvas.ID, widget.ID)
	default:
		logger.Verbose("Skipping non-media widget type: %s", widget.WidgetType)
		return nil // Not a media widget type
	}

	if err != nil {
		logger.Verbose("Failed to get widget details for ID=%s, Type=%s: %v", widget.ID, widget.WidgetType, err)
		return nil
	}

	// Extract hash and other fields using reflection
	hash := ""
	filename := ""
	name := ""

	if widgetValue := reflect.ValueOf(widgetDetails); widgetValue.IsValid() && !widgetValue.IsNil() {
		// Get hash field - if it exists and is not empty, this is a media asset
		if hashField := widgetValue.Elem().FieldByName("Hash"); hashField.IsValid() && hashField.CanInterface() {
			if hashStr, ok := hashField.Interface().(string); ok && hashStr != "" {
				hash = hashStr
				logger.Verbose("Found hash for widget ID=%s: %s", widget.ID, hash)
			}
		}

		// Get filename field
		if filenameField := widgetValue.Elem().FieldByName("OriginalFilename"); filenameField.IsValid() && filenameField.CanInterface() {
			if filenameStr, ok := filenameField.Interface().(string); ok {
				filename = filenameStr
				logger.Verbose("Found filename for widget ID=%s: %s", widget.ID, filename)
			}
		}

		// Get name field (could be Title, Name, etc.)
		if nameField := widgetValue.Elem().FieldByName("Title"); nameField.IsValid() && nameField.CanInterface() {
			if nameStr, ok := nameField.Interface().(string); ok {
				name = nameStr
				logger.Verbose("Found title for widget ID=%s: %s", widget.ID, name)
			}
		} else if nameField := widgetValue.Elem().FieldByName("Name"); nameField.IsValid() && nameField.CanInterface() {
			if nameStr, ok := nameField.Interface().(string); ok {
				name = nameStr
				logger.Verbose("Found name for widget ID=%s: %s", widget.ID, name)
			}
		}
	}

	// Only return asset if it has a hash (media assets only)
	if hash == "" {
		logger.Verbose("No hash found for widget ID=%s, Type=%s - not a media asset", widget.ID, widget.WidgetType)
		return nil
	}

	return &AssetInfo{
		Hash:             hash,
		WidgetType:       widget.WidgetType,
		OriginalFilename: filename,
		CanvasID:         canvas.ID,
		CanvasName:       canvas.Name,
		WidgetID:         widget.ID,
		WidgetName:       name,
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

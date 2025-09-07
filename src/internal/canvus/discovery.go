package canvus

import (
	"context"
	"fmt"
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

// extractMediaAssets extracts media assets from widgets by calling specific API endpoints
func extractMediaAssets(ctx context.Context, session *canvussdk.Session, canvas canvussdk.Canvas) []AssetInfo {
	var assets []AssetInfo

	// Get images
	images, err := session.ListImages(ctx, fmt.Sprintf("%d", canvas.ID))
	if err == nil {
		for _, image := range images {
			if image.Hash != "" {
				asset := AssetInfo{
					Hash:             image.Hash,
					WidgetType:       "Image",
					OriginalFilename: image.OriginalFilename,
					CanvasID:         fmt.Sprintf("%d", canvas.ID),
					CanvasName:       canvas.Name,
					WidgetID:         image.ID,
					WidgetName:       image.Title,
				}
				assets = append(assets, asset)
			}
		}
	}

	// Get PDFs (note: PDF struct may be incomplete, but we'll try)
	pdfs, err := session.ListPDFs(ctx, fmt.Sprintf("%d", canvas.ID))
	if err == nil {
		for _, pdf := range pdfs {
			// Since PDF struct is incomplete, we'll need to handle this differently
			// For now, we'll skip PDFs until we can determine the correct structure
			_ = pdf
		}
	}

	// Get Videos (note: Video struct may be incomplete, but we'll try)
	videos, err := session.ListVideos(ctx, fmt.Sprintf("%d", canvas.ID))
	if err == nil {
		for _, video := range videos {
			// Since Video struct is incomplete, we'll need to handle this differently
			// For now, we'll skip Videos until we can determine the correct structure
			_ = video
		}
	}

	return assets
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

package canvus

import (
	"context"
	"fmt"
)

// GetMipmapInfo retrieves mipmap information for a given asset hash and page.
// Requires 'canvas-id' and 'Private-Token' headers.
func (s *Session) GetMipmapInfo(ctx context.Context, canvasID, publicHashHex string, page *int) (*MipmapInfo, error) {
	var info MipmapInfo
	path := fmt.Sprintf("mipmaps/%s", publicHashHex)
	params := map[string]interface{}{}
	if page != nil {
		params["page"] = *page
	}
	err := s.doRequestWithHeaders(ctx, "GET", path, nil, &info, params, map[string]string{"canvas-id": canvasID}, false)
	if err != nil {
		return nil, fmt.Errorf("GetMipmapInfo: %w", err)
	}
	return &info, nil
}

// GetMipmapLevel retrieves a specific mipmap level image (WebP format).
// Returns binary data. Requires 'canvas-id' and 'Private-Token' headers.
func (s *Session) GetMipmapLevel(ctx context.Context, canvasID, publicHashHex string, level int, page *int) ([]byte, error) {
	path := fmt.Sprintf("mipmaps/%s/%d", publicHashHex, level)
	params := map[string]interface{}{}
	if page != nil {
		params["page"] = *page
	}
	var data []byte
	err := s.doRequestWithHeaders(ctx, "GET", path, nil, &data, params, map[string]string{"canvas-id": canvasID}, true)
	if err != nil {
		return nil, fmt.Errorf("GetMipmapLevel: %w", err)
	}
	return data, nil
}

// GetAssetByHash retrieves an asset file by its hash.
// Returns binary data. Requires 'canvas-id' and 'Private-Token' headers.
func (s *Session) GetAssetByHash(ctx context.Context, canvasID, publicHashHex string) ([]byte, error) {
	path := fmt.Sprintf("assets/%s", publicHashHex)
	var data []byte
	err := s.doRequestWithHeaders(ctx, "GET", path, nil, &data, nil, map[string]string{"canvas-id": canvasID}, true)
	if err != nil {
		return nil, fmt.Errorf("GetAssetByHash: %w", err)
	}
	return data, nil
}

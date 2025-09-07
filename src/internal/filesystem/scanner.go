package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileInfo represents information about a file in the assets folder
type FileInfo struct {
	Path         string `json:"path"`
	Hash         string `json:"hash"`
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	RelativePath string `json:"relative_path"` // Relative path from assets root (preserves folder structure)
}

// ScanResult represents the result of filesystem scanning
type ScanResult struct {
	Files    []FileInfo `json:"files"`
	HashMap  map[string]FileInfo `json:"hash_map"`
	TotalSize int64     `json:"total_size"`
	Error    error      `json:"error,omitempty"`
}

// ScanAssetsFolder scans the assets folder and builds a hash map
func ScanAssetsFolder(assetsPath string) (*ScanResult, error) {
	result := &ScanResult{
		Files:    make([]FileInfo, 0),
		HashMap:  make(map[string]FileInfo),
		TotalSize: 0,
	}

	// Check if assets path exists
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("assets folder does not exist: %s", assetsPath)
	}

	// Walk the directory tree
	err := filepath.Walk(assetsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Extract hash from filename
		hash := extractHashFromFilename(info.Name())
		if hash == "" {
			return nil // Skip files that don't match hash pattern
		}

		// Calculate relative path from assets root to preserve folder structure
		relPath, err := filepath.Rel(assetsPath, path)
		if err != nil {
			// If we can't calculate relative path, use just the filename
			relPath = info.Name()
		}

		fileInfo := FileInfo{
			Path:         path,
			Hash:         hash,
			Filename:     info.Name(),
			Size:         info.Size(),
			RelativePath: relPath,
		}

		result.Files = append(result.Files, fileInfo)
		result.HashMap[hash] = fileInfo
		result.TotalSize += info.Size()

		return nil
	})

	if err != nil {
		result.Error = err
		return result, err
	}

	return result, nil
}

// extractHashFromFilename extracts the hash from a filename
// Expected format: {hash}.{extension}
func extractHashFromFilename(filename string) string {
	// Remove extension
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}

	hash := strings.TrimSuffix(filename, ext)

	// Basic validation: hash should be alphanumeric and reasonable length
	if len(hash) < 8 || len(hash) > 64 {
		return ""
	}

	// Check if hash contains only valid characters (alphanumeric)
	for _, char := range hash {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return ""
		}
	}

	return hash
}

// FindMissingAssets compares discovered assets with filesystem contents
func FindMissingAssets(discoveredAssets []string, scanResult *ScanResult) []string {
	var missingAssets []string

	for _, hash := range discoveredAssets {
		if _, exists := scanResult.HashMap[hash]; !exists {
			missingAssets = append(missingAssets, hash)
		}
	}

	return missingAssets
}

// ParallelScanAssetsFolder scans the assets folder using parallel processing
func ParallelScanAssetsFolder(assetsPath string, numWorkers int) (*ScanResult, error) {
	result := &ScanResult{
		Files:    make([]FileInfo, 0),
		HashMap:  make(map[string]FileInfo),
		TotalSize: 0,
	}

	// Check if assets path exists
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("assets folder does not exist: %s", assetsPath)
	}

	// Channel to collect file paths
	filePaths := make(chan string, 1000)

	// Channel to collect results
	results := make(chan FileInfo, 1000)

	// Channel to signal completion
	done := make(chan bool)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filePaths {
				info, err := os.Stat(path)
				if err != nil {
					continue
				}

				// Skip directories
				if info.IsDir() {
					continue
				}

				// Extract hash from filename
				hash := extractHashFromFilename(info.Name())
				if hash == "" {
					continue // Skip files that don't match hash pattern
				}

				// Calculate relative path from assets root to preserve folder structure
				relPath, err := filepath.Rel(assetsPath, path)
				if err != nil {
					// If we can't calculate relative path, use just the filename
					relPath = info.Name()
				}

				fileInfo := FileInfo{
					Path:         path,
					Hash:         hash,
					Filename:     info.Name(),
					Size:         info.Size(),
					RelativePath: relPath,
				}

				results <- fileInfo
			}
		}()
	}

	// Start result collector
	go func() {
		defer close(done)
		for fileInfo := range results {
			result.Files = append(result.Files, fileInfo)
			result.HashMap[fileInfo.Hash] = fileInfo
			result.TotalSize += fileInfo.Size
		}
	}()

	// Walk the directory tree and send paths to workers
	err := filepath.Walk(assetsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Send file paths to workers
		if !info.IsDir() {
			filePaths <- path
		}

		return nil
	})

	close(filePaths)
	wg.Wait()
	close(results)
	<-done

	if err != nil {
		result.Error = err
		return result, err
	}

	return result, nil
}

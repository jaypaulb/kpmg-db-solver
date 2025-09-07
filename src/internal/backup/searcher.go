package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jaypaulb/kpmg-db-solver/internal/logging"
)

// BackupFile represents a found backup file
type BackupFile struct {
	Path         string    // Full path to the backup file
	Hash         string    // Asset hash (filename without extension)
	Extension    string    // File extension
	ModifiedTime time.Time // File modification time
	Size         int64     // File size in bytes
	RelativePath string    // Relative path from backup root (preserves folder structure)
}

// SearchResult contains the results of a backup search
type SearchResult struct {
	FoundFiles    map[string][]BackupFile // Hash -> list of backup files (newest first)
	MissingHashes []string                // Hashes that were not found in any backup
	TotalSearched int                     // Total number of backup directories searched
	TotalFiles    int                     // Total number of files found
}

// Searcher handles searching for backup files
type Searcher struct {
	backupRootFolder string
	logger           *logging.Logger
}

// NewSearcher creates a new backup searcher
func NewSearcher(backupRootFolder string) *Searcher {
	return &Searcher{
		backupRootFolder: backupRootFolder,
		logger:           logging.GetLogger(),
	}
}

// SearchForAssets searches for missing assets in backup folders
// Returns a map of hash -> list of backup files (newest first)
// Looks for backup folders with pattern: {timestamp}_{date}_{version}_mt-canvus_backup\assets\
func (s *Searcher) SearchForAssets(missingHashes []string) (*SearchResult, error) {
	result := &SearchResult{
		FoundFiles:    make(map[string][]BackupFile),
		MissingHashes: make([]string, 0),
		TotalSearched: 0,
		TotalFiles:    0,
	}

	if len(missingHashes) == 0 {
		s.logger.Info("No missing assets to search for")
		return result, nil
	}

	s.logger.Info("🔍 Searching for %d missing assets in backup folder: %s", len(missingHashes), s.backupRootFolder)

	// Create a set of missing hashes for efficient lookup
	missingSet := make(map[string]bool)
	for _, hash := range missingHashes {
		missingSet[hash] = true
	}

	// Check if backup root folder exists
	if !s.pathExists(s.backupRootFolder) {
		s.logger.Warn("Backup folder does not exist: %s", s.backupRootFolder)
		return result, nil
	}

	// Find backup folders with the expected pattern and search their assets subfolder
	err := s.searchBackupFolders(s.backupRootFolder, missingSet, result)
	if err != nil {
		s.logger.Error("Error searching backup folders %s: %v", s.backupRootFolder, err)
		return result, err
	}

	result.TotalSearched = 1 // We searched one root folder

	// Identify hashes that were not found
	for hash := range missingSet {
		if _, found := result.FoundFiles[hash]; !found {
			result.MissingHashes = append(result.MissingHashes, hash)
		}
	}

	s.logger.Info("✅ Backup search completed:")
	s.logger.Info("   📁 Folders searched: %d", result.TotalSearched)
	s.logger.Info("   📄 Files found: %d", result.TotalFiles)
	s.logger.Info("   ✅ Assets found: %d", len(result.FoundFiles))
	s.logger.Info("   ❌ Assets still missing: %d", len(result.MissingHashes))

	return result, nil
}

// searchBackupFolders finds backup folders with the expected pattern and searches their assets subfolder
func (s *Searcher) searchBackupFolders(backupRoot string, missingSet map[string]bool, result *SearchResult) error {
	// Look for backup folders with pattern: {timestamp}_{date}_{version}_mt-canvus_backup
	entries, err := os.ReadDir(backupRoot)
	if err != nil {
		return fmt.Errorf("failed to read backup root directory: %w", err)
	}

	backupFolders := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && s.isBackupFolder(entry.Name()) {
			backupPath := filepath.Join(backupRoot, entry.Name())
			assetsPath := filepath.Join(backupPath, "assets")
			
			// Check if this backup folder has an assets subfolder
			if s.pathExists(assetsPath) {
				backupFolders = append(backupFolders, assetsPath)
				s.logger.Verbose("Found backup assets folder: %s", assetsPath)
			}
		}
	}

	if len(backupFolders) == 0 {
		s.logger.Warn("No backup folders with assets subfolder found in: %s", backupRoot)
		return nil
	}

	// Search each backup assets folder
	for _, assetsPath := range backupFolders {
		err := s.searchBackupFolder(assetsPath, missingSet, result)
		if err != nil {
			s.logger.Error("Error searching backup folder %s: %v", assetsPath, err)
			continue
		}
	}

	return nil
}

// isBackupFolder checks if a folder name matches the backup pattern
func (s *Searcher) isBackupFolder(folderName string) bool {
	// Pattern: {timestamp}_{date}_{version}_mt-canvus_backup
	// Example: 1757261054_2025_09_07_3.3.0_mt-canvus_backup
	return strings.Contains(folderName, "_mt-canvus_backup")
}

// searchBackupFolder recursively searches a backup folder for missing assets
func (s *Searcher) searchBackupFolder(backupRoot string, missingSet map[string]bool, result *SearchResult) error {
	return filepath.Walk(backupRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log error but continue searching
			s.logger.Verbose("Error accessing %s: %v", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Extract hash from filename (filename without extension)
		filename := info.Name()
		ext := filepath.Ext(filename)
		hash := strings.TrimSuffix(filename, ext)

		// Check if this hash is one we're looking for
		if missingSet[hash] {
			// Calculate relative path from the assets subfolder to preserve folder structure
			// This ensures the path starts from assets\ and aligns with the target assets\ folder
			relPath, err := filepath.Rel(backupRoot, path)
			if err != nil {
				// If we can't calculate relative path, use just the filename
				relPath = info.Name()
			}
			
			backupFile := BackupFile{
				Path:         path,
				Hash:         hash,
				Extension:    ext,
				ModifiedTime: info.ModTime(),
				Size:         info.Size(),
				RelativePath: relPath,
			}

			// Add to results (will be sorted by modification time later)
			result.FoundFiles[hash] = append(result.FoundFiles[hash], backupFile)
			result.TotalFiles++

			s.logger.Verbose("Found backup: %s (hash: %s, size: %d bytes, modified: %s)",
				path, hash, info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		}

		return nil
	})
}

// SortBackupFiles sorts backup files by modification time (newest first)
func (s *Searcher) SortBackupFiles(result *SearchResult) {
	for hash, files := range result.FoundFiles {
		// Sort by modification time (newest first)
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModifiedTime.After(files[j].ModifiedTime)
		})
		result.FoundFiles[hash] = files
	}
}

// GetBestBackupFile returns the newest backup file for a given hash
func (s *Searcher) GetBestBackupFile(result *SearchResult, hash string) *BackupFile {
	files, exists := result.FoundFiles[hash]
	if !exists || len(files) == 0 {
		return nil
	}

	// Files are already sorted by modification time (newest first)
	return &files[0]
}

// pathExists checks if a path exists
func (s *Searcher) pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

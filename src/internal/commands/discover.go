package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jaypaulb/kpmg-db-solver/internal/backup"
	"github.com/jaypaulb/kpmg-db-solver/internal/canvus"
	"github.com/jaypaulb/kpmg-db-solver/internal/config"
	"github.com/jaypaulb/kpmg-db-solver/internal/filesystem"
	"github.com/jaypaulb/kpmg-db-solver/internal/logging"
	canvussdk "canvus-go-api/canvus"
)

// DiscoverCommand handles the discover command
type DiscoverCommand struct {
	config *config.Config
}

// NewDiscoverCommand creates a new discover command
func NewDiscoverCommand(cfg *config.Config) *DiscoverCommand {
	return &DiscoverCommand{
		config: cfg,
	}
}

// Execute runs the discover command
func (cmd *DiscoverCommand) Execute() error {
	logger := logging.GetLogger()

	logger.Info("🔍 Starting asset discovery...")
	logger.Info("📡 Connecting to Canvus Server: %s", cmd.config.CanvusServer.URL)
	logger.Info("📁 Scanning assets folder: %s", cmd.config.Paths.AssetsFolder)

	// Create Canvus session using existing SDK
	ctx := context.Background()
	var session *canvussdk.Session
	if cmd.config.CanvusServer.InsecureTLS {
		session = canvussdk.NewSession(cmd.config.GetCanvusAPIURL(), canvussdk.WithInsecureTLS())
	} else {
		session = canvussdk.NewSession(cmd.config.GetCanvusAPIURL())
	}

	// Authenticate using existing SDK
	logger.Info("🔐 Authenticating with Canvus Server...")
	err := session.Login(ctx, cmd.config.CanvusServer.Username, cmd.config.CanvusServer.Password)
	if err != nil {
		logger.Error("Authentication failed: %v", err)
		return fmt.Errorf("authentication failed: %w", err)
	}
	defer session.Logout(ctx)

	logger.Info("✅ Successfully authenticated with Canvus Server")

	// Discover assets from API
	logger.Info("📊 Discovering assets from Canvus API...")
	discoveryResult, err := canvus.DiscoverAllAssets(session, cmd.config.Performance.MaxConcurrentAPI)
	if err != nil {
		logger.Error("Asset discovery failed: %v", err)
		return fmt.Errorf("asset discovery failed: %w", err)
	}

	logger.Info("📈 Found %d canvases with %d total media assets",
		len(discoveryResult.Canvases), len(discoveryResult.Assets))

	// Get unique assets
	uniqueAssets := discoveryResult.GetUniqueAssets()
	logger.Info("🔗 Unique assets (deduplicated): %d", len(uniqueAssets))

	// Extract asset hashes for filesystem comparison
	assetHashes := make([]string, len(uniqueAssets))
	for i, asset := range uniqueAssets {
		assetHashes[i] = asset.Hash
	}

	// Scan filesystem
	logger.Info("💾 Scanning assets folder...")
	scanResult, err := filesystem.ScanAssetsFolder(cmd.config.Paths.AssetsFolder)
	if err != nil {
		logger.Error("Filesystem scan failed: %v", err)
		return fmt.Errorf("filesystem scan failed: %w", err)
	}

	logger.Info("📂 Found %d files in assets folder (%.2f MB total)",
		len(scanResult.Files), float64(scanResult.TotalSize)/(1024*1024))

	// Find missing assets
	missingAssets := filesystem.FindMissingAssets(assetHashes, scanResult)
	logger.Info("❌ Missing assets: %d", len(missingAssets))

	// Search for missing assets in backup folders
	var backupSearchResult *backup.SearchResult
	if len(missingAssets) > 0 {
		logger.Info("🔍 Searching for missing assets in backup folder...")
		searcher := backup.NewSearcher(cmd.config.Paths.BackupRootFolder)
		backupSearchResult, err = searcher.SearchForAssets(missingAssets)
		if err != nil {
			logger.Error("Backup search failed: %v", err)
			return fmt.Errorf("backup search failed: %w", err)
		}

		// Sort backup files by modification time (newest first)
		searcher.SortBackupFiles(backupSearchResult)

		// Offer to restore found assets
		if len(backupSearchResult.FoundFiles) > 0 {
			logger.Info("💾 Found %d missing assets in backup folder", len(backupSearchResult.FoundFiles))
			logger.Info("🔄 Would you like to restore these assets? (y/N): ")

			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
				logger.Info("🔄 Restoring assets...")
				restorer := backup.NewRestorer(cmd.config.Paths.AssetsFolder)
				restoreResult, err := restorer.RestoreAssets(backupSearchResult)
				if err != nil {
					logger.Error("Asset restoration failed: %v", err)
					return fmt.Errorf("asset restoration failed: %w", err)
				}

				logger.Info("✅ Restoration completed:")
				logger.Info("   ✅ Files restored: %d", len(restoreResult.RestoredFiles))
				logger.Info("   ❌ Files failed: %d", len(restoreResult.FailedFiles))
				logger.Info("   📊 Total bytes: %d", restoreResult.TotalBytes)

				if len(restoreResult.Errors) > 0 {
					logger.Warn("Restoration errors:")
					for _, errMsg := range restoreResult.Errors {
						logger.Warn("   %s", errMsg)
					}
				}
			} else {
				logger.Info("⏭️ Skipping asset restoration")
			}
		}
	}

	// Generate reports
	if len(missingAssets) > 0 {
		logger.Info("📋 Generating reports...")
		err = cmd.generateReports(discoveryResult, missingAssets, uniqueAssets, backupSearchResult)
		if err != nil {
			logger.Error("Report generation failed: %v", err)
			return fmt.Errorf("report generation failed: %w", err)
		}
	} else {
		logger.Info("✅ No missing assets found!")
	}

	// Print summary
	cmd.printSummary(discoveryResult, scanResult, missingAssets)

	return nil
}

// generateReports generates detailed and CSV reports
func (cmd *DiscoverCommand) generateReports(discoveryResult *canvus.DiscoveryResult, missingAssets []string, uniqueAssets []canvus.AssetInfo, backupSearchResult *backup.SearchResult) error {
	// Create missing assets map for quick lookup
	missingMap := make(map[string]bool)
	for _, hash := range missingAssets {
		missingMap[hash] = true
	}

	// Filter unique assets to only include missing ones
	var missingAssetInfos []canvus.AssetInfo
	for _, asset := range uniqueAssets {
		if missingMap[asset.Hash] {
			missingAssetInfos = append(missingAssetInfos, asset)
		}
	}

	// Generate detailed report
	err := cmd.generateDetailedReport(missingAssetInfos, backupSearchResult)
	if err != nil {
		return fmt.Errorf("failed to generate detailed report: %w", err)
	}

	// Generate CSV report
	err = cmd.generateCSVReport(missingAssetInfos, backupSearchResult)
	if err != nil {
		return fmt.Errorf("failed to generate CSV report: %w", err)
	}

	return nil
}

// generateDetailedReport generates a detailed text report
func (cmd *DiscoverCommand) generateDetailedReport(missingAssets []canvus.AssetInfo, backupSearchResult *backup.SearchResult) error {
	reportPath := "missing_assets_report.txt"

	// Group assets by canvas
	canvasMap := make(map[string][]canvus.AssetInfo)
	for _, asset := range missingAssets {
		canvasMap[asset.CanvasName] = append(canvasMap[asset.CanvasName], asset)
	}

	// Generate report content
	content := fmt.Sprintf("KPMG DB Solver - Missing Assets Report\n")
	content += fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("Total Missing Assets: %d\n\n", len(missingAssets))

	for canvasName, assets := range canvasMap {
		content += fmt.Sprintf("Canvas: %s (ID: %s)\n", canvasName, assets[0].CanvasID)
		for _, asset := range assets {
			content += fmt.Sprintf("  Widget: %s (ID: %s, Type: %s)\n", asset.WidgetName, asset.WidgetID, asset.WidgetType)
			content += fmt.Sprintf("    Hash: %s\n", asset.Hash)
			if asset.OriginalFilename != "" {
				content += fmt.Sprintf("    Original Filename: %s\n", asset.OriginalFilename)
			}

			// Add backup status
			if backupSearchResult != nil {
				if backupFiles, found := backupSearchResult.FoundFiles[asset.Hash]; found && len(backupFiles) > 0 {
					bestBackup := backupFiles[0] // Newest file
					content += fmt.Sprintf("    Backup Status: ✅ Found in backup\n")
					content += fmt.Sprintf("    Backup Path: %s\n", bestBackup.Path)
					content += fmt.Sprintf("    Backup Size: %d bytes\n", bestBackup.Size)
					content += fmt.Sprintf("    Backup Modified: %s\n", bestBackup.ModifiedTime.Format("2006-01-02 15:04:05"))
				} else {
					content += fmt.Sprintf("    Backup Status: ❌ Not found in any backup\n")
				}
			}
			content += "\n"
		}
	}

	// Write report to file
	err := writeFile(reportPath, content)
	if err != nil {
		return err
	}

	fmt.Printf("📄 Detailed report saved to: %s\n", reportPath)
	return nil
}

// generateCSVReport generates a CSV report
func (cmd *DiscoverCommand) generateCSVReport(missingAssets []canvus.AssetInfo, backupSearchResult *backup.SearchResult) error {
	reportPath := "missing_assets.csv"

	// Generate CSV content
	content := "Hash,WidgetType,OriginalFilename,CanvasID,CanvasName,WidgetID,WidgetName,BackupStatus,BackupPath,BackupSize,BackupModified\n"

	for _, asset := range missingAssets {
		backupStatus := "Not Found"
		backupPath := ""
		backupSize := ""
		backupModified := ""

		if backupSearchResult != nil {
			if backupFiles, found := backupSearchResult.FoundFiles[asset.Hash]; found && len(backupFiles) > 0 {
				bestBackup := backupFiles[0] // Newest file
				backupStatus = "Found"
				backupPath = bestBackup.Path
				backupSize = fmt.Sprintf("%d", bestBackup.Size)
				backupModified = bestBackup.ModifiedTime.Format("2006-01-02 15:04:05")
			}
		}

		content += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			asset.Hash,
			asset.WidgetType,
			asset.OriginalFilename,
			asset.CanvasID,
			asset.CanvasName,
			asset.WidgetID,
			asset.WidgetName,
			backupStatus,
			backupPath,
			backupSize,
			backupModified,
		)
	}

	// Write CSV to file
	err := writeFile(reportPath, content)
	if err != nil {
		return err
	}

	fmt.Printf("📊 CSV report saved to: %s\n", reportPath)
	return nil
}

// printSummary prints a summary of the discovery results
func (cmd *DiscoverCommand) printSummary(discoveryResult *canvus.DiscoveryResult, scanResult *filesystem.ScanResult, missingAssets []string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 DISCOVERY SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("⏱️  Discovery Duration: %v\n", discoveryResult.Duration)
	fmt.Printf("📈 Total Canvases: %d\n", len(discoveryResult.Canvases))
	fmt.Printf("🎯 Total Media Assets: %d\n", len(discoveryResult.Assets))
	fmt.Printf("🔗 Unique Assets: %d\n", len(discoveryResult.GetUniqueAssets()))
	fmt.Printf("💾 Files in Assets Folder: %d\n", len(scanResult.Files))
	fmt.Printf("💽 Total Assets Size: %.2f MB\n", float64(scanResult.TotalSize)/(1024*1024))
	fmt.Printf("❌ Missing Assets (Filesystem): %d\n", len(missingAssets))

	// Show server validation results if available
	if discoveryResult.ServerValidation != nil {
		fmt.Printf("🔍 Server Validation: %d/%d assets exist on server\n",
			discoveryResult.ServerValidation.ExistingAssets,
			discoveryResult.ServerValidation.TotalAssets)
		if discoveryResult.ServerValidation.MissingAssets > 0 {
			fmt.Printf("❌ Missing Assets (Server): %d\n", discoveryResult.ServerValidation.MissingAssets)
		}
	}

	if len(discoveryResult.Errors) > 0 {
		fmt.Printf("⚠️  Errors Encountered: %d\n", len(discoveryResult.Errors))
		for _, err := range discoveryResult.Errors {
			fmt.Printf("   - %s\n", err)
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}

// writeFile writes content to a file
func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filename, err)
	}

	return nil
}

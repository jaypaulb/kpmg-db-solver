package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jaypaulb/kpmg-db-solver/internal/canvus"
	"github.com/jaypaulb/kpmg-db-solver/internal/config"
	"github.com/jaypaulb/kpmg-db-solver/internal/filesystem"
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
	fmt.Println("🔍 Starting asset discovery...")
	fmt.Printf("📡 Connecting to Canvus Server: %s\n", cmd.config.CanvusServer.URL)
	fmt.Printf("📁 Scanning assets folder: %s\n", cmd.config.Paths.AssetsFolder)

	// Create Canvus session using existing SDK
	ctx := context.Background()
	session := canvussdk.NewSession(cmd.config.GetCanvusAPIURL())

	// Authenticate using existing SDK
	fmt.Println("🔐 Authenticating with Canvus Server...")
	err := session.Login(ctx, cmd.config.CanvusServer.Username, cmd.config.CanvusServer.Password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	defer session.Logout(ctx)

	fmt.Printf("✅ Authenticated successfully\n")

	// Discover assets from API
	fmt.Println("📊 Discovering assets from Canvus API...")
	discoveryResult, err := canvus.DiscoverAllAssets(session, cmd.config.Performance.MaxConcurrentAPI)
	if err != nil {
		return fmt.Errorf("asset discovery failed: %w", err)
	}

	fmt.Printf("📈 Found %d canvases with %d total media assets\n",
		len(discoveryResult.Canvases), len(discoveryResult.Assets))

	// Get unique assets
	uniqueAssets := discoveryResult.GetUniqueAssets()
	fmt.Printf("🔗 Unique assets (deduplicated): %d\n", len(uniqueAssets))

	// Extract asset hashes for filesystem comparison
	assetHashes := make([]string, len(uniqueAssets))
	for i, asset := range uniqueAssets {
		assetHashes[i] = asset.Hash
	}

	// Scan filesystem
	fmt.Println("💾 Scanning assets folder...")
	scanResult, err := filesystem.ScanAssetsFolder(cmd.config.Paths.AssetsFolder)
	if err != nil {
		return fmt.Errorf("filesystem scan failed: %w", err)
	}

	fmt.Printf("📂 Found %d files in assets folder (%.2f MB total)\n",
		len(scanResult.Files), float64(scanResult.TotalSize)/(1024*1024))

	// Find missing assets
	missingAssets := filesystem.FindMissingAssets(assetHashes, scanResult)
	fmt.Printf("❌ Missing assets: %d\n", len(missingAssets))

	// Generate reports
	if len(missingAssets) > 0 {
		fmt.Println("📋 Generating reports...")
		err = cmd.generateReports(discoveryResult, missingAssets, uniqueAssets)
		if err != nil {
			return fmt.Errorf("report generation failed: %w", err)
		}
	} else {
		fmt.Println("✅ No missing assets found!")
	}

	// Print summary
	cmd.printSummary(discoveryResult, scanResult, missingAssets)

	return nil
}

// generateReports generates detailed and CSV reports
func (cmd *DiscoverCommand) generateReports(discoveryResult *canvus.DiscoveryResult, missingAssets []string, uniqueAssets []canvus.AssetInfo) error {
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
	err := cmd.generateDetailedReport(missingAssetInfos)
	if err != nil {
		return fmt.Errorf("failed to generate detailed report: %w", err)
	}

	// Generate CSV report
	err = cmd.generateCSVReport(missingAssetInfos)
	if err != nil {
		return fmt.Errorf("failed to generate CSV report: %w", err)
	}

	return nil
}

// generateDetailedReport generates a detailed text report
func (cmd *DiscoverCommand) generateDetailedReport(missingAssets []canvus.AssetInfo) error {
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
		content += fmt.Sprintf("Canvas: %s (ID: %d)\n", canvasName, assets[0].CanvasID)
		for _, asset := range assets {
			content += fmt.Sprintf("  Widget: %s (ID: %d, Type: %s)\n", asset.WidgetName, asset.WidgetID, asset.WidgetType)
			content += fmt.Sprintf("    Hash: %s\n", asset.Hash)
			if asset.OriginalFilename != "" {
				content += fmt.Sprintf("    Original Filename: %s\n", asset.OriginalFilename)
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
func (cmd *DiscoverCommand) generateCSVReport(missingAssets []canvus.AssetInfo) error {
	reportPath := "missing_assets.csv"

	// Generate CSV content
	content := "Hash,WidgetType,OriginalFilename,CanvasID,CanvasName,WidgetID,WidgetName\n"

	for _, asset := range missingAssets {
		content += fmt.Sprintf("%s,%s,%s,%d,%s,%d,%s\n",
			asset.Hash,
			asset.WidgetType,
			asset.OriginalFilename,
			asset.CanvasID,
			asset.CanvasName,
			asset.WidgetID,
			asset.WidgetName,
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
	fmt.Printf("❌ Missing Assets: %d\n", len(missingAssets))

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

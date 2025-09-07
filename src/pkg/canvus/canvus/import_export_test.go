package canvus

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestImportExportRoundTrip verifies that widgets and assets can be exported and imported correctly.
func TestImportExportRoundTrip(t *testing.T) {
	ctx := context.Background()
	// Setup: create a test session, folder, and canvas
	session, cleanupSession, err := getTestUserSession(ctx)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	defer cleanupSession()

	folder, err := session.CreateFolder(ctx, CreateFolderRequest{Name: "test_folder_import_export"})
	if err != nil {
		t.Fatalf("Failed to create test folder: %v", err)
	}
	defer session.DeleteFolder(ctx, folder.ID)

	sourceCanvas, err := session.CreateCanvas(ctx, CreateCanvasRequest{Name: "test_canvas_import_export", FolderID: folder.ID})
	if err != nil {
		t.Fatalf("Failed to create test canvas: %v", err)
	}
	defer session.DeleteCanvas(ctx, sourceCanvas.ID)

	// 1. Create source anchor at (0,0,1920,1080)
	sourceAnchorReq := map[string]interface{}{
		"widget_type": "anchor",
		"location":    map[string]interface{}{"x": 0, "y": 0},
		"size":        map[string]interface{}{"width": 1920, "height": 1080},
		"title":       "source anchor",
	}
	sourceAnchor, err := session.CreateAnchor(ctx, sourceCanvas.ID, sourceAnchorReq)
	if err != nil {
		t.Fatalf("Failed to create source anchor: %v", err)
	}

	// 2. Create note, image, pdf, video within the anchor zone
	// Note
	_, err = session.CreateNote(ctx, sourceCanvas.ID, map[string]interface{}{
		"title": "test note", "widget_type": "note", "text": "hello",
		"location": map[string]interface{}{"x": 100, "y": 100}, "size": map[string]interface{}{"width": 300, "height": 300},
	})
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}
	// Image
	imgFile, err := os.Open("../test_files/test_image.jpg")
	if err != nil {
		t.Fatalf("Failed to open test image: %v", err)
	}
	imgBytes, err := io.ReadAll(imgFile)
	imgFile.Close()
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}
	imgMeta := map[string]interface{}{"title": "test image", "widget_type": "image",
		"location": map[string]interface{}{"x": 200, "y": 200}, "size": map[string]interface{}{"width": 400, "height": 400},
	}
	imgBody, imgContentType, err := buildMultipartBody(imgMeta, "data", "test_image.jpg", imgBytes)
	if err != nil {
		t.Fatalf("Failed to build multipart body for image: %v", err)
	}
	if _, err := session.CreateImage(ctx, sourceCanvas.ID, imgBody, imgContentType); err != nil {
		t.Fatalf("Failed to create image: %v", err)
	}
	// PDF
	pdfFile, err := os.Open("../test_files/test_pdf.pdf")
	if err != nil {
		t.Fatalf("Failed to open test pdf: %v", err)
	}
	pdfBytes, err := io.ReadAll(pdfFile)
	pdfFile.Close()
	if err != nil {
		t.Fatalf("Failed to read test pdf: %v", err)
	}
	pdfMeta := map[string]interface{}{"title": "test pdf", "widget_type": "pdf",
		"location": map[string]interface{}{"x": 300, "y": 300}, "size": map[string]interface{}{"width": 500, "height": 500},
	}
	pdfBody, pdfContentType, err := buildMultipartBody(pdfMeta, "data", "test_pdf.pdf", pdfBytes)
	if err != nil {
		t.Fatalf("Failed to build multipart body for pdf: %v", err)
	}
	if _, err := session.CreatePDF(ctx, sourceCanvas.ID, pdfBody, pdfContentType); err != nil {
		t.Fatalf("Failed to create pdf: %v", err)
	}
	// Video
	videoFile, err := os.Open("../test_files/test_video.mp4")
	if err != nil {
		t.Fatalf("Failed to open test video: %v", err)
	}
	videoBytes, err := io.ReadAll(videoFile)
	videoFile.Close()
	if err != nil {
		t.Fatalf("Failed to read test video: %v", err)
	}
	videoMeta := map[string]interface{}{"title": "test video", "widget_type": "video",
		"location": map[string]interface{}{"x": 400, "y": 400}, "size": map[string]interface{}{"width": 600, "height": 600},
	}
	videoBody, videoContentType, err := buildMultipartBody(videoMeta, "data", "test_video.mp4", videoBytes)
	if err != nil {
		t.Fatalf("Failed to build multipart body for video: %v", err)
	}
	if _, err := session.CreateVideo(ctx, sourceCanvas.ID, videoBody, videoContentType); err != nil {
		t.Fatalf("Failed to create video: %v", err)
	}

	// 3. Use WidgetsContainId to get all widgets contained in the anchor
	zone, err := WidgetsContainId(ctx, session, sourceCanvas.ID, sourceAnchor.ID, nil, 0)
	if err != nil {
		t.Fatalf("WidgetsContainId failed: %v", err)
	}
	var widgetIDs []string
	if zone.Container.ID != "" {
		widgetIDs = append(widgetIDs, zone.Container.ID)
	}
	for _, w := range zone.Contents {
		if w.ID != "" {
			widgetIDs = append(widgetIDs, w.ID)
		}
	}
	// zone.SharedCanvasID is now available for export logic if needed

	// 4. Export those widgets using ExportWidgetsToFolder, with the anchor's bounding box as the region
	anchorRect := WidgetBoundingBox(Widget{Location: sourceAnchor.Location, Size: sourceAnchor.Size})
	exportFolder, err := session.ExportWidgetsToFolder(ctx, sourceCanvas.ID, widgetIDs, anchorRect, zone.SharedCanvasID, "")
	if err != nil {
		t.Fatalf("ExportWidgetsToFolder failed: %v", err)
	}
	exportPath := filepath.Join(exportFolder, "export.json")

	// Import into a new canvas
	importDir := filepath.Join("tests", "importdata")
	os.MkdirAll(importDir, 0755)
	importCanvas, err := session.CreateCanvas(ctx, CreateCanvasRequest{Name: "import_canvas", FolderID: folder.ID})
	if err != nil {
		t.Fatalf("Failed to create import canvas: %v", err)
	}
	defer session.DeleteCanvas(ctx, importCanvas.ID)
	targetAnchorReq := map[string]interface{}{
		"widget_type": "anchor",
		"location":    map[string]interface{}{"x": 10000, "y": 10000},
		"size":        map[string]interface{}{"width": 1280, "height": 720},
		"title":       "target anchor",
	}
	targetAnchor, err := session.CreateAnchor(ctx, importCanvas.ID, targetAnchorReq)
	if err != nil {
		t.Fatalf("Failed to create target anchor: %v", err)
	}
	targetRect := WidgetBoundingBox(Widget{Location: targetAnchor.Location, Size: targetAnchor.Size})

	// 6. Import from the export folder to the target anchor's bounding box
	importBytes, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("Failed to read exported data: %v", err)
	}
	var importedSet ExportedWidgetSet
	err = json.Unmarshal(importBytes, &importedSet)
	if err != nil {
		t.Fatalf("Failed to unmarshal imported data: %v", err)
	}
	_, err = session.ImportWidgetsToRegion(ctx, importCanvas.ID, &importedSet, targetRect)
	if err != nil {
		t.Fatalf("ImportWidgetsToRegion failed: %v", err)
	}

	// 7. Use WidgetsContainId on the target anchor to verify the imported widgets are present
	importedZone, err := WidgetsContainId(ctx, session, importCanvas.ID, targetAnchor.ID, nil, 0)
	if err != nil {
		t.Fatalf("WidgetsContainId (imported) failed: %v", err)
	}
	// Debug: log all widgets found after import
	if testing.Verbose() {
		for i, w := range importedZone.Contents {
			println("[DEBUG] Imported widget", i, "ID=", w.ID, "Type=", w.WidgetType, "ParentID=", w.ParentID)
		}
	}
	println("[DEBUG] Import canvas ID:", importCanvas.ID)
	println("[DEBUG] Imported widget count:", len(importedZone.Contents))
	for _, w := range importedZone.Contents {
		println("[DEBUG] Widget type:", w.WidgetType, "ID:", w.ID)
	}
	if len(importedZone.Contents) < 4 {
		t.Errorf("Expected at least 4 widgets after import, got %d", len(importedZone.Contents))
	}
	// Check that all types are present
	importedTypes := map[string]bool{}
	for _, w := range importedZone.Contents {
		importedTypes[strings.ToLower(w.WidgetType)] = true
	}
	for _, typ := range []string{"note", "image", "pdf", "video"} {
		if !importedTypes[typ] {
			t.Errorf("Missing widget type after import: %s", typ)
		}
	}
}

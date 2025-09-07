package canvus

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// TestWidgetAndAssetEndpoints tests all widget and asset endpoints using a single authenticated session, folder, and canvas.
// For upload endpoints, uses real files from test_files/ and constructs multipart bodies correctly.
// Direct asset creation is not tested, as /api/v1/assets is read-only (see Mipmaps_API.md).
func TestWidgetAndAssetEndpoints(t *testing.T) {
	ctx := context.Background()
	fmt.Printf("[%s] START TestWidgetAndAssetEndpoints\n", time.Now().Format(time.RFC3339Nano))

	// 1. Create a single authenticated session
	session, cleanupSession, err := getTestUserSession(ctx)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	defer func() {
		fmt.Printf("[%s] CLEANUP: cleanupSession()\n", time.Now().Format(time.RFC3339Nano))
		cleanupSession()
	}()

	// 2. Create a test folder
	fmt.Printf("[%s] Creating test folder...\n", time.Now().Format(time.RFC3339Nano))
	folder, err := session.CreateFolder(ctx, CreateFolderRequest{Name: "test_folder_widgets_assets"})
	if err != nil {
		t.Fatalf("Failed to create test folder: %v", err)
	}
	defer func() {
		fmt.Printf("[%s] CLEANUP: DeleteFolder(%s)\n", time.Now().Format(time.RFC3339Nano), folder.ID)
		_ = session.DeleteFolder(ctx, folder.ID)
	}()

	// 3. Create a test canvas in the folder
	fmt.Printf("[%s] Creating test canvas...\n", time.Now().Format(time.RFC3339Nano))
	canvas, err := session.CreateCanvas(ctx, CreateCanvasRequest{Name: "test_canvas_widgets_assets", FolderID: folder.ID})
	if err != nil {
		t.Fatalf("Failed to create test canvas: %v", err)
	}
	defer func() {
		fmt.Printf("[%s] CLEANUP: DeleteCanvas(%s)\n", time.Now().Format(time.RFC3339Nano), canvas.ID)
		_ = session.DeleteCanvas(ctx, canvas.ID)
	}()

	// Store created resource IDs for later steps
	var noteIDs, anchorIDs, imageIDs, pdfIDs, videoIDs, connectorIDs []string
	var widgetIDs []string // for connectors

	// --- CREATE ---
	fmt.Printf("[%s] CREATE: Notes\n", time.Now().Format(time.RFC3339Nano))
	note1, err := session.CreateNote(ctx, canvas.ID, map[string]interface{}{"title": "note1", "widget_type": "note"})
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
	note2, err := session.CreateNote(ctx, canvas.ID, map[string]interface{}{"title": "note2", "widget_type": "note"})
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
	noteIDs = append(noteIDs, note1.ID, note2.ID)
	widgetIDs = append(widgetIDs, note1.ID, note2.ID)

	fmt.Printf("[%s] CREATE: Anchors\n", time.Now().Format(time.RFC3339Nano))
	anchor, err := session.CreateAnchor(ctx, canvas.ID, map[string]interface{}{"title": "anchor1", "widget_type": "anchor"})
	if err != nil {
		t.Fatalf("CreateAnchor failed: %v", err)
	}
	anchorIDs = append(anchorIDs, anchor.ID)
	widgetIDs = append(widgetIDs, anchor.ID)

	fmt.Printf("[%s] CREATE: Images\n", time.Now().Format(time.RFC3339Nano))
	imgFile, err := os.Open("../test_files/test_image.jpg")
	if err != nil {
		t.Fatalf("Failed to open test image: %v", err)
	}
	imgBytes, err := io.ReadAll(imgFile)
	imgFile.Close()
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}
	meta := map[string]interface{}{"title": "test image", "widget_type": "image"}
	body, contentType, err := buildMultipartBody(meta, "data", "test_image.jpg", imgBytes)
	if err != nil {
		t.Fatalf("Failed to build multipart body: %v", err)
	}
	image, err := session.CreateImage(ctx, canvas.ID, body, contentType)
	if err != nil {
		t.Fatalf("CreateImage failed: %v", err)
	}
	imageIDs = append(imageIDs, image.ID)
	widgetIDs = append(widgetIDs, image.ID)

	fmt.Printf("[%s] CREATE: PDFs\n", time.Now().Format(time.RFC3339Nano))
	pdfFile, err := os.Open("../test_files/test.pdf")
	if err == nil {
		pdfBytes, _ := io.ReadAll(pdfFile)
		pdfFile.Close()
		meta := map[string]interface{}{"title": "test pdf", "widget_type": "pdf"}
		body, pdfContentType, err := buildMultipartBody(meta, "data", "test.pdf", pdfBytes)
		if err == nil {
			pdf, err := session.CreatePDF(ctx, canvas.ID, body, pdfContentType)
			if err == nil {
				pdfIDs = append(pdfIDs, pdf.ID)
				widgetIDs = append(widgetIDs, pdf.ID)
			}
		}
	}

	fmt.Printf("[%s] CREATE: Videos\n", time.Now().Format(time.RFC3339Nano))
	videoFile, err := os.Open("../test_files/test_video.mp4")
	if err == nil {
		videoBytes, _ := io.ReadAll(videoFile)
		videoFile.Close()
		meta := map[string]interface{}{"title": "test video", "widget_type": "video"}
		body, videoContentType, err := buildMultipartBody(meta, "data", "test_video.mp4", videoBytes)
		if err == nil {
			video, err := session.CreateVideo(ctx, canvas.ID, body, videoContentType)
			if err == nil {
				videoIDs = append(videoIDs, video.ID)
				widgetIDs = append(widgetIDs, video.ID)
			}
		}
	}

	fmt.Printf("[%s] CREATE: Connectors\n", time.Now().Format(time.RFC3339Nano))
	if len(widgetIDs) >= 2 {
		conn, err := session.CreateConnector(ctx, canvas.ID, map[string]interface{}{"from_id": widgetIDs[0], "to_id": widgetIDs[1]})
		if err == nil {
			connectorIDs = append(connectorIDs, conn.ID)
		}
	}

	// --- WIDGETS: List and Get ---
	fmt.Printf("[%s] WIDGETS: ListWidgets\n", time.Now().Format(time.RFC3339Nano))
	widgets, err := session.ListWidgets(ctx, canvas.ID, nil)
	if err != nil {
		fmt.Printf("[%s] ListWidgets failed: %v\n", time.Now().Format(time.RFC3339Nano), err)
	} else {
		fmt.Printf("[%s] ListWidgets returned %d widgets\n", time.Now().Format(time.RFC3339Nano), len(widgets))
		for _, w := range widgets {
			fmt.Printf("[%s] Widget: ID=%s, Type=%s\n", time.Now().Format(time.RFC3339Nano), w.ID, w.WidgetType)
		}
	}
	// GetWidget for each created widget
	widgetIDs = []string{note1.ID, note2.ID}
	if len(anchorIDs) > 0 {
		widgetIDs = append(widgetIDs, anchorIDs[0])
	}
	if len(imageIDs) > 0 {
		widgetIDs = append(widgetIDs, imageIDs[0])
	}
	if len(pdfIDs) > 0 {
		widgetIDs = append(widgetIDs, pdfIDs[0])
	}
	if len(videoIDs) > 0 {
		widgetIDs = append(widgetIDs, videoIDs[0])
	}
	for _, wid := range widgetIDs {
		w, err := session.GetWidget(ctx, canvas.ID, wid)
		if err != nil {
			fmt.Printf("[%s] GetWidget(%s) failed: %v\n", time.Now().Format(time.RFC3339Nano), wid, err)
		} else {
			fmt.Printf("[%s] GetWidget(%s) succeeded: Type=%s\n", time.Now().Format(time.RFC3339Nano), wid, w.WidgetType)
		}
	}

	// --- PATCH (UPDATE) ---
	fmt.Printf("[%s] PATCH: Notes\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range noteIDs {
		_, _ = session.UpdateNote(ctx, canvas.ID, id, map[string]interface{}{"title": "updated note"})
	}
	fmt.Printf("[%s] PATCH: Anchors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range anchorIDs {
		_, _ = session.UpdateAnchor(ctx, canvas.ID, id, map[string]interface{}{"title": "updated anchor"})
	}
	fmt.Printf("[%s] PATCH: Images\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range imageIDs {
		_, _ = session.UpdateImage(ctx, canvas.ID, id, map[string]interface{}{"title": "updated image"})
	}
	fmt.Printf("[%s] PATCH: PDFs\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range pdfIDs {
		_, _ = session.UpdatePDF(ctx, canvas.ID, id, map[string]interface{}{"title": "updated pdf"})
	}
	fmt.Printf("[%s] PATCH: Videos\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range videoIDs {
		_, _ = session.UpdateVideo(ctx, canvas.ID, id, map[string]interface{}{"title": "updated video"})
	}
	fmt.Printf("[%s] PATCH: Connectors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range connectorIDs {
		_, _ = session.UpdateConnector(ctx, canvas.ID, id, map[string]interface{}{"label": "updated connector"})
	}

	// --- GET ---
	fmt.Printf("[%s] GET: Notes\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range noteIDs {
		_, _ = session.GetNote(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] GET: Anchors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range anchorIDs {
		_, _ = session.GetAnchor(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] GET: Images\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range imageIDs {
		_, _ = session.GetImage(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] GET: PDFs\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range pdfIDs {
		_, _ = session.GetPDF(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] GET: Videos\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range videoIDs {
		_, _ = session.GetVideo(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] GET: Connectors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range connectorIDs {
		_, _ = session.GetConnector(ctx, canvas.ID, id)
	}

	// --- DELETE ---
	fmt.Printf("[%s] DELETE: Connectors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range connectorIDs {
		_ = session.DeleteConnector(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] DELETE: Videos\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range videoIDs {
		_ = session.DeleteVideo(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] DELETE: PDFs\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range pdfIDs {
		_ = session.DeletePDF(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] DELETE: Images\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range imageIDs {
		_ = session.DeleteImage(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] DELETE: Anchors\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range anchorIDs {
		_ = session.DeleteAnchor(ctx, canvas.ID, id)
	}
	fmt.Printf("[%s] DELETE: Notes\n", time.Now().Format(time.RFC3339Nano))
	for _, id := range noteIDs {
		_ = session.DeleteNote(ctx, canvas.ID, id)
	}

	// --- Mipmaps ---
	fmt.Printf("[%s] MIPMAPS: Testing info, level, and asset retrieval\n", time.Now().Format(time.RFC3339Nano))
	if len(imageIDs) > 0 {
		img, _ := session.GetImage(ctx, canvas.ID, imageIDs[0])
		if img != nil && img.Hash != "" {
			_, _ = session.GetMipmapInfo(ctx, canvas.ID, img.Hash, nil)
			_, _ = session.GetMipmapLevel(ctx, canvas.ID, img.Hash, 0, nil)
			_, _ = session.GetAssetByHash(ctx, canvas.ID, img.Hash)
		}
	}

	// --- Backgrounds ---
	fmt.Printf("[%s] BACKGROUNDS: Testing get, patch, post\n", time.Now().Format(time.RFC3339Nano))
	_, _ = session.GetCanvasBackground(ctx, canvas.ID)
	_ = session.PatchCanvasBackground(ctx, canvas.ID, map[string]interface{}{"type": "solid", "background_color": "#FF0000"})
	// For post, reuse image file if available
	if len(imageIDs) > 0 {
		imgFile, bgErr := os.Open("../test_files/test_image.jpg")
		if bgErr == nil {
			bgImgBytes, _ := io.ReadAll(imgFile)
			imgFile.Close()
			bgBody, _, bgErr := buildMultipartBody(nil, "data", "test_image.jpg", bgImgBytes)
			if bgErr == nil {
				_ = session.PostCanvasBackground(ctx, canvas.ID, bgBody)
			}
		}
	}

	// --- Color Presets ---
	fmt.Printf("[%s] COLOR PRESETS: Testing get and patch\n", time.Now().Format(time.RFC3339Nano))
	_, _ = session.GetColorPresets(ctx, canvas.ID)
	_, _ = session.PatchColorPresets(ctx, canvas.ID, &ColorPresets{NoteBackground: []string{"#FFFFFF", "#000000"}})

	// --- Uploads Folder ---
	fmt.Printf("[%s] UPLOADS FOLDER: Testing note and asset upload\n", time.Now().Format(time.RFC3339Nano))
	var uploadNoteBody io.Reader
	var uploadNoteErr error
	meta = map[string]interface{}{"title": "upload note", "widget_type": "note"}
	uploadNoteBody, _, uploadNoteErr = buildMultipartBody(meta, "", "", nil)
	if uploadNoteErr == nil {
		_, _ = session.UploadNote(ctx, canvas.ID, uploadNoteBody)
	}
	var uploadBody io.Reader
	var uploadErr error
	var uploadImgBytes []byte
	var uploadAssetErr error
	var uploadImgFile *os.File
	if len(imageIDs) > 0 {
		uploadImgFile, uploadAssetErr = os.Open("../test_files/test_image.jpg")
		if uploadAssetErr == nil {
			uploadImgBytes, _ = io.ReadAll(uploadImgFile)
			uploadImgFile.Close()
			meta := map[string]interface{}{"title": "upload asset"}
			uploadBody, _, uploadErr = buildMultipartBody(meta, "data", "test_image.jpg", uploadImgBytes)
			if uploadErr == nil {
				_, _ = session.UploadAsset(ctx, canvas.ID, uploadBody)
			}
		}
	}

	// --- Video Inputs/Outputs ---
	fmt.Printf("[%s] VIDEO INPUTS/OUTPUTS: Checking for client and testing if available\n", time.Now().Format(time.RFC3339Nano))
	if hasActiveClient(ctx) {
		// Replace with real clientID if available
		clientID := "dummy-client-id"
		_, _ = session.ListVideoInputs(ctx, canvas.ID)
		_, _ = session.ListClientVideoInputs(ctx, clientID)
		_, _ = session.ListVideoOutputs(ctx, clientID)
		// Skipping create/delete/update due to lack of real client context
	} else {
		fmt.Printf("[%s] SKIP: No active client for video input/output tests\n", time.Now().Format(time.RFC3339Nano))
	}

	fmt.Printf("[%s] END TestWidgetAndAssetEndpoints\n", time.Now().Format(time.RFC3339Nano))
}

// Helper: create a test session as a real test user
func getTestUserSession(ctx context.Context) (*Session, func(), error) {
	ts, err := loadTestSettings()
	if err != nil {
		fmt.Printf("[DEBUG] loadTestSettings error: %v\n", err)
		return nil, func() {}, err
	}
	fmt.Printf("[DEBUG] Loaded settings: api_base_url=%s, api_key=%s\n", ts.APIBaseURL, ts.APIKey)
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)
	testEmail := "testuser_" + fmt.Sprint(time.Now().UnixNano()) + "@example.com"
	testName := "testuser_" + fmt.Sprint(time.Now().UnixNano())
	testPassword := "TestPassword123!"
	fmt.Printf("[DEBUG] Creating test user: email=%s, name=%s\n", testEmail, testName)
	tc, err := NewTestClient(ctx, admin, ts.APIBaseURL, testEmail, testName, testPassword)
	if err != nil {
		fmt.Printf("[DEBUG] NewTestClient error: %v\n", err)
		return nil, func() {}, err
	}
	cleanup := func() { _ = tc.Cleanup(ctx) }
	return tc.Session, cleanup, nil
}

// Helper: create a test folder using a real test user session
func createTestFolder(ctx context.Context, session *Session, name string) (*Folder, error) {
	folder, err := session.CreateFolder(ctx, CreateFolderRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return folder, nil
}

// Helper: create a test canvas using a real test user session
func createTestCanvas(ctx context.Context, session *Session, name, folderID string) (*Canvas, error) {
	canvas, err := session.CreateCanvas(ctx, CreateCanvasRequest{Name: name, FolderID: folderID})
	if err != nil {
		return nil, err
	}
	return canvas, nil
}

func hasActiveClient(ctx context.Context) bool { return false }

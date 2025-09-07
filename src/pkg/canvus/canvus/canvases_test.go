package canvus

import (
	"context"
	"testing"
	"time"
)

// TestClient is used for all Canvas API tests except those requiring admin rights.
// Only admin operations (user create/delete, unblock, approve, audit-log, licenses) use the admin client.
// See PRD.md for client type documentation.

func TestCanvasLifecycle(t *testing.T) {
	ctx := context.Background()
	// Use TestClient for all non-admin tests
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)

	testEmail := "testcanvas_" + time.Now().Format("20060102150405") + "@example.com"
	testName := "testcanvas_" + time.Now().Format("150405")
	testPassword := "TestPassword123!"
	tc, err := NewTestClient(ctx, admin, ts.APIBaseURL, testEmail, testName, testPassword)
	if err != nil {
		t.Fatalf("failed to create TestClient: %v", err)
	}
	defer func() { _ = tc.Cleanup(ctx) }()
	session := tc.Session

	// Create a test canvas
	canvasName := "testcanvas_" + time.Now().Format("20060102150405")
	t.Logf("[CreateCanvas] Sending: Name=%q", canvasName)
	canvas, err := session.CreateCanvas(ctx, CreateCanvasRequest{Name: canvasName})
	if err != nil {
		t.Fatalf("failed to create canvas: %v", err)
	}
	// Clean up: delete the canvas at the end
	defer func() { _ = session.DeleteCanvas(ctx, canvas.ID) }()

	// List canvases and check the new canvas is present
	t.Logf("[ListCanvases] Listing all canvases")
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		t.Errorf("failed to list canvases: %v", err)
	}
	found := false
	for _, c := range canvases {
		if c.ID == canvas.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created canvas not found in list")
	}

	// Get the canvas by ID
	t.Logf("[GetCanvas] ID=%q", canvas.ID)
	got, err := session.GetCanvas(ctx, canvas.ID)
	if err != nil {
		t.Errorf("failed to get canvas: %v", err)
	}
	if got.Name != canvasName {
		t.Errorf("expected canvas name %q, got %q", canvasName, got.Name)
	}

	// Update the canvas name (PATCH /canvases/:id)
	newName := canvasName + "_renamed"
	t.Logf("[UpdateCanvas] ID=%q, Name=%q", canvas.ID, newName)
	updated, err := session.UpdateCanvas(ctx, canvas.ID, UpdateCanvasRequest{Name: newName})
	if err != nil {
		t.Errorf("failed to update canvas name: %v", err)
	} else {
		if updated.Name != newName {
			t.Errorf("expected updated canvas name %q, got %q", newName, updated.Name)
		}
	}

	// Update the canvas mode (PATCH /canvases/:id)
	t.Logf("[UpdateCanvas] ID=%q, Mode=%q", canvas.ID, "demo")
	updatedMode, err := session.UpdateCanvas(ctx, canvas.ID, UpdateCanvasRequest{Mode: "demo"})
	if err != nil {
		t.Errorf("failed to update canvas mode: %v", err)
	} else {
		if updatedMode.Mode != "demo" {
			t.Errorf("expected mode 'demo', got %q", updatedMode.Mode)
		}
	}

	// Save and restore demo state (should not error)
	t.Logf("[SaveDemoState] ID=%q", canvas.ID)
	err = session.SaveDemoState(ctx, canvas.ID)
	if err != nil {
		t.Errorf("failed to save demo state: %v", err)
	}
	t.Logf("[RestoreDemoCanvas] ID=%q", canvas.ID)
	err = session.RestoreDemoCanvas(ctx, canvas.ID)
	if err != nil {
		t.Errorf("failed to restore demo canvas: %v", err)
	}

	// Get preview (may be empty, but should not error)
	t.Logf("[GetCanvasPreview] ID=%q", canvas.ID)
	_, err = session.GetCanvasPreview(ctx, canvas.ID)
	if err != nil {
		t.Logf("GetCanvasPreview: %v (may be expected if no preview)", err)
	}

	// Permissions: get and set
	t.Logf("[GetCanvasPermissions] ID=%q", canvas.ID)
	perms, err := session.GetCanvasPermissions(ctx, canvas.ID)
	if err != nil {
		t.Errorf("failed to get canvas permissions: %v", err)
	}
	perms.EditorsCanShare = false
	perms.LinkPermission = "view"
	t.Logf("[SetCanvasPermissions] ID=%q, LinkPermission=%q", canvas.ID, perms.LinkPermission)
	updatedPerms, err := session.SetCanvasPermissions(ctx, canvas.ID, *perms)
	if err != nil {
		t.Errorf("failed to set canvas permissions: %v", err)
	}
	if updatedPerms.LinkPermission != "view" {
		t.Errorf("expected LinkPermission 'view', got %q", updatedPerms.LinkPermission)
	}

	// Move/copy/trash: create a second folder, move/copy canvas, then trash
	folderName := "testfolder_for_canvas_" + time.Now().Format("20060102150405")
	t.Logf("[CreateFolder] Name=%q", folderName)
	folder, err := session.CreateFolder(ctx, CreateFolderRequest{Name: folderName})
	if err != nil {
		t.Fatalf("failed to create folder for move/copy: %v", err)
	}
	defer func() { _ = session.DeleteFolder(ctx, folder.ID) }()

	// Move canvas
	t.Logf("[MoveCanvas] CanvasID=%q, FolderID=%q", canvas.ID, folder.ID)
	moved, err := session.MoveCanvas(ctx, canvas.ID, MoveOrCopyCanvasRequest{FolderID: folder.ID, Conflicts: "replace"})
	if err != nil {
		t.Errorf("failed to move canvas: %v", err)
	} else if moved.FolderID != folder.ID {
		t.Errorf("expected moved canvas folder_id %q, got %q", folder.ID, moved.FolderID)
	}

	// Copy canvas
	t.Logf("[CopyCanvas] CanvasID=%q, FolderID=%q", canvas.ID, folder.ID)
	copied, err := session.CopyCanvas(ctx, canvas.ID, MoveOrCopyCanvasRequest{FolderID: folder.ID, Conflicts: "replace"})
	if err != nil {
		t.Errorf("failed to copy canvas: %v", err)
	}
	defer func() { _ = session.DeleteCanvas(ctx, copied.ID) }()

	// Trash the copied canvas
	t.Logf("[TrashCanvas] CanvasID=%q, TrashID=%q", copied.ID, "trash."+folder.ID)
	trashed, err := session.TrashCanvas(ctx, copied.ID, "trash."+folder.ID)
	if err != nil {
		t.Errorf("failed to trash canvas: %v", err)
	}
	if !trashed.InTrash {
		t.Errorf("expected canvas to be in trash, got in_trash=%v", trashed.InTrash)
	}
}

func TestCanvasInvalidCases(t *testing.T) {
	ctx := context.Background()
	// Use TestClient for all non-admin tests
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)

	// Get non-existent canvas
	_, err = admin.GetCanvas(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for non-existent canvas, got nil")
	}

	// Delete non-existent canvas
	err = admin.DeleteCanvas(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for deleting non-existent canvas, got nil")
	}

	// Update non-existent canvas
	_, err = admin.UpdateCanvas(ctx, "nonexistent-canvas-id", UpdateCanvasRequest{Name: "newname"})
	if err == nil {
		t.Errorf("expected error for updating non-existent canvas, got nil")
	}

	// Move non-existent canvas
	_, err = admin.MoveCanvas(ctx, "nonexistent-canvas-id", MoveOrCopyCanvasRequest{FolderID: "some-folder-id"})
	if err == nil {
		t.Errorf("expected error for moving non-existent canvas, got nil")
	}

	// Copy non-existent canvas
	_, err = admin.CopyCanvas(ctx, "nonexistent-canvas-id", MoveOrCopyCanvasRequest{FolderID: "some-folder-id"})
	if err == nil {
		t.Errorf("expected error for copying non-existent canvas, got nil")
	}

	// Trash non-existent canvas
	_, err = admin.TrashCanvas(ctx, "nonexistent-canvas-id", "trash.some-folder-id")
	if err == nil {
		t.Errorf("expected error for trashing non-existent canvas, got nil")
	}

	// Get permissions of non-existent canvas
	_, err = admin.GetCanvasPermissions(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for getting permissions of non-existent canvas, got nil")
	}

	// Set permissions of non-existent canvas
	perms := CanvasPermissions{EditorsCanShare: false}
	_, err = admin.SetCanvasPermissions(ctx, "nonexistent-canvas-id", perms)
	if err == nil {
		t.Errorf("expected error for setting permissions of non-existent canvas, got nil")
	}

	// Save/restore demo state for non-existent canvas
	err = admin.SaveDemoState(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for saving demo state of non-existent canvas, got nil")
	}
	err = admin.RestoreDemoCanvas(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for restoring demo state of non-existent canvas, got nil")
	}

	// Get preview for non-existent canvas
	_, err = admin.GetCanvasPreview(ctx, "nonexistent-canvas-id")
	if err == nil {
		t.Errorf("expected error for getting preview of non-existent canvas, got nil")
	}
}

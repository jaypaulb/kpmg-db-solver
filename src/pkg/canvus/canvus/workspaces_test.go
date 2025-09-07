package canvus

import (
	"context"
	"testing"
)

func TestWorkspaceLifecycle(t *testing.T) {
	ctx := context.Background()
	// Use the same test settings loader as other integration tests
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Use a user client with permissions from settings.json
	uc, err := NewUserClient(ctx, ts.APIBaseURL, ts.TestUser.Username, ts.TestUser.Password)
	if err != nil {
		t.Fatalf("failed to create UserClient: %v", err)
	}
	defer func() { _ = uc.Cleanup(ctx) }()
	session := uc.Session

	// List clients and pick the first one (or skip if none)
	clients, err := session.ListClients(ctx)
	if err != nil || len(clients) == 0 {
		t.Skip("No clients available for workspace tests")
	}
	clientID := clients[0].ID

	t.Run("ListWorkspaces", func(t *testing.T) {
		workspaces, err := session.ListWorkspaces(ctx, clientID)
		if err != nil {
			t.Fatalf("ListWorkspaces failed: %v", err)
		}
		if len(workspaces) == 0 {
			t.Skip("No workspaces available on client for further tests")
		}
	})

	t.Run("GetWorkspace", func(t *testing.T) {
		selector := WorkspaceSelector{Index: intPtr(0)}
		ws, err := session.GetWorkspace(ctx, clientID, selector)
		if err != nil {
			t.Fatalf("GetWorkspace failed: %v", err)
		}
		if ws.Index != 0 {
			t.Errorf("expected workspace index 0, got %d", ws.Index)
		}
	})

	t.Run("ToggleInfoPanel", func(t *testing.T) {
		selector := WorkspaceSelector{Index: intPtr(0)}
		err := session.ToggleWorkspaceInfoPanel(ctx, clientID, selector)
		if err != nil {
			t.Errorf("ToggleWorkspaceInfoPanel failed: %v", err)
		}
	})

	t.Run("TogglePinned", func(t *testing.T) {
		selector := WorkspaceSelector{Index: intPtr(0)}
		err := session.ToggleWorkspacePinned(ctx, clientID, selector)
		if err != nil {
			t.Errorf("ToggleWorkspacePinned failed: %v", err)
		}
	})

	t.Run("UpdateWorkspaceViewport", func(t *testing.T) {
		selector := WorkspaceSelector{Index: intPtr(0)}
		ws, err := session.GetWorkspace(ctx, clientID, selector)
		if err != nil {
			t.Fatalf("GetWorkspace for viewport update failed: %v", err)
		}
		if ws.Size == nil {
			t.Skip("Workspace has no size info; skipping viewport test")
		}
		rect := &Rectangle{
			X:      0,
			Y:      0,
			Width:  ws.Size.Width,
			Height: ws.Size.Height,
		}
		_, err = session.UpdateWorkspace(ctx, clientID, selector, UpdateWorkspaceRequest{ViewRectangle: rect})
		if err != nil {
			t.Errorf("UpdateWorkspace (viewport) failed: %v", err)
		}
	})

	// TODO: Add OpenCanvasOnWorkspace and widget-centric viewport tests when widget endpoints are implemented
}

func intPtr(i int) *int { return &i }

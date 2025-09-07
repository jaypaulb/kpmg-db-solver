package canvus

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// WorkspaceWidgetGetter allows getting a widget by ID for workspace viewport logic.
type WorkspaceWidgetGetter interface {
	GetWidget(ctx context.Context, clientID string, widgetID string) (*Widget, error)
}

// resolveWorkspaceIndex resolves a workspace index from a WorkspaceSelector.
func (c *Session) resolveWorkspaceIndex(ctx context.Context, clientID string, selector WorkspaceSelector) (int, error) {
	if selector.Index != nil {
		return *selector.Index, nil
	}
	workspaces, err := c.ListWorkspaces(ctx, clientID)
	if err != nil {
		return 0, err
	}
	if selector.Name != nil {
		for _, ws := range workspaces {
			if ws.WorkspaceName == *selector.Name {
				return ws.Index, nil
			}
		}
		return 0, fmt.Errorf("workspace with name %q not found", *selector.Name)
	}
	if selector.User != nil {
		for _, ws := range workspaces {
			if ws.User == *selector.User {
				return ws.Index, nil
			}
		}
		return 0, fmt.Errorf("workspace for user %q not found", *selector.User)
	}
	// Default to index 0
	return 0, nil
}

// ListWorkspaces retrieves all workspaces for a client.
func (c *Session) ListWorkspaces(ctx context.Context, clientID string) ([]Workspace, error) {
	var workspaces []Workspace
	endpoint := fmt.Sprintf("clients/%s/workspaces", clientID)
	err := c.doRequest(ctx, "GET", endpoint, nil, &workspaces, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListWorkspaces: %w", err)
	}
	return workspaces, nil
}

// GetWorkspace retrieves a single workspace by index.
func (c *Session) GetWorkspace(ctx context.Context, clientID string, selector WorkspaceSelector) (*Workspace, error) {
	idx, err := c.resolveWorkspaceIndex(ctx, clientID, selector)
	if err != nil {
		return nil, err
	}
	var ws Workspace
	endpoint := fmt.Sprintf("clients/%s/workspaces/%d", clientID, idx)
	err = c.doRequest(ctx, "GET", endpoint, nil, &ws, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetWorkspace: %w", err)
	}
	return &ws, nil
}

// UpdateWorkspace updates workspace parameters.
func (c *Session) UpdateWorkspace(ctx context.Context, clientID string, selector WorkspaceSelector, req UpdateWorkspaceRequest) (*Workspace, error) {
	idx, err := c.resolveWorkspaceIndex(ctx, clientID, selector)
	if err != nil {
		return nil, err
	}
	var ws Workspace
	endpoint := fmt.Sprintf("clients/%s/workspaces/%d", clientID, idx)
	err = c.doRequest(ctx, "PATCH", endpoint, req, &ws, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateWorkspace: %w", err)
	}
	return &ws, nil
}

// ToggleWorkspaceInfoPanel toggles the info_panel_visible state.
func (c *Session) ToggleWorkspaceInfoPanel(ctx context.Context, clientID string, selector WorkspaceSelector) error {
	ws, err := c.GetWorkspace(ctx, clientID, selector)
	if err != nil {
		return err
	}
	newVal := !ws.InfoPanelVisible
	_, err = c.UpdateWorkspace(ctx, clientID, selector, UpdateWorkspaceRequest{InfoPanelVisible: &newVal})
	return err
}

// ToggleWorkspacePinned toggles the pinned state.
func (c *Session) ToggleWorkspacePinned(ctx context.Context, clientID string, selector WorkspaceSelector) error {
	ws, err := c.GetWorkspace(ctx, clientID, selector)
	if err != nil {
		return err
	}
	newVal := !ws.Pinned
	_, err = c.UpdateWorkspace(ctx, clientID, selector, UpdateWorkspaceRequest{Pinned: &newVal})
	return err
}

// SetWorkspaceViewport sets the workspace viewport for a client by coordinates or by centering on a widget.
//
// Parameters:
//
//	ctx      - context for cancellation and deadlines
//	client   - a WorkspaceWidgetGetter (usually a Session) to fetch widgets
//	apiClient- the Session (API client) to update the workspace
//	clientID - the ID of the client whose workspace is being updated
//	selector - WorkspaceSelector to identify the workspace (by index, name, or user)
//	opts     - SetViewportOptions specifying either WidgetID (to center on a widget) or X, Y, Width, Height (to set explicit viewport)
//
// Returns an error if the operation fails.
//
// Usage Example:
//
//	err := canvus.SetWorkspaceViewport(ctx, session, session, "client123", canvus.WorkspaceSelector{Index: ptr(0)}, canvus.SetViewportOptions{WidgetID: ptr("widget456")})
//	if err != nil {
//	    log.Fatal(err)
//	}
func SetWorkspaceViewport(ctx context.Context, client WorkspaceWidgetGetter, apiClient *Session, clientID string, selector WorkspaceSelector, opts SetViewportOptions) error {
	var rect *Rectangle
	if opts.WidgetID != nil {
		widget, err := client.GetWidget(ctx, clientID, *opts.WidgetID)
		if err != nil {
			return err
		}
		margin := opts.Margin
		if margin == 0 {
			margin = 20 // default margin
		}
		rect = &Rectangle{
			X:      widget.Location.X - margin,
			Y:      widget.Location.Y - margin,
			Width:  widget.Size.Width + 2*margin,
			Height: widget.Size.Height + 2*margin,
		}
	} else if opts.X != nil && opts.Y != nil && opts.Width != nil && opts.Height != nil {
		rect = &Rectangle{
			X:      *opts.X,
			Y:      *opts.Y,
			Width:  *opts.Width,
			Height: *opts.Height,
		}
	} else {
		return errors.New("must provide either WidgetID or all of X, Y, Width, Height")
	}
	_, err := apiClient.UpdateWorkspace(ctx, clientID, selector, UpdateWorkspaceRequest{ViewRectangle: rect})
	return err
}

// OpenCanvasOnWorkspace opens a canvas and optionally centers viewport.
func (c *Session) OpenCanvasOnWorkspace(ctx context.Context, clientID string, selector WorkspaceSelector, opts OpenCanvasOptions) error {
	idx, err := c.resolveWorkspaceIndex(ctx, clientID, selector)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("clients/%s/workspaces/%d/open-canvas", clientID, idx)
	payload := map[string]interface{}{
		"canvas_id": opts.CanvasID,
	}
	if opts.ServerID != "" {
		payload["server_id"] = opts.ServerID
	}
	if opts.UserEmail != "" {
		payload["user_email"] = opts.UserEmail
	}
	// Open canvas
	err = c.doRequest(ctx, "POST", endpoint, payload, nil, nil, false)
	if err != nil {
		return fmt.Errorf("OpenCanvasOnWorkspace: %w", err)
	}

	// Poll for workspace update (canvas ID change)
	timeout := 10 * time.Second
	interval := 200 * time.Millisecond
	if opts.PollTimeout > 0 {
		timeout = opts.PollTimeout
	}
	if opts.PollInterval > 0 {
		interval = opts.PollInterval
	}
	var ws *Workspace
	start := time.Now()
	for {
		ws, err = c.GetWorkspace(ctx, clientID, selector)
		if err != nil {
			return fmt.Errorf("OpenCanvasOnWorkspace: polling GetWorkspace failed: %w", err)
		}
		if ws.CanvasID == opts.CanvasID {
			break
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("OpenCanvasOnWorkspace: timed out waiting for canvas ID %s (last seen: %s)", opts.CanvasID, ws.CanvasID)
		}
		time.Sleep(interval)
	}

	// Optionally set viewport
	if opts.CenterX != nil && opts.CenterY != nil {
		rect := &Rectangle{
			X:      *opts.CenterX,
			Y:      *opts.CenterY,
			Width:  ws.Size.Width,
			Height: ws.Size.Height,
		}
		_, err = c.UpdateWorkspace(ctx, clientID, selector, UpdateWorkspaceRequest{ViewRectangle: rect})
		if err != nil {
			return fmt.Errorf("OpenCanvasOnWorkspace: failed to set viewport: %w", err)
		}
	} else if opts.WidgetID != nil {
		err := SetWorkspaceViewport(ctx, c, c, clientID, selector, SetViewportOptions{WidgetID: opts.WidgetID})
		if err != nil {
			return fmt.Errorf("OpenCanvasOnWorkspace: failed to center on widget: %w", err)
		}
	}
	return nil
}

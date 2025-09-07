package canvus

import (
	"context"
	"fmt"
	"io"
)

// ListWidgets retrieves all widgets for a given canvas. If filter is non-nil, results are filtered client-side.
// This endpoint is read-only: POST, PATCH, DELETE are not supported on /canvases/{id}/widgets.
func (s *Session) ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error) {
	var widgets []Widget
	path := fmt.Sprintf("canvases/%s/widgets", canvasID)
	err := s.doRequest(ctx, "GET", path, nil, &widgets, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListWidgets: %w", err)
	}
	if filter != nil {
		widgets = FilterSlice(widgets, filter)
	}
	return widgets, nil
}

// GetWidget retrieves a widget by ID for a given canvas.
func (s *Session) GetWidget(ctx context.Context, canvasID, widgetID string) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	err := s.doRequest(ctx, "GET", path, nil, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetWidget: %w", err)
	}
	return &widget, nil
}

// CreateWidget creates a widget of the specified type by dispatching to the correct resource-specific method.
// Supported widget types: "note", "anchor", "image", "pdf", "video", "connector".
// For image/pdf/video, req must be a multipart body (io.Reader) and contentType must be provided.
// For other types, req is a map[string]interface{}.
// Returns the created widget as *Widget, or an error if the type is unsupported or the operation fails.
func (s *Session) CreateWidget(ctx context.Context, canvasID string, req interface{}, contentType ...string) (*Widget, error) {
	// For multipart types, req is io.Reader and contentType must be set
	if m, ok := req.(map[string]interface{}); ok {
		widgetType, ok := m["widget_type"].(string)
		if !ok {
			return nil, fmt.Errorf("CreateWidget: missing or invalid widget_type in request")
		}
		switch widgetType {
		case "note":
			note, err := s.CreateNote(ctx, canvasID, m)
			if err != nil {
				return nil, err
			}
			return &Widget{
				ID:         note.ID,
				WidgetType: note.WidgetType,
				ParentID:   note.ParentID,
				Location:   note.Location,
				Size:       note.Size,
				Pinned:     note.Pinned,
				Scale:      note.Scale,
				State:      note.State,
				Depth:      note.Depth,
			}, nil
		case "anchor":
			anchor, err := s.CreateAnchor(ctx, canvasID, m)
			if err != nil {
				return nil, err
			}
			return &Widget{
				ID:         anchor.ID,
				WidgetType: anchor.WidgetType,
				ParentID:   anchor.ParentID,
				Location:   anchor.Location,
				Size:       anchor.Size,
				Pinned:     anchor.Pinned,
				Scale:      anchor.Scale,
				State:      anchor.State,
				Depth:      anchor.Depth,
			}, nil
		case "connector":
			connector, err := s.CreateConnector(ctx, canvasID, m)
			if err != nil {
				return nil, err
			}
			return &Widget{ID: connector.ID, WidgetType: "connector"}, nil
		case "image", "pdf", "video":
			return nil, fmt.Errorf("CreateWidget: for widget_type %q, req must be a multipart body (io.Reader) and contentType must be set", widgetType)
		default:
			return nil, fmt.Errorf("CreateWidget: unsupported widget_type: %s", widgetType)
		}
	} else if rdr, ok := req.(io.Reader); ok {
		if len(contentType) == 0 {
			return nil, fmt.Errorf("CreateWidget: contentType must be provided for multipart widget creation")
		}
		ct := contentType[0]
		// For multipart types, require a map[string]interface{} with widget_type in the JSON part
		// The caller must build the multipart body correctly
		// We cannot infer widget_type from the body, so the caller must know which type is being created
		// We'll try all three and return the first that succeeds
		// (In practice, the CLI/SDK should know which one is being created)
		image, err := s.CreateImage(ctx, canvasID, rdr, ct)
		if err == nil {
			return &Widget{
				ID:         image.ID,
				WidgetType: image.WidgetType,
				ParentID:   image.ParentID,
				Location:   image.Location,
				Size:       image.Size,
				Pinned:     image.Pinned,
				Scale:      image.Scale,
				State:      image.State,
				Depth:      image.Depth,
			}, nil
		}
		pdf, err := s.CreatePDF(ctx, canvasID, rdr, ct)
		if err == nil {
			return &Widget{ID: pdf.ID, WidgetType: "pdf"}, nil
		}
		video, err := s.CreateVideo(ctx, canvasID, rdr, ct)
		if err == nil {
			return &Widget{ID: video.ID, WidgetType: "video"}, nil
		}
		return nil, fmt.Errorf("CreateWidget: failed to create image, pdf, or video: %w", err)
	}
	return nil, fmt.Errorf("CreateWidget: req must be a map[string]interface{} (for note, anchor, connector) or io.Reader (for image, pdf, video)")
}

// UpdateWidget updates a widget of the specified type by dispatching to the correct resource-specific method.
// The req must be a map[string]interface{} with a "widget_type" key.
func (s *Session) UpdateWidget(ctx context.Context, canvasID, widgetID string, req map[string]interface{}) (*Widget, error) {
	widgetType, ok := req["widget_type"].(string)
	if !ok {
		return nil, fmt.Errorf("UpdateWidget: missing or invalid widget_type in request")
	}
	switch widgetType {
	case "note":
		note, err := s.UpdateNote(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{
			ID:         note.ID,
			WidgetType: note.WidgetType,
			ParentID:   note.ParentID,
			Location:   note.Location,
			Size:       note.Size,
			Pinned:     note.Pinned,
			Scale:      note.Scale,
			State:      note.State,
			Depth:      note.Depth,
		}, nil
	case "anchor":
		anchor, err := s.UpdateAnchor(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{
			ID:         anchor.ID,
			WidgetType: anchor.WidgetType,
			ParentID:   anchor.ParentID,
			Location:   anchor.Location,
			Size:       anchor.Size,
			Pinned:     anchor.Pinned,
			Scale:      anchor.Scale,
			State:      anchor.State,
			Depth:      anchor.Depth,
		}, nil
	case "image":
		image, err := s.UpdateImage(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{
			ID:         image.ID,
			WidgetType: image.WidgetType,
			ParentID:   image.ParentID,
			Location:   image.Location,
			Size:       image.Size,
			Pinned:     image.Pinned,
			Scale:      image.Scale,
			State:      image.State,
			Depth:      image.Depth,
		}, nil
	case "pdf":
		pdf, err := s.UpdatePDF(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{ID: pdf.ID, WidgetType: "pdf"}, nil
	case "video":
		video, err := s.UpdateVideo(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{ID: video.ID, WidgetType: "video"}, nil
	case "connector":
		connector, err := s.UpdateConnector(ctx, canvasID, widgetID, req)
		if err != nil {
			return nil, err
		}
		return &Widget{ID: connector.ID, WidgetType: "connector"}, nil
	default:
		return nil, fmt.Errorf("UpdateWidget: unsupported widget_type: %s", widgetType)
	}
}

// DeleteWidget deletes a widget of the specified type by dispatching to the correct resource-specific method.
// widgetType must be provided (not inferred from the server).
func (s *Session) DeleteWidget(ctx context.Context, canvasID, widgetID, widgetType string) error {
	switch widgetType {
	case "note":
		return s.DeleteNote(ctx, canvasID, widgetID)
	case "anchor":
		return s.DeleteAnchor(ctx, canvasID, widgetID)
	case "image":
		return s.DeleteImage(ctx, canvasID, widgetID)
	case "pdf":
		return s.DeletePDF(ctx, canvasID, widgetID)
	case "video":
		return s.DeleteVideo(ctx, canvasID, widgetID)
	case "connector":
		return s.DeleteConnector(ctx, canvasID, widgetID)
	default:
		return fmt.Errorf("DeleteWidget: unsupported widget_type: %s", widgetType)
	}
}

// PatchParentID updates the parent ID of a widget (parenting).
func (s *Session) PatchParentID(ctx context.Context, canvasID, widgetID, parentID string) (*Widget, error) {
	var widget Widget
	path := fmt.Sprintf("canvases/%s/widgets/%s", canvasID, widgetID)
	req := map[string]interface{}{"parent_id": parentID}
	err := s.doRequest(ctx, "PATCH", path, req, &widget, nil, false)
	if err != nil {
		return nil, fmt.Errorf("PatchParentID: %w", err)
	}
	return &widget, nil
}

// WidgetMatch represents a widget match result across canvases.
type WidgetMatch struct {
	CanvasID string
	WidgetID string
	Widget   Widget
}

// WidgetsLister defines the interface for listing canvases and widgets.
type WidgetsLister interface {
	ListCanvases(ctx context.Context, filter *Filter) ([]Canvas, error)
	ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error)
}

// FindWidgetsAcrossCanvases searches all canvases for widgets matching the given query.
// The query supports exact, wildcard, and partial string matches (see Filter abstraction).
// Returns a slice of WidgetMatch with CanvasID, WidgetID, and the Widget itself.
func FindWidgetsAcrossCanvases(ctx context.Context, lister WidgetsLister, query map[string]interface{}) ([]WidgetMatch, error) {
	canvases, err := lister.ListCanvases(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("FindWidgetsAcrossCanvases: failed to list canvases: %w", err)
	}
	filter := &Filter{Criteria: query}
	var matches []WidgetMatch
	for _, canvas := range canvases {
		widgets, err := lister.ListWidgets(ctx, canvas.ID, filter)
		if err != nil {
			return nil, fmt.Errorf("FindWidgetsAcrossCanvases: failed to list widgets for canvas %s: %w", canvas.ID, err)
		}
		for _, w := range widgets {
			matches = append(matches, WidgetMatch{
				CanvasID: canvas.ID,
				WidgetID: w.ID,
				Widget:   w,
			})
		}
	}
	return matches, nil
}

// WidgetZone represents a container widget and its contained widgets, plus canvas and SharedCanvas context.
type WidgetZone struct {
	CanvasID       string
	SharedCanvasID string
	Container      Widget
	Contents       []Widget
}

// WidgetsContainId returns a WidgetZone: the source widget as Container, and all widgets fully contained within it as Contents.
// For all returned widgets, if ParentID matches the SharedCanvas ID, it is set to "".
func WidgetsContainId(ctx context.Context, s *Session, canvasID string, widgetID string, widget *Widget, tolerance float64) (WidgetZone, error) {
	var srcWidget Widget
	if widget != nil {
		srcWidget = *widget
	} else {
		if widgetID == "" {
			return WidgetZone{}, fmt.Errorf("WidgetsContainId: widgetID must be provided if widget is nil")
		}
		w, err := s.GetWidget(ctx, canvasID, widgetID)
		if err != nil {
			return WidgetZone{}, fmt.Errorf("WidgetsContainId: failed to fetch widget: %w", err)
		}
		srcWidget = *w
	}

	// Fetch all widgets on the same canvas
	widgets, err := s.ListWidgets(ctx, canvasID, nil)
	if err != nil {
		return WidgetZone{}, fmt.Errorf("WidgetsContainId: failed to list widgets: %w", err)
	}

	// Find SharedCanvas ID
	var sharedCanvasID string
	for _, w := range widgets {
		if w.WidgetType == "SharedCanvas" {
			sharedCanvasID = w.ID
			break
		}
	}

	srcRect := WidgetBoundingBox(srcWidget)
	// Expand bounding box by tolerance
	srcRect.X -= tolerance
	srcRect.Y -= tolerance
	srcRect.Width += 2 * tolerance
	srcRect.Height += 2 * tolerance

	var contained []Widget
	for _, w := range widgets {
		if w.ID == srcWidget.ID {
			// Debug log for skipping self
			fmt.Printf("[CONTAINS] Skipping self: %s\n", w.ID)
			continue // skip self
		}
		if w.WidgetType == "SharedCanvas" {
			continue // skip SharedCanvas in results
		}
		if WidgetContainsRect(srcRect, w) {
			// Normalize ParentID if it matches SharedCanvas
			if sharedCanvasID != "" && w.ParentID == sharedCanvasID {
				w.ParentID = ""
			}
			contained = append(contained, w)
		}
	}
	// Normalize ParentID for the container as well
	if sharedCanvasID != "" && srcWidget.ParentID == sharedCanvasID {
		srcWidget.ParentID = ""
	}
	return WidgetZone{CanvasID: canvasID, SharedCanvasID: sharedCanvasID, Container: srcWidget, Contents: contained}, nil
}

// WidgetContainsRect returns true if the given rectangle fully contains the widget's bounding box.
func WidgetContainsRect(rect Rectangle, w Widget) bool {
	return Contains(rect, WidgetBoundingBox(w))
}

// WidgetsTouchId returns a WidgetZone: the source widget as Container, and all widgets that touch it as Contents.
// For all returned widgets, if ParentID matches the SharedCanvas ID, it is set to "".
func WidgetsTouchId(ctx context.Context, s *Session, canvasID string, widgetID string, widget *Widget, tolerance float64) (WidgetZone, error) {
	var srcWidget Widget
	if widget != nil {
		srcWidget = *widget
	} else {
		if widgetID == "" {
			return WidgetZone{}, fmt.Errorf("WidgetsTouchId: widgetID must be provided if widget is nil")
		}
		w, err := s.GetWidget(ctx, canvasID, widgetID)
		if err != nil {
			return WidgetZone{}, fmt.Errorf("WidgetsTouchId: failed to fetch widget: %w", err)
		}
		srcWidget = *w
	}

	// Fetch all widgets on the same canvas
	widgets, err := s.ListWidgets(ctx, canvasID, nil)
	if err != nil {
		return WidgetZone{}, fmt.Errorf("WidgetsTouchId: failed to list widgets: %w", err)
	}

	// Find SharedCanvas ID
	var sharedCanvasID string
	for _, w := range widgets {
		if w.WidgetType == "SharedCanvas" {
			sharedCanvasID = w.ID
			break
		}
	}

	srcRect := WidgetBoundingBox(srcWidget)
	// Expand bounding box by tolerance
	srcRect.X -= tolerance
	srcRect.Y -= tolerance
	srcRect.Width += 2 * tolerance
	srcRect.Height += 2 * tolerance

	var touched []Widget
	for _, w := range widgets {
		if w.ID == srcWidget.ID {
			fmt.Printf("[TOUCHES] Skipping self: %s\n", w.ID)
			continue // skip self
		}
		if w.WidgetType == "SharedCanvas" {
			continue // skip SharedCanvas in results
		}
		if Touches(srcRect, WidgetBoundingBox(w)) {
			// Normalize ParentID if it matches SharedCanvas
			if sharedCanvasID != "" && w.ParentID == sharedCanvasID {
				w.ParentID = ""
			}
			touched = append(touched, w)
		}
	}
	// Normalize ParentID for the container as well
	if sharedCanvasID != "" && srcWidget.ParentID == sharedCanvasID {
		srcWidget.ParentID = ""
	}
	return WidgetZone{CanvasID: canvasID, SharedCanvasID: sharedCanvasID, Container: srcWidget, Contents: touched}, nil
}

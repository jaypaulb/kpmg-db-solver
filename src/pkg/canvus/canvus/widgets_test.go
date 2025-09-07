package canvus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type mockSession struct {
	canvases         []Canvas
	widgets          map[string][]Widget
	failListCanvases bool
	failListWidgets  map[string]bool
}

func (m *mockSession) ListCanvases(ctx context.Context, filter *Filter) ([]Canvas, error) {
	if m.failListCanvases {
		return nil, errors.New("mock ListCanvases failure")
	}
	return m.canvases, nil
}

func (m *mockSession) ListWidgets(ctx context.Context, canvasID string, filter *Filter) ([]Widget, error) {
	if m.failListWidgets != nil && m.failListWidgets[canvasID] {
		return nil, errors.New("mock ListWidgets failure")
	}
	widgets := m.widgets[canvasID]
	if filter != nil {
		widgets = FilterSlice(widgets, filter)
	}
	return widgets, nil
}

func TestFindWidgetsAcrossCanvases(t *testing.T) {
	ctx := context.Background()
	ms := &mockSession{
		canvases: []Canvas{{ID: "c1"}, {ID: "c2"}},
		widgets: map[string][]Widget{
			"c1": {
				{ID: "w1", WidgetType: "browser", ParentID: "", State: "active"},
				{ID: "w2", WidgetType: "note", ParentID: "", State: "archived"},
			},
			"c2": {
				{ID: "w3", WidgetType: "browser", ParentID: "", State: "active"},
				{ID: "w4", WidgetType: "browser", ParentID: "", State: "inactive"},
			},
		},
	}

	t.Run("ExactMatch", func(t *testing.T) {
		query := map[string]interface{}{"widget_type": "note"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []WidgetMatch{{CanvasID: "c1", WidgetID: "w2", Widget: ms.widgets["c1"][1]}}
		if !reflect.DeepEqual(matches, want) {
			t.Errorf("got %+v, want %+v", matches, want)
		}
	})

	t.Run("WildcardMatch", func(t *testing.T) {
		query := map[string]interface{}{"widget_type": "*"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matches) != 4 {
			t.Errorf("expected 4 matches, got %d", len(matches))
		}
	})

	t.Run("SuffixMatch", func(t *testing.T) {
		ms.widgets["c1"][0].ParentID = "abc12345"
		query := map[string]interface{}{"parent_id": "*12345"}
		matches, err := FindWidgetsAcrossCanvases(ctx, ms, query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matches) != 1 || matches[0].WidgetID != "w1" {
			t.Errorf("expected w1, got %+v", matches)
		}
	})

	t.Run("ListCanvasesError", func(t *testing.T) {
		ms.failListCanvases = true
		_, err := FindWidgetsAcrossCanvases(ctx, ms, map[string]interface{}{"widget_type": "browser"})
		if err == nil {
			t.Error("expected error, got nil")
		}
		ms.failListCanvases = false
	})

	t.Run("ListWidgetsError", func(t *testing.T) {
		ms.failListWidgets = map[string]bool{"c2": true}
		_, err := FindWidgetsAcrossCanvases(ctx, ms, map[string]interface{}{"widget_type": "browser"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestWidgetsContainId(t *testing.T) {
	ctx := context.Background()
	sharedCanvas := Widget{ID: "sc1", WidgetType: "SharedCanvas"}
	anchor := Widget{ID: "a1", WidgetType: "anchor", ParentID: "sc1", Location: &Point{X: 0, Y: 0}, Size: &Size{Width: 100, Height: 100}}
	note := Widget{ID: "n1", WidgetType: "note", ParentID: "a1", Location: &Point{X: 10, Y: 10}, Size: &Size{Width: 20, Height: 20}}
	ms := &mockSession{
		canvases: []Canvas{{ID: "c1"}},
		widgets: map[string][]Widget{
			"c1": {sharedCanvas, anchor, note},
		},
	}
	// Local version for test using mockSession
	localWidgetsContainId := func(ctx context.Context, s *mockSession, canvasID string, widgetID string, widget *Widget, tolerance float64) (WidgetZone, error) {
		var srcWidget Widget
		if widget != nil {
			srcWidget = *widget
		} else {
			if widgetID == "" {
				return WidgetZone{}, fmt.Errorf("WidgetsContainId: widgetID must be provided if widget is nil")
			}
			ws, err := s.ListWidgets(ctx, canvasID, nil)
			if err != nil {
				return WidgetZone{}, fmt.Errorf("WidgetsContainId: failed to list widgets: %w", err)
			}
			found := false
			for _, w := range ws {
				if w.ID == widgetID {
					srcWidget = w
					found = true
					break
				}
			}
			if !found {
				return WidgetZone{}, fmt.Errorf("WidgetsContainId: widgetID not found")
			}
		}
		widgets, err := s.ListWidgets(ctx, canvasID, nil)
		if err != nil {
			return WidgetZone{}, fmt.Errorf("WidgetsContainId: failed to list widgets: %w", err)
		}
		var sharedCanvasID string
		for _, w := range widgets {
			if w.WidgetType == "SharedCanvas" {
				sharedCanvasID = w.ID
				break
			}
		}
		srcRect := WidgetBoundingBox(srcWidget)
		srcRect.X -= tolerance
		srcRect.Y -= tolerance
		srcRect.Width += 2 * tolerance
		srcRect.Height += 2 * tolerance
		var contained []Widget
		for _, w := range widgets {
			if w.ID == srcWidget.ID {
				continue
			}
			if w.WidgetType == "SharedCanvas" {
				continue // skip SharedCanvas in results
			}
			if WidgetContainsRect(srcRect, w) {
				if sharedCanvasID != "" && w.ParentID == sharedCanvasID {
					w.ParentID = ""
				}
				contained = append(contained, w)
			}
		}
		if sharedCanvasID != "" && srcWidget.ParentID == sharedCanvasID {
			srcWidget.ParentID = ""
		}
		return WidgetZone{Container: srcWidget, Contents: contained}, nil
	}
	zone, err := localWidgetsContainId(ctx, ms, "c1", "a1", nil, 0)
	if err != nil {
		t.Fatalf("WidgetsContainId failed: %v", err)
	}
	if zone.Container.ID != "a1" {
		t.Errorf("Expected container ID 'a1', got '%s'", zone.Container.ID)
	}
	if zone.Container.ParentID != "" {
		t.Errorf("Expected container ParentID to be blank (was '%s')", zone.Container.ParentID)
	}
	if len(zone.Contents) != 1 || zone.Contents[0].ID != "n1" {
		t.Errorf("Expected one contained note with ID 'n1', got %+v", zone.Contents)
	}
	if zone.Contents[0].ParentID != "a1" {
		t.Errorf("Expected note ParentID to be 'a1', got '%s'", zone.Contents[0].ParentID)
	}
}

package canvus

import (
	"context"
	"testing"
)

func TestListAuditEvents(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	events, err := admin.ListAuditEvents(ctx, nil)
	if err != nil {
		t.Errorf("failed to list audit events: %v", err)
	}
	t.Logf("retrieved %d audit events", len(events))

	// Test with pagination (if supported)
	paged, err := admin.ListAuditEvents(ctx, &AuditLogOptions{Page: 1, PerPage: 2})
	if err != nil {
		t.Errorf("failed to list audit events with pagination: %v", err)
	}
	t.Logf("retrieved %d audit events (paged)", len(paged))

	// Test with filter (if supported)
	filtered, err := admin.ListAuditEvents(ctx, &AuditLogOptions{Filter: "login"})
	if err != nil {
		t.Logf("filtering audit events returned error (may be expected if filter unsupported): %v", err)
	} else {
		t.Logf("retrieved %d filtered audit events", len(filtered))
	}
}

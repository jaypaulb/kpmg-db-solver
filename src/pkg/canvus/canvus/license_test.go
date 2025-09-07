package canvus

import (
	"context"
	"testing"
)

func TestGetLicenseInfo(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	info, err := admin.GetLicenseInfo(ctx)
	if err != nil {
		t.Errorf("failed to get license info: %v", err)
	}
	if info == nil {
		t.Errorf("expected license info, got nil")
	}
	// Optionally check some fields if known
}

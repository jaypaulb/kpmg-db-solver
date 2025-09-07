package canvus

import (
	"context"
	"testing"
	"time"
)

func TestServerConfigLifecycle(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Get current config (for later restore)
	orig, err := admin.GetServerConfig(ctx)
	if err != nil {
		t.Fatalf("failed to get server config: %v", err)
	}

	// Always restore original config, even if test fails
	defer func() {
		if _, err := admin.UpdateServerConfig(ctx, ServerConfig{ServerName: orig.ServerName}); err != nil {
			t.Errorf("failed to restore original serverName: %v", err)
		}
	}()

	// Update a value (e.g., server_name)
	newName := "TestServer_" + time.Now().Format("150405")
	update := ServerConfig{
		ServerName: newName,
	}
	updated, err := admin.UpdateServerConfig(ctx, update)
	if err != nil {
		t.Errorf("failed to update server config: %v", err)
	} else if updated.ServerName != newName {
		t.Errorf("expected server_name %q, got %q", newName, updated.ServerName)
	}

	// Send test email (should not error, but may not be supported in all environments)
	err = admin.SendTestEmail(ctx)
	if err != nil {
		t.Logf("SendTestEmail failed (may be expected if SMTP not configured): %v", err)
	}
}

func TestServerConfigInvalidCases(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Try to update with invalid data (e.g., invalid field)
	update := ServerConfig{
		ServerName: "",
	}
	_, err = admin.UpdateServerConfig(ctx, update)
	// This may or may not error depending on server rules, so just log
	if err != nil {
		t.Logf("UpdateServerConfig with empty name returned error (may be expected): %v", err)
	}
}

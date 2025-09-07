package canvus

import (
	"context"
	"testing"
)

func TestGetServerInfo(t *testing.T) {
	ctx := context.Background()
	client, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	info, err := client.GetServerInfo(ctx)
	if err != nil {
		t.Errorf("failed to get server info: %v", err)
	}
	if info == nil {
		t.Errorf("expected server info, got nil")
	}
	t.Logf("ServerInfo: %+v", info)
}

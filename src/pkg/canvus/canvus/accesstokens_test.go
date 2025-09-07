package canvus

import (
	"context"
	"testing"
	"time"
)

func TestAccessTokenLifecycle(t *testing.T) {
	ctx := context.Background()
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)

	// Create a test user to own the token
	email := "tokenuser_" + time.Now().Format("20060102150405") + "@example.com"
	name := "tokenuser_" + time.Now().Format("150405")
	password := "TestPassword123!"
	user, err := admin.CreateUser(ctx, CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer func() { _ = admin.DeleteUser(ctx, user.ID) }()

	// Create access token
	desc := "test token " + time.Now().Format("150405")
	token, err := admin.CreateAccessToken(ctx, user.ID, CreateAccessTokenRequest{
		Description: desc,
	})
	if err != nil {
		t.Fatalf("failed to create access token: %v", err)
	}
	defer func() { _ = admin.DeleteAccessToken(ctx, user.ID, token.ID) }()

	// Retrieve access token
	got, err := admin.GetAccessToken(ctx, user.ID, token.ID)
	if err != nil {
		t.Errorf("failed to get access token: %v", err)
	}
	if got.ID != token.ID {
		t.Errorf("expected token ID %q, got %q", token.ID, got.ID)
	}
	if got.Description != desc {
		t.Errorf("expected description %q, got %q", desc, got.Description)
	}

	// Delete access token (already deferred)
}

func TestAccessTokenInvalid(t *testing.T) {
	ctx := context.Background()
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)
	// Try to get a non-existent token
	_, err = admin.GetAccessToken(ctx, 999999, "nonexistent-token-id")
	if err == nil {
		t.Fatalf("expected error for invalid token id, got nil")
	}
}

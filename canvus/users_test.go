package canvus

import (
	"context"
	"testing"
	"time"
)

func TestUserLifecycle(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	email := "testuser_" + time.Now().Format("20060102150405") + "@example.com"
	name := "testuser_" + time.Now().Format("150405")
	password := "TestPassword123!"

	// Create user
	user, err := admin.CreateUser(ctx, CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer func() {
		_ = admin.DeleteUser(ctx, user.ID)
	}()

	// Retrieve user
	got, err := admin.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("failed to get user: %v", err)
	}
	if got.Email != email {
		t.Errorf("expected email %q, got %q", email, got.Email)
	}
	if got.Name != name {
		t.Errorf("expected name %q, got %q", name, got.Name)
	}

	// Login as user
	testClient := NewSession(admin.BaseURL)
	err = testClient.Login(ctx, email, password)
	if err != nil {
		t.Errorf("failed to login as user: %v", err)
	}

	// Logout
	err = testClient.Logout(ctx)
	if err != nil && err.Error() != "EOF" {
		t.Errorf("failed to logout: %v", err)
	} else if err != nil && err.Error() == "EOF" {
		t.Logf("logout returned EOF, which is expected if the API returns no response body")
	}

	// Delete user (already deferred)
}

func TestUserLoginInvalid(t *testing.T) {
	ctx := context.Background()
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	testClient := NewSession(ts.APIBaseURL)
	err = testClient.Login(ctx, "nonexistent@example.com", "wrongpassword")
	if err == nil {
		t.Fatalf("expected error for invalid login, got nil")
	}
}

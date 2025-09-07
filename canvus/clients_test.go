package canvus

import (
	"context"
	"testing"
	"time"
)

func TestTestClientLifecycle(t *testing.T) {
	ctx := context.Background()
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)

	testEmail := "testclient_" + time.Now().Format("20060102150405") + "@example.com"
	testName := "testclient_" + time.Now().Format("150405")
	testPassword := "TestPassword123!"

	tc, err := NewTestClient(ctx, admin, ts.APIBaseURL, testEmail, testName, testPassword)
	if err != nil {
		t.Fatalf("failed to create TestClient: %v", err)
	}
	defer func() { _ = tc.Cleanup(ctx) }()

	// Perform an action as the test client (e.g., get self user info)
	user, err := tc.Session.GetUser(ctx, tc.userID)
	if err != nil {
		t.Errorf("TestClient could not get user: %v", err)
	}
	if user.Email != testEmail {
		t.Errorf("expected email %q, got %q", testEmail, user.Email)
	}
	if user.Name != testName {
		t.Errorf("expected name %q, got %q", testName, user.Name)
	}
}

func TestUserClientLifecycle(t *testing.T) {
	ctx := context.Background()
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	admin := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)

	email := "userclient_" + time.Now().Format("20060102150405") + "@example.com"
	name := "userclient_" + time.Now().Format("150405")
	password := "TestPassword123!"

	// Create user for this test
	user, err := admin.CreateUser(ctx, CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer func() { _ = admin.DeleteUser(ctx, user.ID) }()

	uc, err := NewUserClient(ctx, ts.APIBaseURL, email, password)
	if err != nil {
		t.Fatalf("failed to create UserClient: %v", err)
	}
	defer func() { _ = uc.Cleanup(ctx) }()

	// Perform an action as the user client
	got, err := uc.Session.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("UserClient could not get user: %v", err)
	}
	if got.Email != email {
		t.Errorf("expected email %q, got %q", email, got.Email)
	}
	if got.Name != name {
		t.Errorf("expected name %q, got %q", name, got.Name)
	}
}

func TestNewClientFromConfig(t *testing.T) {
	ts, err := loadTestSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}
	client := NewSessionFromConfig(ts.APIBaseURL, ts.APIKey)
	if client == nil {
		t.Fatalf("NewClientFromConfig returned nil client")
	}
	// Optionally, test a simple action (e.g., list users)
	ctx := context.Background()
	_, err = client.ListUsers(ctx)
	if err != nil {
		t.Errorf("ListUsers failed: %v", err)
	}
}

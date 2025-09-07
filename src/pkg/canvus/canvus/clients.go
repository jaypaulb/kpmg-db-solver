package canvus

import (
	"context"
	"fmt"
)

// ClientInfo represents a client in the Canvus system.
type ClientInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	// Add other fields as needed
}

// CreateClientRequest is the payload for creating a new client.
type CreateClientRequest struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	// Add other fields as needed
}

// UpdateClientRequest is the payload for updating an existing client.
type UpdateClientRequest struct {
	Name *string `json:"name,omitempty"`
	// Add other fields as needed
}

// ListClients retrieves all clients from the Canvus API.
func (c *Session) ListClients(ctx context.Context) ([]ClientInfo, error) {
	var clients []ClientInfo
	err := c.doRequest(ctx, "GET", "clients", nil, &clients, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListClients: %w", err)
	}
	return clients, nil
}

// GetClient retrieves a client by ID from the Canvus API.
func (c *Session) GetClient(ctx context.Context, id string) (*ClientInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("GetClient: id is required")
	}
	var client ClientInfo
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "GET", endpoint, nil, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetClient: %w", err)
	}
	return &client, nil
}

// CreateClient creates a new client in the Canvus API.
func (c *Session) CreateClient(ctx context.Context, req CreateClientRequest) (*ClientInfo, error) {
	var client ClientInfo
	err := c.doRequest(ctx, "POST", "clients", req, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateClient: %w", err)
	}
	return &client, nil
}

// UpdateClient updates an existing client by ID in the Canvus API.
func (c *Session) UpdateClient(ctx context.Context, id string, req UpdateClientRequest) (*ClientInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("UpdateClient: id is required")
	}
	var client ClientInfo
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "PATCH", endpoint, req, &client, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateClient: %w", err)
	}
	return &client, nil
}

// DeleteClient deletes a client by ID in the Canvus API.
func (c *Session) DeleteClient(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("DeleteClient: id is required")
	}
	endpoint := fmt.Sprintf("clients/%s", id)
	err := c.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteClient: %w", err)
	}
	return nil
}

// TestClient wraps a Session and manages a temporary test user and token.
// TestClient is used in tests to create, login, and clean up a dedicated test user.
type TestClient struct {
	Session     *Session
	userID      int64
	email       string
	password    string
	cleanupFunc func(context.Context) error
}

// UserClient wraps a Session and manages a temporary token for an existing user.
// UserClient is used in tests to login as an existing user and clean up the session.
type UserClient struct {
	Session     *Session
	cleanupFunc func(context.Context) error
}

// NewTestClient creates a new test user, logs in as that user, and returns a TestClient.
// The test user and token are deleted on Cleanup.
func NewTestClient(ctx context.Context, adminSession *Session, baseURL, testEmail, testUsername, testPassword string) (*TestClient, error) {
	// 1. Create user
	user, err := adminSession.CreateUser(ctx, CreateUserRequest{
		Name:     testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	if err != nil {
		return nil, err
	}
	// 2. Login as new user
	testSession := NewSession(baseURL)
	err = testSession.Login(ctx, testEmail, testPassword)
	if err != nil {
		// Cleanup user if login fails
		_ = adminSession.DeleteUser(ctx, user.ID)
		return nil, err
	}
	cleanup := func(ctx context.Context) error {
		_ = testSession.Logout(ctx)
		return adminSession.DeleteUser(ctx, user.ID)
	}
	return &TestClient{
		Session:     testSession,
		userID:      user.ID,
		email:       testEmail,
		password:    testPassword,
		cleanupFunc: cleanup,
	}, nil
}

// Cleanup logs out and deletes the test user.
func (tc *TestClient) Cleanup(ctx context.Context) error {
	if tc.cleanupFunc != nil {
		err := tc.cleanupFunc(ctx)
		if err != nil {
			fmt.Printf("[CLEANUP ERROR] TestClient cleanup failed for %s: %v\n", tc.email, err)
		}
		return err
	}
	return nil
}

// NewUserClient logs in as an existing user and returns a UserClient with a temporary token.
// The token is invalidated on Cleanup.
func NewUserClient(ctx context.Context, baseURL, email, password string) (*UserClient, error) {
	session := NewSession(baseURL)
	err := session.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	cleanup := func(ctx context.Context) error {
		return session.Logout(ctx)
	}
	return &UserClient{
		Session:     session,
		cleanupFunc: cleanup,
	}, nil
}

// Cleanup logs out and invalidates the token.
func (uc *UserClient) Cleanup(ctx context.Context) error {
	if uc.cleanupFunc != nil {
		err := uc.cleanupFunc(ctx)
		if err != nil {
			fmt.Printf("[CLEANUP ERROR] UserClient cleanup failed: %v\n", err)
		}
		return err
	}
	return nil
}

// NewSessionFromConfig creates a Session using credentials from a config/settings file.
// This is the standard persistent client; no automatic cleanup is performed.
func NewSessionFromConfig(baseURL, apiKey string) *Session {
	return NewSession(baseURL, WithAPIKey(apiKey))
}

package canvus

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

// User represents a user in the Canvus system.
// User contains basic identity and contact information.
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
	Approved  bool   `json:"approved"`
	Blocked   bool   `json:"blocked"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login"`
	State     string `json:"state"`
	// Add other fields as needed
}

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Admin    *bool  `json:"admin,omitempty"`
	Approved *bool  `json:"approved,omitempty"`
	Blocked  *bool  `json:"blocked,omitempty"`
	// Add other fields as needed
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	Admin    *bool   `json:"admin,omitempty"`
	Approved *bool   `json:"approved,omitempty"`
	Blocked  *bool   `json:"blocked,omitempty"`
	// Add other fields as needed
}

// ListUsers retrieves all users from the Canvus API.
func (s *Session) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := s.doRequest(ctx, "GET", "users", nil, &users, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: %w", err)
	}
	return users, nil
}

// GetUser retrieves a user by ID from the Canvus API.
func (s *Session) GetUser(ctx context.Context, id int64) (*User, error) {
	if id == 0 {
		return nil, fmt.Errorf("GetUser: id is required")
	}
	var user User
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "GET", endpoint, nil, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user in the Canvus API.
func (s *Session) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	marker := time.Now().UnixNano()
	fmt.Printf("[DEBUG] CreateUser called: marker=%d\n%s\n", marker, debug.Stack())
	ctx = context.WithValue(ctx, "createUserMarker", marker)
	var user User
	err := s.doRequest(ctx, "POST", "users", req, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user by ID in the Canvus API.
func (s *Session) UpdateUser(ctx context.Context, id int64, req UpdateUserRequest) (*User, error) {
	if id == 0 {
		return nil, fmt.Errorf("UpdateUser: id is required")
	}
	var user User
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "PATCH", endpoint, req, &user, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser: %w", err)
	}
	return &user, nil
}

// DeleteUser deletes a user by ID in the Canvus API.
func (s *Session) DeleteUser(ctx context.Context, id int64) error {
	if id == 0 {
		return fmt.Errorf("DeleteUser: id is required")
	}
	endpoint := fmt.Sprintf("users/%d", id)
	err := s.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}
	return nil
}

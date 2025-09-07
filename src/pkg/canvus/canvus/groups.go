package canvus

import (
	"context"
	"fmt"
)

// GroupMember represents a user in a group.
type GroupMember struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
	Approved  bool   `json:"approved"`
	Blocked   bool   `json:"blocked"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login"`
	State     string `json:"state"`
}

// CreateGroupRequest is the payload for creating a group.
type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// AddUserToGroupRequest is the payload for adding a user to a group.
type AddUserToGroupRequest struct {
	ID int `json:"id"`
}

// Group represents a user group in Canvus.
type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ListGroups retrieves all groups from the Canvus API.
func (s *Session) ListGroups(ctx context.Context) ([]Group, error) {
	var groups []Group
	err := s.doRequest(ctx, "GET", "groups", nil, &groups, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListGroups: %w", err)
	}
	return groups, nil
}

// GetGroup retrieves a single group by ID from the Canvus API.
func (s *Session) GetGroup(ctx context.Context, id int) (*Group, error) {
	var group Group
	path := fmt.Sprintf("groups/%d", id)
	err := s.doRequest(ctx, "GET", path, nil, &group, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetGroup: %w", err)
	}
	return &group, nil
}

// CreateGroup creates a new group in the Canvus API.
func (s *Session) CreateGroup(ctx context.Context, req CreateGroupRequest) (*Group, error) {
	var group Group
	err := s.doRequest(ctx, "POST", "groups", req, &group, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateGroup: %w", err)
	}
	return &group, nil
}

// DeleteGroup deletes a group by ID in the Canvus API.
func (s *Session) DeleteGroup(ctx context.Context, id int) error {
	path := fmt.Sprintf("groups/%d", id)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

// AddUserToGroup adds a user to a group.
func (s *Session) AddUserToGroup(ctx context.Context, groupID int, userID int) error {
	path := fmt.Sprintf("groups/%d/members", groupID)
	body := AddUserToGroupRequest{ID: userID}
	return s.doRequest(ctx, "POST", path, body, nil, nil, false)
}

// ListGroupMembers lists all users in a group.
func (s *Session) ListGroupMembers(ctx context.Context, groupID int) ([]GroupMember, error) {
	path := fmt.Sprintf("groups/%d/members", groupID)
	var members []GroupMember
	err := s.doRequest(ctx, "GET", path, nil, &members, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListGroupMembers: %w", err)
	}
	return members, nil
}

// RemoveUserFromGroup removes a user from a group.
func (s *Session) RemoveUserFromGroup(ctx context.Context, groupID int, userID int) error {
	path := fmt.Sprintf("groups/%d/members/%d", groupID, userID)
	return s.doRequest(ctx, "DELETE", path, nil, nil, nil, false)
}

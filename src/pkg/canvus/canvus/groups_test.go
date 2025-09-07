package canvus

import (
	"context"
	"testing"
	"time"
)

func TestGroupLifecycle(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Create a unique group
	groupName := "testgroup_" + time.Now().Format("20060102150405")
	desc := "Test group created at " + time.Now().Format("15:04:05")
	group, err := admin.CreateGroup(ctx, CreateGroupRequest{
		Name:        groupName,
		Description: desc,
	})
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}
	// Ensure group is deleted at the end
	defer func() { _ = admin.DeleteGroup(ctx, group.ID) }()

	// List groups and check the new group is present
	groups, err := admin.ListGroups(ctx)
	if err != nil {
		t.Errorf("failed to list groups: %v", err)
	}
	found := false
	for _, g := range groups {
		if g.ID == group.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created group not found in list")
	}

	// Get the group by ID
	got, err := admin.GetGroup(ctx, group.ID)
	if err != nil {
		t.Errorf("failed to get group: %v", err)
	}
	if got.Name != groupName {
		t.Errorf("expected group name %q, got %q", groupName, got.Name)
	}
	if got.Description != desc {
		t.Errorf("expected description %q, got %q", desc, got.Description)
	}

	// Create a test user to add to the group
	testEmail := "groupuser_" + time.Now().Format("20060102150405") + "@example.com"
	testName := "groupuser_" + time.Now().Format("150405")
	testPassword := "TestPassword123!"
	user, err := admin.CreateUser(ctx, CreateUserRequest{
		Name:     testName,
		Email:    testEmail,
		Password: testPassword,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer func() { _ = admin.DeleteUser(ctx, user.ID) }()

	// Add user to group
	err = admin.AddUserToGroup(ctx, group.ID, int(user.ID))
	if err != nil {
		t.Errorf("failed to add user to group: %v", err)
	}

	// List group members and check the user is present
	members, err := admin.ListGroupMembers(ctx, group.ID)
	if err != nil {
		t.Errorf("failed to list group members: %v", err)
	}
	userFound := false
	for _, m := range members {
		if m.ID == int(user.ID) {
			userFound = true
			break
		}
	}
	if !userFound {
		t.Errorf("added user not found in group members")
	}

	// Remove user from group
	err = admin.RemoveUserFromGroup(ctx, group.ID, int(user.ID))
	if err != nil {
		t.Errorf("failed to remove user from group: %v", err)
	}

	// List group members again to ensure user is removed
	members, err = admin.ListGroupMembers(ctx, group.ID)
	if err != nil {
		t.Errorf("failed to list group members after removal: %v", err)
	}
	for _, m := range members {
		if m.ID == int(user.ID) {
			t.Errorf("user still present in group after removal")
		}
	}

	// Delete group (already deferred)
}

func TestGroupInvalidCases(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Get non-existent group
	_, err = admin.GetGroup(ctx, 99999999)
	if err == nil {
		t.Errorf("expected error for non-existent group, got nil")
	}

	// Delete non-existent group
	err = admin.DeleteGroup(ctx, 99999999)
	if err == nil {
		t.Errorf("expected error for deleting non-existent group, got nil")
	}

	// Add user to non-existent group
	err = admin.AddUserToGroup(ctx, 99999999, 99999999)
	if err == nil {
		t.Errorf("expected error for adding user to non-existent group, got nil")
	}

	// Remove user from non-existent group
	err = admin.RemoveUserFromGroup(ctx, 99999999, 99999999)
	if err == nil {
		t.Errorf("expected error for removing user from non-existent group, got nil")
	}

	// List members of non-existent group
	_, err = admin.ListGroupMembers(ctx, 99999999)
	if err == nil {
		t.Errorf("expected error for listing members of non-existent group, got nil")
	}
}

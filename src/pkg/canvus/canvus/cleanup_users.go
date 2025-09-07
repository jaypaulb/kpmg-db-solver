// cleanup_users.go
// Standalone utility: deletes all test users and test folders (names starting with testuser_ or testfolder_)
// WARNING: This file is NOT a test and should NOT be run as part of normal test suites.
// Run manually only when you want to clean up test artifacts.
package canvus

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		fmt.Printf("Failed to load test settings: %v\n", err)
		os.Exit(1)
	}

	// Clean up test users
	users, err := admin.ListUsers(ctx)
	if err != nil {
		fmt.Printf("Failed to list users: %v\n", err)
		os.Exit(1)
	}
	for _, user := range users {
		if strings.HasPrefix(user.Email, "testuser_") || strings.HasPrefix(user.Name, "testuser_") {
			fmt.Printf("Deleting test user: %s (%s)\n", user.Name, user.Email)
			if err := admin.DeleteUser(ctx, user.ID); err != nil {
				fmt.Printf("  [ERROR] Failed to delete user %s: %v\n", user.Email, err)
			}
		}
	}

	// Clean up test folders
	folders, err := admin.ListFolders(ctx)
	if err != nil {
		fmt.Printf("Failed to list folders: %v\n", err)
		os.Exit(1)
	}
	for _, folder := range folders {
		if strings.HasPrefix(folder.Name, "testfolder_") {
			fmt.Printf("Deleting test folder: %s\n", folder.Name)
			if err := admin.DeleteFolder(ctx, folder.ID); err != nil {
				fmt.Printf("  [ERROR] Failed to delete folder %s: %v\n", folder.Name, err)
			}
		}
	}

	fmt.Println("Cleanup complete.")
}

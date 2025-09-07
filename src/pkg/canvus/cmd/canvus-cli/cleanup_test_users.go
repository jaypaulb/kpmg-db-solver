package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"canvus-go-api/canvus" // Updated import path
)

// Config represents the structure of the settings.json file
type Config struct {
	APIBaseURL string `json:"api_base_url"`
	APIKey     string `json:"api_key"`
}

func main() {
	// Read configuration from settings.json
	configFile, err := os.ReadFile("settings.json")
	if err != nil {
		log.Fatalf("Error reading settings.json: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling settings.json: %v", err)
	}

	apiBaseURL := config.APIBaseURL
	apiKey := config.APIKey

	if apiBaseURL == "" || apiKey == "" {
		log.Fatal("API Base URL or API Key not found in settings.json")
	}

	// Create a new Canvus session
	session := canvus.NewSession(apiBaseURL, canvus.WithAPIKey(apiKey))
	ctx := context.Background()

	log.Printf("Connecting to Canvus API at %s", apiBaseURL)

	// List all users
	users, err := session.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}

	log.Printf("Found %d users.", len(users))

	deletedCount := 0
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Email), "@example.com") {
			log.Printf("Deleting user: %s (ID: %d)", user.Email, user.ID)
			err := session.DeleteUser(ctx, user.ID)
			if err != nil {
				log.Printf("Error deleting user %s (ID: %d): %v", user.Email, user.ID, err)
			} else {
				deletedCount++
			}
		}
	}

	log.Printf("Cleanup complete. Deleted %d users with '@example.com' in their email.", deletedCount)
}

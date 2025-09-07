// test_helpers.go
// Shared helpers for test setup and admin session loading.
package canvus

import (
	"encoding/json"
	"os"
)

// testSettings holds configuration for test and admin sessions.
type testSettings struct {
	APIBaseURL string `json:"api_base_url"`
	APIKey     string `json:"api_key"`
	TestUser   struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"test_user"`
}

// loadTestSettings loads test settings from ../settings.json.
func loadTestSettings() (*testSettings, error) {
	f, err := os.Open("../settings.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s testSettings
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// getTestAdminClientFromSettings returns an admin Session and test settings from config.
func getTestAdminClientFromSettings() (*Session, *testSettings, error) {
	ts, err := loadTestSettings()
	if err != nil {
		return nil, nil, err
	}
	return NewSessionFromConfig(ts.APIBaseURL, ts.APIKey), ts, nil
}

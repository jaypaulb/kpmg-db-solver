package canvus

import (
	"context"
	"fmt"
)

// AccessToken represents an API access token for a user.
type AccessToken struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	PlainToken  string `json:"plain_token,omitempty"`
}

// CreateAccessTokenRequest is the payload for creating a new access token.
type CreateAccessTokenRequest struct {
	Description string `json:"description"`
}

// ListAccessTokens retrieves all access tokens for a user from the Canvus API.
func (s *Session) ListAccessTokens(ctx context.Context, userID int64) ([]AccessToken, error) {
	var tokens []AccessToken
	endpoint := fmt.Sprintf("users/%d/access-tokens", userID)
	err := s.doRequest(ctx, "GET", endpoint, nil, &tokens, nil, false)
	if err != nil {
		return nil, fmt.Errorf("ListAccessTokens: %w", err)
	}
	return tokens, nil
}

// GetAccessToken retrieves an access token by ID for a user from the Canvus API.
func (s *Session) GetAccessToken(ctx context.Context, userID int64, tokenID string) (*AccessToken, error) {
	if tokenID == "" {
		return nil, fmt.Errorf("GetAccessToken: tokenID is required")
	}
	var token AccessToken
	endpoint := fmt.Sprintf("users/%d/access-tokens/%s", userID, tokenID)
	err := s.doRequest(ctx, "GET", endpoint, nil, &token, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetAccessToken: %w", err)
	}
	return &token, nil
}

// CreateAccessToken creates a new access token for a user in the Canvus API.
func (s *Session) CreateAccessToken(ctx context.Context, userID int64, req CreateAccessTokenRequest) (*AccessToken, error) {
	var token AccessToken
	endpoint := fmt.Sprintf("users/%d/access-tokens", userID)
	err := s.doRequest(ctx, "POST", endpoint, req, &token, nil, false)
	if err != nil {
		return nil, fmt.Errorf("CreateAccessToken: %w", err)
	}
	return &token, nil
}

// DeleteAccessToken deletes an access token by ID for a user in the Canvus API.
func (s *Session) DeleteAccessToken(ctx context.Context, userID int64, tokenID string) error {
	if tokenID == "" {
		return fmt.Errorf("DeleteAccessToken: tokenID is required")
	}
	endpoint := fmt.Sprintf("users/%d/access-tokens/%s", userID, tokenID)
	err := s.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, false)
	if err != nil {
		return fmt.Errorf("DeleteAccessToken: %w", err)
	}
	return nil
}

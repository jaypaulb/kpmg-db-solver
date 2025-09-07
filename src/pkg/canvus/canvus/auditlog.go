package canvus

import (
	"context"
	"encoding/json"
	"fmt"
)

// AuditEvent represents an audit log event in the Canvus system.
type AuditEvent struct {
	ID        json.Number `json:"id"`
	Timestamp string      `json:"timestamp,omitempty"`
	UserID    int64       `json:"user_id,omitempty"`
	Action    string      `json:"action,omitempty"`
	Resource  string      `json:"resource,omitempty"`
	Details   string      `json:"details,omitempty"`
	// Add other fields as needed based on the API response
}

// ListAuditEvents retrieves audit log events with optional pagination and filtering from the Canvus API.
func (s *Session) ListAuditEvents(ctx context.Context, opts *AuditLogOptions) ([]AuditEvent, error) {
	var events []AuditEvent
	query := map[string]string{}
	if opts != nil {
		if opts.PerPage > 0 {
			query["per_page"] = fmt.Sprintf("%d", opts.PerPage)
		}
	}
	err := s.doRequest(ctx, "GET", "audit-log", nil, &events, query, false)
	if err != nil {
		return nil, fmt.Errorf("ListAuditEvents: %w", err)
	}
	return events, nil
}

// Package canvus provides option types for SDK methods (pagination, filtering, etc.).
package canvus

// ListOptions specifies options for list endpoints (pagination, filtering, etc.).
type ListOptions struct {
	Limit  int    // Maximum number of items to return
	Offset int    // Offset for pagination
	Filter string // Optional filter string
}

// GetOptions specifies options for get endpoints (e.g., subscribe to updates).
type GetOptions struct {
	Subscribe bool // Whether to subscribe to updates (if supported)
}

// SubscribeOptions specifies options for streaming/subscription endpoints.
type SubscribeOptions struct {
	Annotations bool // Whether to include annotations
	// Add more fields as needed
}

// AuditLogOptions specifies options for querying the audit log.
type AuditLogOptions struct {
	Page    int    // Page number
	PerPage int    // Items per page
	Filter  string // Filter string
}

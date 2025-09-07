# Canvus Go SDK

A modern, idiomatic Go SDK for the Canvus API. Provides full access to all Canvus endpoints, strong typing, authentication, and developer-friendly abstractions.

## Features
- Full API coverage: system, canvas, workspace, widget, and asset endpoints
- Strongly typed request/response models
- Modular, idiomatic Go design
- Authentication (API key, login, token refresh)
- Context support for all requests
- Centralized error handling
- **Centralized response validation & retry logic:** All API calls are validated against the request, with robust error handling and automatic retries for transient failures.
- **Geometry utilities:** Utilities for determining widget containment and overlap (e.g., `WidgetsContainId`, `WidgetsTouchId`, `contains`, `touches`).
- Pagination and streaming helpers
- Comprehensive tests and examples

## Installation

```
go get github.com/jaypaulb/Canvus-Go-API/canvus
```

## Usage

```go
import (
    "context"
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Create a new session using an API key
    session := canvus.NewSession("https://your-canvus-server/api/v1", canvus.WithAPIKey("YOUR_API_KEY"))
    ctx := context.Background()

    // List canvases
    canvases, err := session.ListCanvases(ctx)
    if err != nil {
        panic(err)
    }
    for _, c := range canvases {
        println(c.ID, c.Name)
    }
}
```

## Authentication

The SDK supports three authentication flows:
- **API Key:** Use a static access token (recommended for most use cases)
- **Username/Password:** Obtain a temporary token via login
- **Token Refresh:** Prolong token lifetime via `/users/login`

Example (login):
```go
session := canvus.NewSession("https://your-canvus-server/api/v1")
err := session.Login(ctx, "user@example.com", "password")
if err != nil {
    // handle error
}
```

## Examples

- List folders:
  ```go
  folders, err := session.ListFolders(ctx)
  // ...
  ```
- Create a canvas:
  ```go
  canvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{Name: "My Canvas", FolderID: "..."})
  // ...
  ```
- Error handling:
  ```go
  if err != nil {
      if apiErr, ok := err.(*canvus.APIError); ok {
          fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
      } else {
          fmt.Println("Other error:", err)
      }
  }
  ```

## Testing

- All endpoints are covered by integration tests.
- To run tests, configure `settings.json` with your server URL and credentials, then run:
  ```
  go test ./canvus/...
  ```

## Contributing

Contributions are welcome! Please see `CONTRIBUTING.md` for guidelines.

## License

MIT License. See LICENSE file.

## Filtering (Client-Side)

You can filter canvases and widgets client-side using the `Filter` abstraction. This allows for flexible, in-memory filtering after fetching results from the API.

### Example: Filter Canvases by Name

```go
filter := &canvus.Filter{Criteria: map[string]interface{}{"name": "My Canvas"}}
canvases, err := session.ListCanvases(ctx, filter)
if err != nil {
    // handle error
}
for _, c := range canvases {
    fmt.Println(c.ID, c.Name)
}
```

### Example: Filter Widgets by Type and Parent

```go
filter := &canvus.Filter{Criteria: map[string]interface{}{
    "widget_type": "note",
    "parent_id":   "parent123",
}}
widgets, err := session.ListWidgets(ctx, canvasID, filter)
if err != nil {
    // handle error
}
for _, w := range widgets {
    fmt.Println(w.ID, w.WidgetType, w.ParentID)
}
```

### Advanced: Wildcards and JSONPath-like Selectors

- Use `"*"` as a value to match any value for a field.
- Use `"$.location.x"` to match nested fields (e.g., widgets with a specific X coordinate).

```go
filter := &canvus.Filter{Criteria: map[string]interface{}{
    "widget_type": "*", // any widget type
    "$.location.x": 100.0, // widgets at x=100
}}
widgets, err := session.ListWidgets(ctx, canvasID, filter)
```

## Widget Search Across All Canvases

You can search for widgets matching a query across all canvases using `FindWidgetsAcrossCanvases`. This utility supports exact, wildcard, prefix/suffix, and partial (contains) string matches, as well as nested field selectors.

### Example: Find All Browser Widgets with URL Suffix

```go
query := map[string]interface{}{
    "widget_type": "browser",
    "url": "*12345", // matches any url ending with 12345
}
matches, err := canvus.FindWidgetsAcrossCanvases(ctx, session, query)
if err != nil {
    // handle error
}
for _, m := range matches {
    fmt.Println(m.CanvasID, m.WidgetID, m.Widget.WidgetType)
}
```

### Wildcard and Partial Match Support
- `"*"` matches any value
- `"abc*"` matches values starting with `abc`
- `"*123"` matches values ending with `123`
- `"*mid*"` matches values containing `mid`

### Nested Field Selectors
- Use JSONPath-like keys (e.g., `"$.location.x"`) to match nested fields.

## Import/Export and Asset Handling

- The SDK supports robust, round-trip-safe import and export of all widget and asset types (notes, images, pdfs, videos, anchors, connectors, etc.).
- Asset files (images, PDFs, videos) are exported as files and referenced in the export JSON. Import reads these files and creates the corresponding widgets with correct spatial and parent/child relationships.
- The import/export logic is fully covered by integration tests, ensuring that all widgets and assets can be exported from one canvas and imported into another with full fidelity.
- The SDK normalizes widget types and handles all required fields (location, size, etc.) for all widget types.

### Geometry Utilities Example

You can determine which widgets are spatially contained within or touch a given widget using:

```go
zone, err := canvus.WidgetsContainId(ctx, session, canvasID, anchorID, nil, 0)
if err != nil {
    // handle error
}
for _, w := range zone.Contents {
    fmt.Println(w.ID, w.WidgetType)
}
```
See godoc for more details on geometry utilities.

--- 
# Canvus Go SDK: Planned Abstractions & Utilities

This document outlines the key abstractions, utilities, and advanced features planned for the Canvus Go SDK, based on current requirements and ongoing design discussions. It is referenced from `PRD.md` and should be updated as the SDK evolves.

---

## 1. Generic Filtering & Subscription Abstraction

- **Design:** A generic filter struct will be used for all list/get endpoints (e.g., canvases, widgets). Since MCS does not support server-side filtering, all filtering will be performed client-side.
- **Features:**
  - Accepts arbitrary JSON for filter criteria.
  - Supports wildcards (`*`) and JSONPath-like (`$`) keying for flexible matching.
  - Can be reused across endpoints.

## 2. Connector Creation with Widget Spawning

- **Design:** The standard `CreateConnector` function will handle both connecting existing widgets and creating a new widget (of any type) as part of the connector creation process.
- **Features:**
  - If the `src` or `dst` array contains a full widget JSON (not just an ID), the SDK will create the widget, extract its ID, and use it for the connector.
  - Works for all widget types, provided sufficient payload information is present.

## 3. Recursive Widget Search (findWidget/findWidgetIn)

- **Design:**
  - `findWidget`: Recursively searches all canvases for widgets matching arbitrary JSON criteria.
  - `findWidgetIn`: Limits search to a specific folder or canvas, falling back to global search if not found.
- **Features:**
  - Returns full drill-down path: `{server, canvas_id, widget_id}`.
  - Supports arbitrary JSON queries, wildcards, and JSONPath-like selectors.

## 4. Response Validation & Retry Logic

- **Design:**
  - Validates that fields in the request payload match those in the server response (excluding server-generated fields like `id`).
  - Retries up to 3 times on mismatch or transient errors.
  - Does **not** retry on 401 (Unauthorized) or 403 (Forbidden); 404 is considered potentially transient and will be retried.

## 5. Client Creation Patterns

- **Types:**
  - `test_client`: Creates a new user and token, uses them for all actions, deletes both on exit.
  - `user_client`: Creates a temporary token for a specific user (by ID, name, or email), deletes token on exit.
  - `client`: Uses credentials from settings file.
- **Lifecycle:**
  - Clients can be reused for multiple actions within a process.
  - Cleanup (deleting users/tokens) occurs when the client is "closed" (end of process or explicit call).
  - **Open Question:** Options for managing client lifecycle and cleanup (see below).

## 6. Viewport & Widget Location Utilities

- **Design:**
  - Can set can you port centered on a widget (default) or at top-left (0,0) as an option.
  - `open-canvas` accepts an option to open at a specific widget or location.

## 7. Geometry Utilities (contains, touches)

- **Design:**
  - Operate on widgets within a canvas.
  - `contains`: Returns widgets entirely within the bounds of a source widget (with configurable tolerance, default as described).
  - `touches`: Returns widgets that overlap with the source widget (with tolerance).

## 8. Batch Operations (move, copy, migrate, delete, pin/unpin)

- **Design:**
  - Best-effort concurrency: operations are performed in parallel, partial failures are logged, and retry logic is applied as above.
  - `migrate`: For widgets with file content, downloads from source canvas and uploads to destination.

## 9. Open Questions & Design Notes

- **Client Lifecycle Management:**
  - Options for explicit vs. automatic cleanup.
  - How to handle multiple clients in a single process.
- **Geometry Tolerance:**
  - Default and configurable values for "contains" and "touches".
- **Migrate/Clone API Feature:**
  - See `MCS-Feature-Requests.md` for proposed API improvements.

## Workspace Abstractions (2024-06)

### Workspace Selection
- All workspace-related functions accept a `WorkspaceSelector` (by index, name, user, or default to index 0).

### Viewport Setting
- `SetWorkspaceViewport` supports:
  - Direct coordinates (x, y, width, height)
  - Centering/scaling on a widget by widgetID
- **Math for centering/scaling:**
  - Center: viewport.X = widget.X + widget.Width/2 - viewport.Width/2
  - Center: viewport.Y = widget.Y + widget.Height/2 - viewport.Height/2
  - To fit widget: scale = max(widget.Width/(viewport.Width-2*margin), widget.Height/(viewport.Height-2*margin))

### Toggle Commands
- `ToggleWorkspaceInfoPanel` and `ToggleWorkspacePinned` fetch current state and toggle it.

### Open Canvas
- `OpenCanvasOnWorkspace` supports optional x/y for viewport centering after open, and validates canvas change.
- If widgetID is provided, centers viewport on that widget after open.

---

**See also:** [PRD.md](./PRD.md)

---

## To Do

- Expand with code samples and usage patterns as SDK implementation progresses.
- Update with any new abstractions/utilities as they are identified.

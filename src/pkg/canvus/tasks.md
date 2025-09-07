# Canvus Go SDK Development Tasks

This document is the authoritative project plan. **Update this file with every new plan, approach, or after each major prompt cycle, especially before committing to git.**

## Project Objective

- Expand the Go API library to support all Canvus API endpoints and features.
- Build a full-featured Go SDK for Canvus, including utilities, documentation, and developer tools.

---

## Task List

### 1. Project Setup

- [x] Initialize a new git repository (GitHub, using CLI)
- [x] Set up standard branching and commit practices
- [x] Confirm Windows/PowerShell as the development environment
- [x] Add `.gitignore` for Go and project-specific files
- [x] Confirmed: All JSON field names in requests must be lowercase to match the Canvus API. This is required for PATCH/POST to work (e.g., canvas renaming).

**2024-06-10 Summary:**

- Confirmed existing git repository.
- Added Go-specific rules to `.gitignore` (Windows/PowerShell compatible).
- Created `CONTRIBUTING.md` with branching, commit, and Windows compatibility guidelines.

### 2. MCS REST API Analysis

- [x] Collect and review the official MT Canvus Server (MCS) REST API documentation/spec
- [x] List all available API endpoints and their parameters
- [x] Identify required abstractions/utilities to improve developer experience (see [Abstractions.md](./Abstractions.md))
- [x] Map API endpoints to planned SDK features and abstractions (see [Abstractions.md](./Abstractions.md))
- [x] Document findings and update the PRD if needed (see [PRD.md](./PRD.md))

**2024-06-13 Summary:**

- Completed identification, mapping, and documentation of required SDK abstractions and utilities. See [Abstractions.md](./Abstractions.md) for details and [PRD.md](./PRD.md) for cross-references.

### 3. API Coverage Analysis

- [x] Extract all endpoints and features from the Docs

### 4. Go API Library Expansion

For each missing endpoint/feature:

- [x] Define Go method signature and data structures
- [x] Implement the method (for all endpoints below unless noted)
- [x] Add error handling and authentication
- [x] Write unit/integration tests for all endpoints (full CRUD, partial, and custom endpoints, including widgets GET, are now covered by robust integration tests)
- [x] Update documentation (godoc, README, code samples) (Completed 2024-06-16: Full godoc audit, README, and code samples updated)

#### Endpoint Implementation & Testing Order (2024-06-15)

1. **System Management Endpoints**
    - [x] Users: implement and test all actions
    - [x] Access Tokens: implement and test all actions
    - [x] Groups: implement and test all actions
    - [x] Canvas Folders: implement and test all actions
    - [x] Server Config: implement and test all actions
    - [x] License: implement (do not test activation)
    - [x] Audit Log: implement and test all actions
    - [x] Server Info: implement and test all actions
2. **Canvas Endpoints**
    - [x] Implement and test all Canvas actions (CRUD, move, copy, permissions, etc.)
3. **Client & Workspace Endpoints**
    - [x] Implement and test all Client actions
    - [x] Implement and test all Workspace actions
    - [x] Implement logic to launch MT-Canvus-Client with canvas URL for integration tests [ABROGATED: Not part of Go SDK; for integration/E2E testing only. For API tests, ensure a client is active and running. See documentation.]
    - [x] Refactor open-canvas code to poll get workspaces and verify canvas ID update (server does not send full response; must check for matching canvas ID)
4. **Widget & Asset Endpoints**
    - [x] Notes: all CRUD actions implemented
    - [x] Anchors: all CRUD actions implemented
    - [x] Images: all CRUD actions implemented
    - [x] Connectors: all CRUD actions implemented
    - [x] Color Presets: get/patch implemented
    - [x] Uploads Folder: note and asset upload implemented
    - [x] PDFs: all CRUD actions implemented
    - [x] Videos: all CRUD actions implemented
    - [x] Video Inputs: all actions implemented
    - [x] Video Outputs: all actions implemented
    - [x] Mipmaps: info, level, and asset retrieval implemented
    - [x] Backgrounds: get/patch/post implemented
    - [x] Assets: all CRUD actions implemented
    - [x] Parenting (patching parent ID): method stubbed, not tested (known bug)
    - [x] Read-only endpoints (Widgets, Annotations): widgets read-only done, annotations skipped (API incomplete)

#### Testing & Cleanup Policy

- [x] All tests must clean up (permanently delete) resources they create, even on failure. Moving to trash is not sufficient.
- [x] Each test must use unique resource names/IDs to avoid collisions and ensure safe cleanup.

### Status Summary (2024-06-15)

- **All major system, canvas, client, workspace, widget, and asset endpoints are implemented in the Go SDK.**
- **Unit tests and documentation updates are in progress for new endpoints.**
- **Parenting and Annotations endpoints are not fully implemented/tested due to API limitations or known bugs.**
- **Next:**
    - Finalize godoc and README documentation
    - Review for any missing minor endpoints or edge cases
    - Prepare for initial release
    - (Optional) Implement and test widget-centric viewport and open-canvas centering in workspace logic
    - (Optional) Revisit parenting and annotations endpoints if/when API/server supports it

### Required Abstractions/Utilities

- **Authentication:** API key (from env/config), secure handling
- **Context support:** All requests accept `context.Context`
- **Error handling:** Centralized, idiomatic Go error types
- **Pagination/streaming:** Helpers for paginated and streaming endpoints
- **Strong typing:** Request/response models as Go structs
- **Modular structure:** Packages by resource (canvases, folders, widgets, users, etc.)
- **CLI utilities:** (Optional) for common workflows
- **Documentation:** Godoc, README, code samples
- [x] **Testing:** Unit/integration tests for all features

### 5. Build the Go SDK

- [ ] Identify and design SDK utilities (CLI tools, helpers, etc.) (Next focus as of 2024-06-16)
- [ ] Implement SDK features
- [ ] Add code samples and templates
- [ ] Write comprehensive documentation

### 6. Developer Experience & Release

- [ ] Package as a Go module
- [ ] Ensure easy installation and usage
- [ ] Finalize documentation and examples
- [ ] Tag and release initial version

### 7. Advanced Abstractions & Utilities (from Abstractions.md)

> The following advanced abstractions/utilities are required for full parity with the design in Abstractions.md. Work will proceed through these in order. (2024-06-16)

- [x] **Generic Filtering Abstraction**
    - Implemented as Filter struct with support for arbitrary JSON, wildcards, and JSONPath-like selectors for all list/get endpoints. Documented in README and code. Complete as of 2024-06-16.
- [x] **Connector Creation with Widget Spawning**
    - CreateConnector supports widget creation for src/dst if provided as widget JSON. Complete as of 2024-06-16.
- [x] **Recursive Widget Search (findWidget/findWidgetIn)**
    - Implemented as FindWidgetsAcrossCanvases and FilterSlice utilities, supporting recursive and cross-canvas widget search with flexible queries. Complete as of 2024-06-16.
- [x] **Response Validation & Retry Logic**
    - Centralized response validation and retry logic implemented in session.go, with relaxed numeric comparison and robust error handling. Complete as of 2024-06-16.
- [x] **Geometry Utilities (contains, touches)**
    - Add utilities to determine widget containment and overlap within a canvas.
    - Implemented WidgetsContainId with explicit canvasID requirement and detailed documentation. Usage examples provided in godoc.
- [ ] **Batch Operations Utility**
    - Implement batch move/copy/migrate/delete/pin/unpin with concurrency, partial failure logging, and retry logic.

### CLI Command Roadmap
- [ ] widget list
- [ ] widget get
- [ ] widget create
- [ ] widget update
- [ ] widget delete
- [ ] widget contains
- [ ] widget touches
- [ ] user list
- [ ] user get
- [ ] user create
- [x] user delete (Implemented `cleanup_test_users.go`)
- [ ] user activate
- [ ] export region
- [ ] export ids
- [ ] import

> **Note:**
> All core SDK features, abstractions, and utilities are now implemented and tested. The next major focus is implementing the above CLI commands. As of 2024-06-16, only the CLI scaffold and a stub for the cleanup command exist; all other CLI commands are pending implementation.
> The `cleanup_test_users.go` script has been implemented to list and delete users whose emails start with "testuser".

### SDK Abstractions (Current Focus)
- [x] Implement import/export abstractions in SDK (export.go, import.go) (Completed 2024-06-16: Supports all widget types and assets)
- [x] Implement and verify import/export round-trip tests (use /tests/importdata and /tests/exportdata for test data) (Completed 2024-06-16: Round-trip is robust, all widget and asset types are supported, and test coverage is complete. Widget type handling is case-insensitive, and numeric response validation is relaxed for cross-type equality.)

### 8. Refactor Widget CRUD to Dispatch by Type
- [x] Refactor widget CRUD to dispatch by type (Completed 2024-06-16)
### 9. Scaffold CLI and Implement Cleanup Command
- [x] Scaffold CLI and implement cleanup command (Completed 2024-06-16)

---

## 2024-06-11: Greenfield Go SDK Approach & Planning

### New Approach

- The Canvus Go SDK will be built from scratch, following modern Go community best practices.
- Old Python/Go libraries are not used as a roadmap; they serve only as historical context.
- The SDK will be idiomatic, modular, and designed for the broader Go developer community.
- Focus: developer experience, full API coverage, strong documentation, and extensibility.

### Required Abstractions/Utilities

- **Authentication:** API key (from env/config), secure handling
- **Context support:** All requests accept `context.Context`
- **Error handling:** Centralized, idiomatic Go error types
- **Pagination/streaming:** Helpers for paginated and streaming endpoints
- **Strong typing:** Request/response models as Go structs
- **Modular structure:** Packages by resource (canvases, folders, widgets, users, etc.)
- **CLI utilities:** (Optional) for common workflows
- **Documentation:** Godoc, README, code samples
- **Testing:** Unit/integration tests for all features

### API Endpoint → Planned Go SDK Feature Mapping (WIP)

| API Resource         | HTTP Method & Path                              | Planned Go SDK Method Signature                |
|---------------------|-------------------------------------------------|------------------------------------------------|
| Canvases            | GET    /canvases                                | func (c *Client) ListCanvases(ctx context.Context) ([]Canvas, error) |
| Canvases            | GET    /canvases/:id                            | func (c *Client) GetCanvas(ctx context.Context, id string) (Canvas, error) |
| Canvases            | POST   /canvases                                | func (c *Client) CreateCanvas(ctx context.Context, req CreateCanvasRequest) (Canvas, error) |
| Canvases            | PATCH  /canvases/:id                            | func (c *Client) UpdateCanvas(ctx context.Context, id string, req UpdateCanvasRequest) (Canvas, error) |
| Canvases            | DELETE /canvases/:id                            | func (c *Client) DeleteCanvas(ctx context.Context, id string) error |
| ...                 | ...                                             | ...                                            |

> This table will be expanded to cover all endpoints in the official API list. Each resource group will have its own Go type(s) and methods, following idiomatic Go SDK design.

---

## Authentication Methods (SDK & API)

There are two primary ways to authenticate to the Canvus server:

1. **Username/Password Login**
   - Endpoint: `POST /users/login`
   - The client sends a username (email) and password to the endpoint.
   - If password authentication is enabled, the server issues a temporary access token (valid for 24 hours).
   - This token is used for subsequent authenticated requests.
   - Example:

     ```json
     POST /users/login
     { "email": "alice@example.com", "password": "BBBB" }
     ```

   - Response includes a `token` and user info.

2. **Access Token**
   - Endpoint: `POST /access-tokens` (or via Canvus web UI)
   - An access token is created via the API or UI.
   - This token does **not expire** and can be used directly for authentication by including it in the `Private-Token` header.
   - You can also POST an existing token to `/users/login` to validate and prolong its lifetime.

3. **Sign Out**
   - Endpoint: `POST /users/logout`
   - Invalidates the provided access token. If no token is provided in the body, the `Private-Token` header is used.

### SDK Client Authentication Options

The Go SDK client supports three authentication options:

1. Username/password login (temporary token)
2. Static access token (long-lived)
3. Token validation/refresh (prolongs token lifetime)

All authentication logic must be tested for both login and static token flows.

---

## Client Instantiation & Authentication Patterns

The SDK must support three primary client instantiation patterns:

1. **test_client**
   - Uses the main client (from settings) to create and activate a new test user.
   - Logs in as the test user to obtain a temporary token (API key/PrivateToken) via `/users/login`.
   - All actions in the session use this token (sent as `
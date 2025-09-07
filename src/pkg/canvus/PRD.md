# Canvus Go SDK – Product Requirements Document (PRD)

## 1. Project Overview

- **Goal:**  
  Build a robust, idiomatic Go SDK for the Canvus API, providing full API coverage, developer-friendly abstractions, and a seamless experience for Go developers on Windows/PowerShell.
- **Scope:**  
  - Expand the existing Go API library to cover all Canvus API endpoints and features.
  - Package the library as a full SDK, including utilities, documentation, and developer tools.
  - Ensure all development is Windows/PowerShell compatible.
  - Maintain high code quality, test coverage, and clear documentation.

---

## 2. Stakeholders

- **Primary:**  
  - Go developers integrating with Canvus.
  - Internal Canvus engineering and QA teams.
- **Secondary:**  
  - Technical writers (for docs).
  - DevOps (for packaging and release).

---

## 3. Functional Requirements

### 3.1. API Coverage

- The SDK must provide Go methods for every public Canvus API endpoint, including:
  - Canvas management (CRUD)
  - Folder management
  - Widget operations (notes, images, browsers, videos, PDFs, connectors, anchors)
  - User and token management
  - Client and workspace operations
  - Server info/config endpoints
  - Asset management (images, videos, PDFs, uploads, backgrounds, mipmaps, etc.)
  - Parenting (patching parent ID) must be implemented but is not to be tested due to a known server bug.
- All request/response models must be idiomatic Go structs, with proper error handling.
- **Import/Export:** The SDK supports robust, round-trip import/export for all widget and asset types, with case-insensitive widget type handling and relaxed numeric validation for cross-type equality.

### 3.2. SDK Features

- **Authentication:**  
  - Support API key authentication (from env or config).
- **Utilities:**  
  - CLI tools for common tasks (e.g., token management, canvas export/import).
  - Helper functions for common workflows (e.g., batch operations, event subscriptions).
- **Documentation:**  
  - Godoc comments for all exported types and functions.
  - Comprehensive README with setup, usage, and examples.
  - Code samples for common use cases.
- **Testing:**  
  - Unit tests for all major features, following Go testing conventions.
  - Test coverage for happy path, edge cases, and error handling.

### 3.3. Developer Experience

- Easy installation via Go modules.
- Windows/PowerShell compatibility for all scripts and commands.
- Clear error messages and logging.
- Consistent, idiomatic Go API design.

### 3.4. Endpoint Implementation & Testing Order

The following order must be followed for endpoint implementation and testing, to ensure safe, isolated, and reliable test environments:

1. **System Management Endpoints**: Users, Access Tokens, Groups, Canvas Folders, Server Config, License (no activation tests), Audit Log, Server Info
2. **Canvas Endpoints**: All canvas actions (CRUD, move, copy, permissions, etc.)
3. **Client & Workspace Endpoints**: All client and workspace actions, including launching the MT-Canvus-Client with a canvas URL and verifying client/workspace state
4. **Widget & Asset Endpoints**: All widget and asset actions, in dependency order (simple elements first, then those requiring files or references, then read-only endpoints)
5. **Parenting**: Implement but do not test parenting (patching parent ID)

---

### Test Execution Order

To ensure reliable, isolated, and maintainable test runs, tests must be executed in the following order, matching the development and implementation sequence:

1. **System Management Tests**
    - Users
    - Access Tokens
    - Groups
    - Canvas Folders
    - Server Config
    - License (no activation tests)
    - Audit Log
    - Server Info
2. **Canvas Tests**
    - All Canvas actions (CRUD, move, copy, permissions, etc.)
3. **Client & Workspace Tests**
    - All Client actions
    - All Workspace actions
    - MT-Canvus-Client launch and workspace state verification
4. **Widget & Asset Tests**
    - All Widget actions (Notes, Anchors, VideoInputs, VideoOutputs, Color Presets, etc.)
    - All Asset actions (Images, Videos, PDFs, Uploads, Connectors, Backgrounds, MipMaps, Assets, etc.)
    - All read-only endpoints (Widgets, Annotations)
5. **Parenting**
    - Implement but do not test parenting (patching parent ID)

This order must be followed for both development and CI test execution to avoid resource conflicts and ensure proper cleanup.

### 3.5. Subscription & Buffering Requirements

- For endpoints supporting subscriptions (with `?subscription`), the SDK must:
  - Support initial GET and real-time update streaming (one JSON per line, CR as keep-alive)
  - Provide a function to filter updates for specific elements
  - Provide a buffered subscription handler: only emit updates after a configurable period of inactivity, to reduce noise from rapid, sequential updates

### 3.6. Testing & Cleanup Policy

- All tests must clean up (permanently delete) resources they create, even on failure. Moving to trash is not sufficient.
- Each test must use unique resource names/IDs to avoid collisions and ensure safe cleanup.

### 3.7. User-Specific Resource Handling (NEW)

- After login, the SDK extracts and stores the authenticated user's ID from the /users/login response.
- The user ID is used to construct user-specific resource IDs, such as:
  - Trash folder: trash.{userId}
  - User root folder: {userId}
- All SDK methods that operate on user-specific folders (e.g., TrashCanvas, TrashFolder) now use this approach for correctness and multi-user safety.
- This pattern should be used for any future user-specific resource operations.

### 3.8. Client Types and Test Policy (NEW)

- The SDK and tests use three client types:
  1. **Admin Client**: Uses the API key from settings.json. Only used for admin operations: user create/delete, unblock, approve, audit-log, licenses.
  2. **TestClient**: Used for all non-admin tests. Creates a temporary user, logs in, and performs all actions as that user. Cleans up after tests.
  3. **UserClient**: Used for tests that require login as an existing user (not common in standard test flows).
- All Canvas, Folder, Widget, and Asset tests must use the TestClient, not the admin client.
- This policy ensures correct permissions, test isolation, and future maintainability.
- See clients.go for implementation details.

## API/SDK Design Notes

- All JSON field names sent to the Canvus API must be lowercase (e.g., 'name', 'mode'), matching the API's requirements. This is critical for PATCH/POST requests to work as expected.
- **BREAKING CHANGE: Rename Client to Session**
  - The SDK's main struct, previously named Client, will be renamed to Session throughout the codebase.
  - Rationale: The Canvus API has a 'clients' resource, which can cause confusion for SDK users. 'Session' more accurately reflects the purpose of the struct (an authenticated API session) and avoids ambiguity.
  - All references, documentation, and tests will be updated accordingly.
  - Migration note: This is a breaking change and will require users to update their code.

---

## 4. Non-Functional Requirements

- **Performance:**  
  - Efficient HTTP requests, minimal allocations, and fast response parsing.
- **Reliability:**  
  - Robust error handling and retries for transient failures.
- **Maintainability:**  
  - Modular code structure, clear separation of concerns.
  - Adherence to Go best practices and project coding standards.
- **Security:**  
  - Secure handling of API keys and sensitive data.
- **Versioning:**  
  - Semantic versioning for releases.

---

## 5. Constraints

- All development and documentation must be Windows/PowerShell compatible.
- No Linux shell commands or assumptions.
- Use only well-maintained, widely adopted Go modules.
- All code and documentation must be committed to git, following the project's branching and commit message conventions.
- The `tasks.md` file must be updated with every new plan, approach, or after each major prompt cycle, and before every git commit.

---

## 6. Milestones & Deliverables

1. **Project Setup**
   - GitHub repo, `.gitignore`, initial planning docs.
2. **API Coverage Analysis**
   - Coverage matrix (Python vs Go), gap analysis.
3. **Go API Library Expansion**
   - Full endpoint coverage, models, and tests.
4. **SDK Utilities & CLI**
   - Helper tools, code samples, and documentation.
5. **Release**
   - Go module packaging, README, and tagged release.
- [ ] BREAKING: Refactor Client to Session throughout SDK and documentation

---

## 7. Success Criteria

- 100% coverage of public Canvus API endpoints.
- All code passes unit tests and is reviewed for idiomatic Go style.
- SDK is installable and usable by Go developers on Windows.
- Documentation is clear, complete, and up-to-date.
- No Linux-specific commands or scripts in the codebase.
- **Import/Export round-trip is robust and fully tested for all widget and asset types. Widget type handling is case-insensitive, and numeric response validation is relaxed for cross-type equality.**

---

## 8. Open Questions / To Be Determined

- Are there any private/internal Canvus endpoints that should be included?
- **Minimum Go version to support:** Go 1.24.1 (per developer environment)
- **Frameworks/Libraries:**
  - Use standard `net/http` for HTTP operations
  - Use `github.com/go-playground/validator/v10` for data validation if needed
  - Use `GORM` or `sqlc` for ORM/database interaction if required
  - Follow all conventions in @10-golang-coding-standards.mdc
- **Release cadence:** SDK will be updated twice a year, aligned with MT Canvus releases

---

## Authentication (API & SDK)

The Canvus API and SDK support two main authentication flows:

### 1. Username/Password Login

- Endpoint: `POST /users/login`
- The client sends an email and password to the endpoint.
- If password authentication is enabled, the server issues a temporary access token (valid for 24 hours).
- Example request:

  ```json
  POST /users/login
  { "email": "alice@example.com", "password": "BBBB" }
  ```

- Example response:

  ```json
  {
    "token": "...",
    "user": { ... }
  }
  ```

### 2. Access Token

- Endpoint: `POST /access-tokens` (or via Canvus web UI)
- An access token is created via the API or UI.
- This token does **not expire** and can be used directly for authentication by including it in the `Private-Token` header.
- You can also POST an existing token to `/users/login` to validate and prolong its lifetime.

### 3. Sign Out

- Endpoint: `POST /users/logout`
- Invalidates the provided access token. If no token is provided in the body, the `Private-Token` header is used.

### SDK Client Authentication Options

The Go SDK client supports three authentication options:

1. Username/password login (temporary token)
2. Static access token (long-lived)
3. Token validation/refresh (prolongs token lifetime)

The SDK must allow the user to choose their preferred authentication method at initialization.

---

## SDK Abstractions & Utilities

See [Abstractions.md](./Abstractions.md) for a detailed description of planned SDK abstractions, utilities, and advanced features.

## MCS API Feature Requests

See [MCS-Feature-Requests.md](./MCS-Feature-Requests.md) for proposed improvements and suggestions for the MCS API.

---

## Client Instantiation & Authentication Patterns

The SDK must support three primary client instantiation patterns:

1. **test_client**
   - Uses the main client (from settings) to create and activate a new test user.
   - Logs in as the test user to obtain a temporary token (API key/PrivateToken) via `/users/login`.
   - All actions in the session use this token (sent as `Private-Token` header).
   - On completion, logs out (`/users/logout`, invalidates token) and deletes the test user.

2. **user_client**
   - Logs in as an existing user (by email and password) to obtain a temporary token via `/users/login`.
   - All actions use this token for the session (sent as `Private-Token` header).
   - On completion, logs out to invalidate the token.

3. **client**
   - Uses credentials from the settings/config file for all actions.
   - No automatic cleanup or user/token creation.

**Notes:**

- The terms "API key" and "PrivateToken" are used interchangeably; all authentication headers use `Private-Token`.
- Session cleanup for temporary tokens is handled by calling `/users/logout`.

---

## Go SDK Type Definition Approach

- **Resource-specific types** (e.g., ClientInfo, AccessToken, User) are defined in their respective files (e.g., clients.go, accesstokens.go, users.go).
- **Only truly shared types** (used by multiple resources/packages) are kept in types.go.
- This approach ensures modularity, maintainability, and idiomatic Go design, and avoids type redeclaration errors.

---

## Configuration Management Approach

- The SDK and all tests use `settings.json` (not .env or environment variables) for configuration.
- `settings.json` must include:
  - `api_base_url`: The full base URL for the Canvus API (including protocol, e.g., https://...)
  - `api_key`: The admin or test API key
  - `test_user`: An object with `username` and `password` for test user flows
  - Any other required settings (timeouts, enabled features, etc.)
- **Environment variables are not the source of truth for configuration.**
- All test and client code should read from `settings.json` for configuration values.
- This approach ensures consistency, avoids missing/empty variable errors, and is Windows/PowerShell friendly.

---

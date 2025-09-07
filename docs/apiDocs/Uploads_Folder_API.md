# Uploads Folder API

## Overview
The uploads folder API provides endpoints for managing uploaded files and assets in the Canvus system.

## Endpoints

### GET /api/v1/uploads
List all uploaded files in the uploads folder.

**Response:**
```json
[
  {
    "id": "string",
    "filename": "string", 
    "hash": "string",
    "size": "number",
    "uploaded_at": "string",
    "content_type": "string"
  }
]
```

### GET /api/v1/uploads/{id}
Get details of a specific uploaded file.

**Response:**
```json
{
  "id": "string",
  "filename": "string",
  "hash": "string", 
  "size": "number",
  "uploaded_at": "string",
  "content_type": "string",
  "url": "string"
}
```

## Notes
- Uploaded files have hash values that can be used for asset identification
- Files in the uploads folder may be referenced by widgets and canvases
- Hash values are used for deduplication and asset management

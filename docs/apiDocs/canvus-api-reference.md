# Canvus Server API Reference

## Authentication

### POST /users/login
Signs the user in and issues an access token. This endpoint is available without authentication.

**Request Body:**
```json
{
  "email": "alice@example.com",
  "password": "BBBB"
}
```

**Response:**
```json
{
  "token": "lmFU9obmM5v4o6jdCXsRW6v5bLD9w47aGIP4eMRnf3A",
  "user": {
    "admin": false,
    "approved": true,
    "blocked": false,
    "created_at": "2021-07-02T06:38:37.141Z",
    "email": "alice@example.com",
    "id": 1002,
    "last_login": "",
    "name": "Alice",
    "state": "normal"
  }
}
```

### POST /users/logout
Signs the user out by invalidating the provided access token.

**Headers:**
```
Private-Token: <access token>
```

**Request Body:**
```json
{
  "token": "z_Ttm-tcFpiadMUR2A_8kQnkOsl6wmcEKplotULC9fk"
}
```

## Canvas Operations

### GET /canvases
Get all canvases that the authenticated user can access.

**Headers:**
```
Private-Token: <access token>
```

### GET /canvases/{id}/widgets
Get all widgets within a specific canvas.

**Headers:**
```
Private-Token: <access token>
```

## Widget Asset Information

Widgets contain asset information in the following format:

```json
{
  "hash": "347b3c308971",
  "widget_type": "Pdf",
  "original_filename": "document.pdf"
}
```

**Key Fields:**
- `hash`: The asset hash used as filename
- `widget_type`: Type of widget (Pdf, Image, Video, etc.)
- `original_filename`: Original filename with extension

**Media Asset Types:**
- `Pdf`: PDF documents
- `Image`: Image files (various formats)
- `Video`: Video files (various formats)

## Rate Limiting
- Start with 100 requests per second
- Scale back to 50/sec if encountering issues
- Further scale back to 25/sec if still unstable
- Implement retry logic for failed requests

# Canvas Backgrounds API

## Overview
The canvas backgrounds API provides endpoints for managing canvas background images and assets.

## Endpoints

### GET /api/v1/canvases/{id}/background
Get the background configuration for a specific canvas.

**Response:**
```json
{
  "type": "string",
  "image": {
    "hash": "string",
    "fit": "string"
  },
  "color": "string",
  "grid": {
    "enabled": "boolean",
    "size": "number"
  },
  "haze": {
    "color1": "string",
    "color2": "string"
  }
}
```

### PUT /api/v1/canvases/{id}/background
Update the background configuration for a canvas.

## Notes
- Canvas backgrounds can contain image assets with hash values
- Background images are stored as assets and referenced by hash
- These assets need to be included in asset discovery and restoration

# MCS API Feature Requests

This document collects feature requests and suggestions for the MCS API, to be shared with the MT team.

---

## 1. Migration/Clone API for File Content

### Problem
Currently, migrating or copying widgets with file content between canvases requires downloading the file from the source canvas and re-uploading it to the destination. This is inefficient and error-prone.

### Proposed Solution
- Add an API endpoint or option to allow direct migration or cloning of widgets (including their file content) between canvases, using file hashes or server-side references.
- This would eliminate the need for client-side download/upload and improve performance and reliability.

### Benefits
- Faster, more reliable migration of content between canvases.
- Reduced bandwidth and client complexity.
- Enables advanced SDK features (e.g., bulk migration, cloning, backup/restore).

---

## To Do
- Add additional feature requests as they are identified during SDK development. 
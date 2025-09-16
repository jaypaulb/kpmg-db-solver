# KPMG DB Solver (Non-Admin Version)

A Go-based tool for identifying missing Canvus assets and locating them in backup folders.

## Overview

KPMG DB Solver addresses the critical issue of missing asset files in Canvus Server deployments. When asset files are missing from the filesystem but still referenced in the database, widgets fail to load. This non-admin version identifies missing assets and reports their backup locations without requiring elevated privileges.

## Features

- **Automated Asset Discovery**: Queries Canvus API to identify all referenced assets
- **Missing Asset Detection**: Compares API data with filesystem to find missing files
- **Backup Search & Location Reporting**: Searches multiple backup locations and reports file locations
- **Comprehensive Reporting**: Generates detailed reports and CSV exports with backup information
- **Parallel Processing**: Efficient handling of thousands of assets and canvases
- **No Admin Privileges Required**: Read-only access to system directories
- **Windows Deployment**: Standalone executable for Windows 11/Server

## Quick Start

### Prerequisites

- Windows 11 or Windows Server
- Read access to Canvus Server directories (no admin privileges required)
- Access to Canvus Server 3.3.0
- Backup locations with asset files

### Installation

1. Download the latest release
2. Extract to desired location
3. Run normally (no administrator privileges required)

### Usage

```bash
# Run the tool
kpmg-db-solver.exe

# Follow the interactive prompts:
# 1. Enter Canvus Server URL
# 2. Enter username and password
# 3. Specify assets folder path
# 4. Specify backup root folder path
# 5. Choose report options
```

## Configuration

The tool uses interactive prompts for configuration. Key settings:

- **Canvus Server URL**: Full URL to your Canvus Server API
- **Assets Folder**: Path to the active Canvus assets directory
- **Backup Root Folder**: Root directory containing backup folders
- **Verbose Logging**: Optional detailed logging for troubleshooting

## Output

The tool generates:

1. **Detailed Report**: Missing assets grouped by canvas with widget information and backup locations
2. **CSV Export**: Comprehensive list of missing assets with backup status and file locations
3. **Backup Location Report**: All backup file locations for assets that can be restored

## Limitations

This non-admin version provides read-only access and cannot restore assets. For asset restoration, you would need:
- Administrator privileges
- The full version of the tool
- Write access to the Canvus assets directory

## Architecture

- **Go 1.21+**: High-performance, concurrent processing
- **Canvus SDK**: Proven API integration with Canvus Server 3.3.0
- **Parallel Processing**: Simultaneous API calls and filesystem operations
- **Modular Design**: Clean separation of concerns for maintainability

## Documentation

- [Product Requirements](docs/PRD.md)
- [Technical Architecture](docs/TECH_STACK.md)
- [Development Tasks](docs/TASKS.md)

## Support

For issues, questions, or contributions, please refer to the project documentation or contact the development team.

## License

This project is proprietary software developed for KPMG internal use.

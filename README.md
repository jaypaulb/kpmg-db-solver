# KPMG DB Solver

A Go-based tool for identifying and restoring missing Canvus assets from backup locations.

## Overview

KPMG DB Solver addresses the critical issue of missing asset files in Canvus Server deployments. When asset files are missing from the filesystem but still referenced in the database, widgets fail to load. This tool automatically identifies missing assets and restores them from backup locations.

## Features

- **Automated Asset Discovery**: Queries Canvus API to identify all referenced assets
- **Missing Asset Detection**: Compares API data with filesystem to find missing files
- **Backup Search & Recovery**: Searches multiple backup locations and restores missing assets
- **Comprehensive Reporting**: Generates detailed reports and CSV exports
- **Parallel Processing**: Efficient handling of thousands of assets and canvases
- **Windows Deployment**: Standalone executable for Windows 11/Server

## Quick Start

### Prerequisites

- Windows 11 or Windows Server
- Administrator privileges
- Access to Canvus Server 3.3.0
- Backup locations with asset files

### Installation

1. Download the latest release
2. Extract to desired location
3. Run as Administrator

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

1. **Detailed Report**: Missing assets grouped by canvas with widget information
2. **CSV Export**: List of missing asset filenames for restoration tracking
3. **Restoration Log**: Summary of successfully restored assets
4. **Error Report**: Assets that couldn't be found in any backup location

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

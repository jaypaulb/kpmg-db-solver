# Build System

This directory contains the build system and tools for the KPMG DB Solver project.

## Files

- `Makefile` - Main build automation with cross-platform support
- `build-windows.sh` - PowerShell script for Windows builds
- `README.md` - This documentation

## Quick Start

### Using Make (Recommended)

```bash
# Build for current platform
make build

# Build for Windows
make build-windows

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Using Scripts

```bash
# Windows build script
./build-windows.sh
```

## Build Targets

| Target | Description |
|--------|-------------|
| `build` | Build for current platform |
| `build-windows` | Build Windows executable (.exe) |
| `build-linux` | Build Linux executable |
| `build-mac` | Build macOS executable |
| `build-all` | Build for all platforms |
| `test` | Run all tests |
| `test-coverage` | Run tests with coverage report |
| `lint` | Run code formatting and linting |
| `clean` | Remove all build artifacts |
| `deps` | Update Go dependencies |
| `help` | Show available targets |

## Output

All built executables are placed in the `../bin/` directory:

- `kpmg-db-solver` - Linux/macOS executable
- `kpmg-db-solver.exe` - Windows executable
- `kpmg-db-solver-linux` - Linux executable (cross-compiled)
- `kpmg-db-solver-mac` - macOS executable (cross-compiled)

## Requirements

- Go 1.21 or later
- Make (for Makefile targets)
- PowerShell (for Windows build script)

## Cross-Compilation

The build system supports cross-compilation for all major platforms:

- **Windows**: `GOOS=windows GOARCH=amd64`
- **Linux**: `GOOS=linux GOARCH=amd64`
- **macOS**: `GOOS=darwin GOARCH=amd64`

All builds use `CGO_ENABLED=0` for static linking and better compatibility.

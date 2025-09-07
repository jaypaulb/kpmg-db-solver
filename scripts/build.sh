#!/bin/bash

# KPMG DB Solver Build Script
# Builds Windows executable from Linux development environment

set -e

echo "🔨 Building KPMG DB Solver..."

# Set build variables
BINARY_NAME="kpmg-db-solver"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.goVersion=${GO_VERSION}"

echo "📋 Build Information:"
echo "  Version: ${VERSION}"
echo "  Build Time: ${BUILD_TIME}"
echo "  Go Version: ${GO_VERSION}"
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f ${BINARY_NAME} ${BINARY_NAME}.exe
rm -rf build/

# Create build directory
mkdir -p build/

# Build for Windows
echo "🏗️  Building Windows executable..."
GOOS=windows GOARCH=amd64 go build \
    -ldflags "${LDFLAGS}" \
    -o build/${BINARY_NAME}.exe \
    ./src/cmd/kpmg-db-solver

# Build for Linux (development)
echo "🏗️  Building Linux executable..."
go build \
    -ldflags "${LDFLAGS}" \
    -o build/${BINARY_NAME} \
    ./src/cmd/kpmg-db-solver

# Verify builds
echo "✅ Verifying builds..."
if [ -f "build/${BINARY_NAME}.exe" ]; then
    echo "  ✓ Windows executable: build/${BINARY_NAME}.exe"
    ls -lh build/${BINARY_NAME}.exe
else
    echo "  ✗ Windows build failed"
    exit 1
fi

if [ -f "build/${BINARY_NAME}" ]; then
    echo "  ✓ Linux executable: build/${BINARY_NAME}"
    ls -lh build/${BINARY_NAME}
else
    echo "  ✗ Linux build failed"
    exit 1
fi

echo ""
echo "🎉 Build completed successfully!"
echo "   Windows executable: build/${BINARY_NAME}.exe"
echo "   Linux executable: build/${BINARY_NAME}"
echo ""
echo "📦 To test the Windows build:"
echo "   wine build/${BINARY_NAME}.exe --help"
echo ""
echo "📦 To test the Linux build:"
echo "   ./build/${BINARY_NAME} --help"

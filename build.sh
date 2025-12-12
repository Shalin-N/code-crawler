#!/bin/bash

# Build script for Code Crawler
# Builds executables for multiple platforms

VERSION="1.0.0"
APP_NAME="code-crawler"

echo "Building Code Crawler v${VERSION}"
echo "=================================="

# Create dist directory
mkdir -p dist

# Build for different platforms
platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    output_name="${APP_NAME}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "dist/${output_name}" .
    
    if [ $? -eq 0 ]; then
        echo "✓ Successfully built dist/${output_name}"
    else
        echo "✗ Failed to build for ${GOOS}/${GOARCH}"
    fi
done

echo ""
echo "Build complete! Binaries are in the dist/ directory"
echo ""
ls -lh dist/

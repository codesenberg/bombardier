#!/bin/bash
echo "Building Bombardier UI..."

GOOS=$(uname | tr '[:upper:]' '[:lower:]')
GOARCH=amd64

go build -o ../../bombardier-ui

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build completed successfully!"
echo "You can now run ./bombardier-ui"

# Make executable
chmod +x ../../bombardier-ui 
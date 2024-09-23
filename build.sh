#!/bin/bash

# Create distribution directories
mkdir -p dist/invoice-generator-macos/fonts
mkdir -p dist/invoice-generator-windows/fonts

# Ensure font files exist
if [ ! -f "Inter/Inter Variable/Inter.ttf" ] || [ ! -f "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" ]; then
    echo "Error: Font files not found. Please ensure they are in the correct location."
    exit 1
fi

# Build the application for macOS
CGO_ENABLED=1 go build -o dist/invoice-generator-macos/invoice-generator

# Build the application for Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o dist/invoice-generator-windows/invoice-generator.exe

# Copy necessary files for macOS
cp config.yaml dist/invoice-generator-macos/
cp "Inter/Inter Variable/Inter.ttf" dist/invoice-generator-macos/fonts/Inter.ttf
cp "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" dist/invoice-generator-macos/fonts/Inter-Bold.ttf

# Copy necessary files for Windows
cp config.yaml dist/invoice-generator-windows/
cp "Inter/Inter Variable/Inter.ttf" dist/invoice-generator-windows/fonts/Inter.ttf
cp "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" dist/invoice-generator-windows/fonts/Inter-Bold.ttf

# Create ZIP archives
cd dist
zip -r invoice-generator-macos.zip invoice-generator-macos
zip -r invoice-generator-windows.zip invoice-generator-windows

echo "Distribution packages created:"
echo "- dist/invoice-generator-macos.zip"
echo "- dist/invoice-generator-windows.zip"

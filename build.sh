#!/bin/bash

# Create distribution directory
mkdir -p dist/invoice-generator/fonts

# Build the application
go build -o dist/invoice-generator/invoice-generator

# Copy necessary files
cp config.yaml dist/invoice-generator/
cp "Inter/Inter Variable/Inter.ttf" dist/invoice-generator/fonts/Inter.ttf
cp "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" dist/invoice-generator/fonts/Inter-Bold.ttf

# Create ZIP archive
cd dist
zip -r invoice-generator.zip invoice-generator

echo "Distribution package created: dist/invoice-generator.zip"

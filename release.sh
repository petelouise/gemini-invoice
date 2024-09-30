#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if a version number is provided
if [ $# -eq 0 ]; then
    echo "Error: Please provide a version number (e.g., v1.0.0)"
    exit 1
fi

VERSION=$1

# Build the project
echo "Building the project..."
go build -o gemini-invoice main.go pdf.go

# Create a dist directory if it doesn't exist
mkdir -p dist

# Move the built binary to the dist directory
mv gemini-invoice dist/

# Commit changes
git add .
git commit -m "Build for release $VERSION"

# Create a new tag
git tag -a $VERSION -m "Release $VERSION"

# Push changes and tags to remote repository
git push origin main
git push origin $VERSION

# Create a GitHub release using the GitHub CLI
# Make sure you have the GitHub CLI installed and authenticated
gh release create $VERSION ./dist/* --title "Release $VERSION" --notes "Release notes for $VERSION"

echo "Release $VERSION created and published successfully!"

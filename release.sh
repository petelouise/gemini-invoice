#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if a version number is provided
if [ $# -eq 0 ]; then
    echo "Error: Please provide a version number (e.g., v1.0.0)"
    exit 1
fi

VERSION=$1

# Check if the dist directory exists and contains files
if [ ! -d "dist" ] || [ -z "$(ls -A dist)" ]; then
    echo "Error: The dist directory is missing or empty. Please build the project before running this script."
    exit 1
fi

# Create a new tag
git tag -a $VERSION -m "Release $VERSION"

# Push the tag to the remote repository
git push origin $VERSION

# Create a GitHub release using the GitHub CLI
# Make sure you have the GitHub CLI installed and authenticated
gh release create $VERSION --title "Release $VERSION" --notes "Release notes for $VERSION"

# Upload assets separately
for file in ./dist/*; do
    if [ -f "$file" ]; then
        gh release upload $VERSION "$file"
    elif [ -d "$file" ]; then
        zip -r "${file}.zip" "$file"
        gh release upload $VERSION "${file}.zip"
        rm "${file}.zip"
    fi
done

echo "Release $VERSION created and published successfully!"

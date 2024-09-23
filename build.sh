#!/bin/bash

log_timestamp() {
	echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log_timestamp "Build script started"

# Clean up existing dist directory
log_timestamp "Cleaning up dist directory"
rm -rf dist

# Create distribution directories
log_timestamp "Creating distribution directories"
mkdir -p dist/Gemini\ Invoice.app/Contents/MacOS
mkdir -p dist/Gemini\ Invoice.app/Contents/Resources/fonts
mkdir -p dist/Gemini\ Invoice\ Windows/fonts

# Create fonts directory in project root
mkdir -p fonts

# Ensure font files exist
log_timestamp "Checking font files"
if [ ! -f "fonts/Inter.ttf" ] || [ ! -f "fonts/Inter-Bold.ttf" ]; then
	echo "Error: Font files not found. Please ensure they are in the fonts directory."
	exit 1
fi

# Build the application for macOS
log_timestamp "Building for macOS"
CGO_ENABLED=1 go build -o dist/Gemini\ Invoice.app/Contents/MacOS/gemini-invoice
log_timestamp "macOS build completed"

# Build the application for Windows
log_timestamp "Building for Windows"
log_timestamp "Setting up Windows build environment"
export CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++

log_timestamp "Starting Windows build with a 10-minute timeout"
timeout 600s go build -v -o dist/Gemini\ Invoice\ Windows/gemini-invoice.exe

if [ $? -eq 124 ]; then
	log_timestamp "Error: Windows build timed out after 10 minutes"
	exit 1
elif [ $? -ne 0 ]; then
	log_timestamp "Error: Windows build failed"
	exit 1
else
	log_timestamp "Windows build completed successfully"
fi

# Create Info.plist for macOS
log_timestamp "Creating Info.plist for macOS"
cat >dist/Gemini\ Invoice.app/Contents/Info.plist <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>gemini-invoice</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.gemini.invoice</string>
    <key>CFBundleName</key>
    <string>Gemini Invoice</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.12</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOL

# Copy necessary files for macOS
log_timestamp "Copying files for macOS"
cp config.yaml dist/Gemini\ Invoice.app/Contents/Resources/
cp fonts/*.ttf dist/Gemini\ Invoice.app/Contents/Resources/fonts/

# Copy necessary files for Windows
log_timestamp "Copying files for Windows"
cp config.yaml dist/Gemini\ Invoice\ Windows/
cp fonts/*.ttf dist/Gemini\ Invoice\ Windows/fonts/

# Create a shortcut for Windows
log_timestamp "Creating Windows shortcut"
echo '@echo off
start "" "%~dp0gemini-invoice.exe"' >dist/Gemini\ Invoice\ Windows/Gemini\ Invoice.bat

log_timestamp "Distribution packages created:"
echo "- dist/Gemini Invoice.app (macOS)"
echo "- dist/Gemini Invoice Windows (Windows)"

# Optional: Create DMG for macOS
if command -v create-dmg &>/dev/null; then
	log_timestamp "Creating DMG for macOS"
	create-dmg \
		--volname "Gemini Invoice" \
		--volicon "icon.icns" \
		--window-pos 200 120 \
		--window-size 600 300 \
		--icon-size 100 \
		--icon "Gemini Invoice.app" 175 120 \
		--hide-extension "Gemini Invoice.app" \
		--app-drop-link 425 120 \
		"dist/Gemini Invoice.dmg" \
		"dist/Gemini Invoice.app"
	log_timestamp "DMG creation completed"
	echo "- dist/Gemini Invoice.dmg (macOS installer)"
else
	log_timestamp "create-dmg not found. Skipping DMG creation for macOS."
fi

# Windows distribution is ready in the folder: dist/Gemini Invoice Windows
log_timestamp "Build script completed"
echo "Note: For Windows, use the folder 'dist/Gemini Invoice Windows'."

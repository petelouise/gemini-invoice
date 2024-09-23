#!/bin/bash

# Create distribution directories
mkdir -p dist/Invoice\ Generator.app/Contents/MacOS
mkdir -p dist/Invoice\ Generator.app/Contents/Resources/fonts
mkdir -p dist/Invoice\ Generator\ Windows/fonts

# Create fonts directory in project root
mkdir -p fonts

# Ensure font files exist and move them to the fonts directory
if [ -f "Inter/Inter Variable/Inter.ttf" ] && [ -f "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" ]; then
    mv "Inter/Inter Variable/Inter.ttf" fonts/Inter.ttf
    mv "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf" fonts/Inter-Bold.ttf
else
    echo "Error: Font files not found. Please ensure they are in the correct location."
    exit 1
fi

# Build the application for macOS
CGO_ENABLED=1 go build -o dist/Invoice\ Generator.app/Contents/MacOS/invoice-generator

# Build the application for Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o dist/Invoice\ Generator\ Windows/invoice-generator.exe

# Create Info.plist for macOS
cat > dist/Invoice\ Generator.app/Contents/Info.plist << EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>invoice-generator</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.yourcompany.invoice-generator</string>
    <key>CFBundleName</key>
    <string>Invoice Generator</string>
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
cp config.yaml dist/Invoice\ Generator.app/Contents/Resources/
cp fonts/Inter.ttf dist/Invoice\ Generator.app/Contents/Resources/fonts/
cp fonts/Inter-Bold.ttf dist/Invoice\ Generator.app/Contents/Resources/fonts/

# Copy necessary files for Windows
cp config.yaml dist/Invoice\ Generator\ Windows/
cp fonts/Inter.ttf dist/Invoice\ Generator\ Windows/fonts/
cp fonts/Inter-Bold.ttf dist/Invoice\ Generator\ Windows/fonts/

# Create a shortcut for Windows
echo '@echo off
start "" "%~dp0invoice-generator.exe"' > dist/Invoice\ Generator\ Windows/Invoice\ Generator.bat

echo "Distribution packages created:"
echo "- dist/Invoice Generator.app (macOS)"
echo "- dist/Invoice Generator Windows (Windows)"

# Optional: Create DMG for macOS
if command -v create-dmg &> /dev/null; then
    create-dmg \
      --volname "Invoice Generator" \
      --volicon "icon.icns" \
      --window-pos 200 120 \
      --window-size 600 300 \
      --icon-size 100 \
      --icon "Invoice Generator.app" 175 120 \
      --hide-extension "Invoice Generator.app" \
      --app-drop-link 425 120 \
      "dist/Invoice Generator.dmg" \
      "dist/Invoice Generator.app"
    echo "- dist/Invoice Generator.dmg (macOS installer)"
else
    echo "Note: create-dmg not found. Skipping DMG creation for macOS."
fi

# Optional: Create installer for Windows
if command -v iscc &> /dev/null; then
    echo "[Setup]
    AppName=Invoice Generator
    AppVersion=1.0
    DefaultDirName={pf}\Invoice Generator
    DefaultGroupName=Invoice Generator
    OutputDir=dist
    OutputBaseFilename=InvoiceGeneratorSetup

    [Files]
    Source: "dist\Invoice Generator Windows\*"; DestDir: "{app}"; Flags: recursesubdirs

    [Icons]
    Name: "{group}\Invoice Generator"; Filename: "{app}\invoice-generator.exe"
    Name: "{commondesktop}\Invoice Generator"; Filename: "{app}\invoice-generator.exe"
    " > installer.iss

    iscc installer.iss
    rm installer.iss
    echo "- dist/InvoiceGeneratorSetup.exe (Windows installer)"
else
    echo "Note: Inno Setup Compiler (iscc) not found. Skipping installer creation for Windows."
fi

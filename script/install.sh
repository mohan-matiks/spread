#!/bin/sh

set -e  # Exit immediately on error

BIN_NAME="spread"   
FILE_NAME="spread"
CDN_URL="https://cdn-swish.justswish.in"
INSTALL_DIR="/usr/local/bin"

# Download the binary
echo "Downloading $BIN_NAME from $CDN_URL/$FILE_NAME..."
curl -L -o "$BIN_NAME" "$CDN_URL/$FILE_NAME"

# Make it executable
chmod +x "$BIN_NAME"

# Move to install directory
sudo mv "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

# Verify installation
if command -v $BIN_NAME >/dev/null 2>&1; then
    echo "$BIN_NAME installed successfully!"
    echo "Run '$BIN_NAME' to start."
else
    echo "Installation failed."
    exit 1
fi
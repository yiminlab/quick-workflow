#!/bin/bash

# Quick Workflow Go Version - Installation Script
# This script installs qk to /usr/local/bin

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

echo -e "${GREEN}Quick Workflow Installer${NC}"
echo "=========================="
echo ""

# Map architecture
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Map OS
case "$OS" in
    darwin)
        PLATFORM="darwin"
        ;;
    linux)
        PLATFORM="linux"
        ;;
    mingw*|msys*|cygwin*)
        PLATFORM="windows"
        ARCH="amd64"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo "Detected platform: ${PLATFORM}-${ARCH}"
echo ""

# GitHub repository
REPO="yiminlab/quick-workflow"
BINARY_NAME="qkflow"

if [ "$PLATFORM" = "windows" ]; then
    BINARY_NAME="qkflow.exe"
fi

# Get latest release
echo "Fetching latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${YELLOW}Could not fetch latest version from GitHub. Using development version.${NC}"
    LATEST_VERSION="latest"
fi

echo "Latest version: $LATEST_VERSION"
echo ""

# Download URL
if [ "$LATEST_VERSION" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${PLATFORM}-${ARCH}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}-${PLATFORM}-${ARCH}"
fi

if [ "$PLATFORM" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

# Download binary
echo "Downloading from: $DOWNLOAD_URL"
TEMP_FILE=$(mktemp)

if ! curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "${RED}Failed to download binary${NC}"
    echo -e "${YELLOW}You can build from source instead:${NC}"
    echo "  cd go-version"
    echo "  make build"
    exit 1
fi

# Make executable
chmod +x "$TEMP_FILE"

# Install
INSTALL_DIR="/usr/local/bin"

if [ "$PLATFORM" = "windows" ]; then
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi

INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "Installing to: $INSTALL_PATH"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_FILE" "$INSTALL_PATH"
else
    echo "Requesting sudo permission to install to $INSTALL_DIR..."
    sudo mv "$TEMP_FILE" "$INSTALL_PATH"
fi

echo -e "${GREEN}✅ Installation successful!${NC}"
echo ""

# Verify installation
if command -v qkflow &> /dev/null; then
    echo "Installed version:"
    qkflow version
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo "  1. Run: ${YELLOW}qkflow init${NC}"
    echo "  2. Configure your credentials"
    echo "  3. Start using: ${YELLOW}qkflow pr create${NC}"
else
    echo -e "${YELLOW}⚠️  qkflow command not found in PATH${NC}"
    echo ""
    echo "Add this to your ~/.zshrc or ~/.bashrc:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo ""
    echo "Then run:"
    echo "  source ~/.zshrc  # or source ~/.bashrc"
fi

echo ""
echo "📚 Documentation: https://github.com/${REPO}"
echo "🐛 Report issues: https://github.com/${REPO}/issues"


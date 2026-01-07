#!/bin/bash
# vaultix installation script for Linux and macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================"
echo "vaultix Installation Script"
echo "======================================"
echo ""

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux*)     OS_TYPE=Linux;;
    Darwin*)    OS_TYPE=macOS;;
    *)          OS_TYPE="UNKNOWN";;
esac

echo "Detected OS: $OS_TYPE"
echo ""

if [ "$OS_TYPE" = "UNKNOWN" ]; then
    echo -e "${RED}Error: Unsupported operating system${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "Go version: $GO_VERSION"
echo ""

# Build the binary
echo "Building vaultix..."
if go build -o vaultix -ldflags="-s -w" .; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
echo ""

# Get the binary size
BINARY_SIZE=$(du -h vaultix | cut -f1)
echo "Binary size: $BINARY_SIZE"
echo ""

# Determine installation directory
INSTALL_DIR="/usr/local/bin"

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    SUDO=""
else
    SUDO="sudo"
    echo -e "${YELLOW}Note: Installation requires sudo privileges${NC}"
fi

# Install the binary
echo "Installing vaultix to $INSTALL_DIR..."
if $SUDO mv vaultix "$INSTALL_DIR/vaultix"; then
    echo -e "${GREEN}✓ Installation successful${NC}"
else
    echo -e "${RED}✗ Installation failed${NC}"
    exit 1
fi

# Ensure it's executable
$SUDO chmod +x "$INSTALL_DIR/vaultix"

echo ""
echo "======================================"
echo -e "${GREEN}Installation Complete!${NC}"
echo "======================================"
echo ""
echo "vaultix is now installed at: $INSTALL_DIR/vaultix"
echo ""
echo "Verify installation:"
echo "  vaultix --help"
echo ""
echo "Get started:"
echo "  vaultix init ./my_vault"
echo ""
echo "For more information, see README.md"

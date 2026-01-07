#!/bin/bash

# Release script for vaultix
# Usage: ./scripts/release.sh v1.0.0

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

echo "Creating release $VERSION..."

# Build for all platforms
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

mkdir -p dist

for platform in "${PLATFORMS[@]}"; do
    OS=$(echo $platform | cut -d'/' -f1)
    ARCH=$(echo $platform | cut -d'/' -f2)
    OUTPUT="dist/vaultix-$OS-$ARCH"
    
    if [ "$OS" = "windows" ]; then
        OUTPUT="$OUTPUT.exe"
    fi
    
    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT
done

echo ""
echo "✅ Builds complete!"
echo ""
echo "Binaries created in dist/:"
ls -lh dist/
echo ""
echo "Next steps:"
echo "1. Test the binaries"
echo "2. Create a git tag: git tag $VERSION && git push origin $VERSION"
echo "3. Create a GitHub release and upload the binaries"

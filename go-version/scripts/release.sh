#!/bin/bash

# Release script for quick-workflow
# Usage: ./scripts/release.sh v1.0.0

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

cd "$(dirname "$0")/.."

echo -e "${GREEN}Creating release $VERSION${NC}"
echo "================================"
echo ""

# 1. Check if git is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Git working directory is not clean${NC}"
    echo "Please commit or stash your changes first"
    exit 1
fi

# 2. Run tests
echo "Running tests..."
if ! ./scripts/test.sh; then
    echo -e "${RED}Tests failed, aborting release${NC}"
    exit 1
fi

# 3. Update version in code (if needed)
echo "Updating version..."
# Add version update logic here if needed

# 4. Build binaries
echo ""
echo "Building binaries..."
make build-all

# 5. Create git tag
echo ""
echo "Creating git tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

# 6. Push tag
echo "Pushing tag to origin..."
git push origin "$VERSION"

echo ""
echo -e "${GREEN}✅ Release $VERSION created!${NC}"
echo ""
echo "Next steps:"
echo "  1. GitHub Actions will build and upload binaries"
echo "  2. Check: https://github.com/yiminlab/quick-workflow/actions"
echo "  3. Edit release notes: https://github.com/yiminlab/quick-workflow/releases"
echo ""
echo "To rollback:"
echo "  git tag -d $VERSION"
echo "  git push origin :refs/tags/$VERSION"


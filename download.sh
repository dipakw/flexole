#!/bin/sh

set -e

# Colors
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
RESET="\033[0m"
NAME="flexole"

info() { echo "${YELLOW}➤ $1${RESET}"; }
success() { echo "${GREEN}✔ $1${RESET}"; }
error() { echo "${RED}✖ $1${RESET}" >&2; }

# Check for required tools
for cmd in curl unzip; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        error "\"$cmd\" is not installed. Please install it and try again."
        exit 1
    fi
done

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    mingw*|cygwin*|msys*) OS="windows" ;;
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) error "Unsupported OS: $OS"; exit 1 ;;
esac

# Determine architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Fetch latest version from GitHub API
REPO="dipakw/${NAME}"
info "Fetching latest version from GitHub..."

VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": "v\([^"]*\)".*/\1/p')

if [ -z "$VERSION" ]; then
    error "Failed to fetch the latest version from GitHub."
    exit 1
fi

FILENAME="${NAME}-${OS}-${ARCH}.zip"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"

info "Detected OS: $OS"
info "Detected Architecture: $ARCH"
info "Latest version: v${VERSION}"
info "Downloading: $URL"

curl -fsSL -o "$FILENAME" "$URL" || { error "Failed to download $FILENAME"; exit 1; }

info "Unzipping: $FILENAME"
unzip -o -q "$FILENAME" || { error "Failed to unzip $FILENAME"; exit 1; }

# Set executable permissions if needed
if [ "$OS" != "windows" ]; then
    chmod +x "${NAME}"
fi

rm -f "$FILENAME"

success "${NAME} v${VERSION} downloaded successfully!"
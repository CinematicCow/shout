#!/bin/bash

# Shout Installation Script
# Downloads and installs the latest version of Shout

set -e

# Configuration
REPO="CinematicCow/shout"
API_URL="https://api.github.com/repos/$REPO/releases/latest"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
detect_os_arch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')
    
    case $OS in
        darwin) OS="macos" ;;
        linux) OS="linux" ;;
        *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
    esac
    
    case $ARCH in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
    esac
    
    echo "${OS}-${ARCH}"
}

# Check if we can write to install directory
check_permissions() {
    if [ ! -w "$INSTALL_DIR" ]; then
        if command -v sudo >/dev/null 2>&1; then
            SUDO="sudo"
        else
            echo "Cannot write to $INSTALL_DIR and sudo is not available" >&2
            exit 1
        fi
    else
        SUDO=""
    fi
}

# Download and install shout
install_shout() {
    local os_arch=$(detect_os_arch)
    local temp_dir=$(mktemp -d)
    local binary_name="shout-${os_arch}"
    
    # Get latest release info
    local release_info=$(curl -s "$API_URL")
    local version=$(echo "$release_info" | grep '"tag_name":' | sed -E 's/.*"tag_name": ?"([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        echo "Failed to fetch version information" >&2
        exit 1
    fi
    
    # Find the download URL for our platform
    local download_url=$(echo "$release_info" | grep -o '"browser_download_url": ?"[^"]*'${binary_name}'[^"]*"' | sed -E 's/.*"browser_download_url": ?"([^"]+)".*/\1/')
    
    if [ -z "$download_url" ]; then
        echo "No binary found for platform: $os_arch" >&2
        exit 1
    fi
    
    # Download and install
    curl -sL "$download_url" -o "$temp_dir/shout"
    chmod +x "$temp_dir/shout"
    
    # Verify the binary works
    if ! "$temp_dir/shout" --version >/dev/null 2>&1; then
        echo "Downloaded binary is not working correctly" >&2
        exit 1
    fi
    
    # Install the binary
    $SUDO mv "$temp_dir/shout" "$INSTALL_DIR/shout"
    
    # Clean up
    rm -rf "$temp_dir"
    
    echo "Shout $version installed successfully"
}

# Check for dependencies
check_dependencies() {
    if ! command -v curl >/dev/null 2>&1; then
        echo "curl is required but not installed" >&2
        exit 1
    fi
}

# Main installation flow
main() {
    check_dependencies
    check_permissions
    
    # Check if shout is already installed
    if command -v shout >/dev/null 2>&1; then
        echo "Shout is already installed. Use 'shout --version' to check the current version."
        exit 0
    fi
    
    install_shout
}

# Run main function
main "$@"

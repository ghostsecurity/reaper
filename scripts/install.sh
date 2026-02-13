#!/bin/bash
# install.sh
# Cross-platform installer for Ghost Security Reaper
# Supports: Linux (amd64, arm64), macOS (amd64, arm64)

set -e

# Configuration
REPO="ghostsecurity/reaper"
BIN_DIR="${HOME}/.ghost/bin"
BINARY_NAME="reaper"

# Detect platform
detect_platform() {
    local os arch

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        *)          echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        aarch64|arm64)  arch="arm64" ;;
        *)              echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
    esac

    echo "${os}_${arch}"
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

# Check if reaper is already installed and get version
get_installed_version() {
    if [ -x "${BIN_DIR}/${BINARY_NAME}" ]; then
        "${BIN_DIR}/${BINARY_NAME}" version 2>/dev/null | head -1 || echo ""
    else
        echo ""
    fi
}

# Extract bare version number for comparison
# "reaper 1.0.0" -> "1.0.0", "v1.0.0" -> "1.0.0"
normalize_version() {
    echo "$1" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1
}

# Download and install from GitHub
install_from_github() {
    local platform="$1"
    local version="$2"
    local os="${platform%_*}"
    local arch="${platform#*_}"
    local download_url

    # Construct download URL
    download_url="https://github.com/${REPO}/releases/download/${version}/reaper_${os}_${arch}.tar.gz"

    echo "Downloading reaper ${version} for ${platform}..."
    echo "URL: ${download_url}"

    # Create bin directory
    mkdir -p "${BIN_DIR}"

    # Download and extract
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT

    curl -sfL "${download_url}" -o "${tmp_dir}/reaper.tar.gz" || return 1
    tar xzf "${tmp_dir}/reaper.tar.gz" -C "${tmp_dir}"

    # Install binary
    mv "${tmp_dir}/${BINARY_NAME}" "${BIN_DIR}/${BINARY_NAME}"
    chmod +x "${BIN_DIR}/${BINARY_NAME}"

    # macOS: remove quarantine attribute
    if [ "$os" = "darwin" ]; then
        xattr -d com.apple.quarantine "${BIN_DIR}/${BINARY_NAME}" 2>/dev/null || true
    fi

    echo "Installed to: ${BIN_DIR}/${BINARY_NAME}"
    return 0
}

# Main
main() {
    echo "Reaper Installer"
    echo "================"

    # Detect platform
    local platform
    platform=$(detect_platform)
    echo "Platform: ${platform}"

    # Get latest version from GitHub API
    local latest_version
    latest_version=$(get_latest_version)

    if [ -z "$latest_version" ]; then
        echo ""
        echo "ERROR: Could not fetch latest version from GitHub."
        echo "Please ensure network access to github.com/ghostsecurity/reaper"
        exit 1
    fi

    echo "Latest version: ${latest_version}"

    # Check if already installed and up to date
    local installed_version
    installed_version=$(get_installed_version)
    echo "Installed version: ${installed_version:-none}"

    if [ -n "$installed_version" ]; then
        local installed_normalized latest_normalized
        installed_normalized=$(normalize_version "$installed_version")
        latest_normalized=$(normalize_version "$latest_version")

        if [ "$installed_normalized" = "$latest_normalized" ]; then
            echo "Already up to date!"
            echo "Binary path: ${BIN_DIR}/${BINARY_NAME}"
            exit 0
        fi

        echo "Updating from ${installed_normalized} to ${latest_normalized}..."
    fi

    if install_from_github "$platform" "$latest_version"; then
        echo ""
        echo "Verification:"
        "${BIN_DIR}/${BINARY_NAME}" version
        echo ""
        echo "Installation complete!"
        echo "Binary path: ${BIN_DIR}/${BINARY_NAME}"
        exit 0
    fi

    echo ""
    echo "ERROR: Could not install reaper."
    echo "Please ensure network access to github.com/ghostsecurity/reaper"
    exit 1
}

main "$@"

#!/usr/bin/env bash
set -euo pipefail

info() {
    echo "info: $1"
}

warn() {
    echo "warn: $1"
}

error() {
    echo "error: $1"
    exit 1
}

install() {
    # allow overriding the version
    VERSION=${SVART_VERSION:-latest}

    REPOSITORY=michielvermeir/svart
    PLATFORM=`uname -s`
    ARCH=`uname -m`

    if [[ $PLATFORM == "Darwin" ]]; then
        PLATFORM="darwin"
    elif [[ $PLATFORM == "Linux" ]]; then
        PLATFORM="linux"
    fi

    case $ARCH in
        armv8*|arm64*|aarch64*)
            ARCH="arm64"
            ;;
        i686*|x86_64*)
            ARCH="amd64"
            ;;
    esac

    BINARY="svart-${PLATFORM}-${ARCH}"

    # Oddly enough GitHub has different URLs for latest vs specific version
    if [[ $VERSION == "latest" ]]; then
        DOWNLOAD_URL=https://github.com/${REPOSITORY}/releases/latest/download/${BINARY}
    else
        DOWNLOAD_URL=https://github.com/${REPOSITORY}/releases/download/${VERSION}/${BINARY}
    fi

    if command -v jq &> /dev/null; then
        version=$(curl -s https://api.github.com/repos/$REPOSITORY/releases/$VERSION | jq -r .name)
        info "this script will automatically download and install svart (${version})"
    else
        info "this script will automatically download and install svart (${VERSION})"
    fi


    if [ "X$(id -u)" == "X0" ]; then
        warn "this script is running as root.  This is dangerous and unnecessary!"
    fi

    if ! command -v 2> /dev/null; then
        error "you do not have 'curl' installed which is required for this script."
    fi

    if ! command -v 2> /dev/null; then
        error "error: you do not have 'gunzip' installed which is required for this script."
    fi

    TEMP_FILE=`mktemp "${TMPDIR:-/tmp}/.svartinstall.XXXXXXXX"`

    cleanup() {
        rm -f "$TEMP_FILE"
        info "svart installed successfully"
        svart --version
    }

    trap cleanup EXIT

    echo "info: downloading $DOWNLOAD_URL"
    HTTP_CODE=$(curl --progress-bar -SL "$DOWNLOAD_URL" --output "$TEMP_FILE" --write-out "%{http_code}")
    if [[ ${HTTP_CODE} -lt 200 || ${HTTP_CODE} -gt 299 ]]; then
        error "platform ${PLATFORM} (${ARCH}) is unsupported."
    fi

    chmod +x "$TEMP_FILE"

    # Detect when the file cannot be executed due to NOEXEC /tmp.  Taken from rustup
    # https://github.com/rust-lang/rustup/blob/87fa15d13e3778733d5d66058e5de4309c27317b/rustup-init.sh#L158-L159
    if [ ! -x "$TEMP_FILE" ]; then
        printf '%s\n' "Cannot execute $TEMP_FILE (likely because of mounting /tmp as noexec)." 1>&2
        printf '%s\n' "Please copy the file to a location where you can execute binaries and run it manually." 1>&2
        exit 1
    fi

    info "installing svart to /usr/local/bin"
    mv "$TEMP_FILE" /usr/local/bin/svart || sudo mv "$TEMP_FILE" /usr/local/bin/svart
}

install
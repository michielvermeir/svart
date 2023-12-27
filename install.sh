#!/usr/bin/env bash
set -euo pipefail

install() {
    # allow overriding the version
    VERSION=${TFVARS_VERSION:-latest}

    REPOSITORY=michielvermeir/tfvars
    PLATFORM=`uname -s`
    ARCH=`uname -m`

    if [[ $PLATFORM == "Darwin" ]]; then
        PLATFORM="darwin"
    elif [[ $PLATFORM == "Linux" ]]; then
        PLATFORM="linux"
    fi

    if [[ $ARCH == armv8* ]] || [[ $ARCH == arm64* ]] || [[ $ARCH == aarch64* ]]; then
        ARCH="arm64"
    elif [[ $ARCH == i686* ]]; then
        ARCH="amd64"
    fi

    BINARY="tfvars-${PLATFORM}-${ARCH}"

    # Oddly enough GitHub has different URLs for latest vs specific version
    if [[ $VERSION == "latest" ]]; then
        DOWNLOAD_URL=https://github.com/${REPOSITORY}/releases/latest/download/${BINARY}
    else
        DOWNLOAD_URL=https://github.com/${REPOSITORY}/releases/download/${VERSION}/${BINARY}
    fi

    echo "info: script will automatically download and install tfvars (${VERSION}) for you."

    if [ "X$(id -u)" == "X0" ]; then
        echo "warning: this script is running as root.  This is dangerous and unnecessary!"
    fi

    if ! hash curl 2> /dev/null; then
        echo "error: you do not have 'curl' installed which is required for this script."
        exit 1
    fi

    if ! hash gunzip 2> /dev/null; then
        echo "error: you do not have 'gunzip' installed which is required for this script."
        exit 1
    fi

    TEMP_FILE=`mktemp "${TMPDIR:-/tmp}/.tfvarsinstall.XXXXXXXX"`

    cleanup() {
        rm -f "$TEMP_FILE"
    }

    trap cleanup EXIT
    HTTP_CODE=$(curl -SL "$DOWNLOAD_URL" --output "$TEMP_FILE" --write-out "%{http_code}")
    if [[ ${HTTP_CODE} -lt 200 || ${HTTP_CODE} -gt 299 ]]; then
        echo "error: platform ${PLATFORM} (${ARCH}) is unsupported."
        exit 1
    fi

    chmod +x "$TEMP_FILE"

    # Detect when the file cannot be executed due to NOEXEC /tmp.  Taken from rustup
    # https://github.com/rust-lang/rustup/blob/87fa15d13e3778733d5d66058e5de4309c27317b/rustup-init.sh#L158-L159
    if [ ! -x "$TEMP_FILE" ]; then
        printf '%s\n' "Cannot execute $TEMP_FILE (likely because of mounting /tmp as noexec)." 1>&2
        printf '%s\n' "Please copy the file to a location where you can execute binaries and run it manually." 1>&2
        exit 1
    fi

    echo "info: installing tfvars to /usr/local/bin"
    exec sudo mv "$TEMP_FILE" /usr/local/bin/tfvars
}

install
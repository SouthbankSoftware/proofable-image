#!/bin/bash
set -eu

PLATFORM=$(uname)
INSTALL_PATH=${INSTALL_PATH:="$PWD"}
VERSION=${VERSION:="v0.1.1"}

OPTS=""

if [[ "$PLATFORM" = "Darwin" ]]; then
    DOWNLOAD_LINK="https://github.com/SouthbankSoftware/proofable-image/releases/download/${VERSION}/proofable-image_darwin_amd64.tar.gz"
elif [[ "$PLATFORM" = "Linux" ]]; then
    DOWNLOAD_LINK="https://github.com/SouthbankSoftware/proofable-image/releases/download/${VERSION}/proofable-image_linux_amd64.tar.gz"
    OPTS="--overwrite"
else
    echo "unsupported platform \`$PLATFORM\`, please try to build from source: https://github.com/SouthbankSoftware/proofable-image#build-your-own-binary"
    exit 1
fi

echo -e "Installing from \`$DOWNLOAD_LINK\` to \`$INSTALL_PATH\`...\n"

if [[ $(command -v curl) ]]; then
    DOWNLOAD_CMD="curl -L \"$DOWNLOAD_LINK\""
elif [[ $(command -v wget) ]]; then
    DOWNLOAD_CMD="wget -O- \"$DOWNLOAD_LINK\""
else
    echo "neither \`curl\` nor \`wget\` is installed, please download and install the binary manually: https://github.com/SouthbankSoftware/proofable-image/releases"
    exit 1
fi

mkdir -p "$INSTALL_PATH"
eval "$DOWNLOAD_CMD" | tar -zxvC "$INSTALL_PATH" --no-same-owner $OPTS

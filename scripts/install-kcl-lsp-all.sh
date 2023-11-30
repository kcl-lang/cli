#!/usr/bin/env bash

# ------------------------------------------------------------
# Copyright The KCL Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Reference: https://github.com/dapr/cli/tree/master/install
# ------------------------------------------------------------

# sudo is required to copy binary to KCL_INSTALL_DIR for linux
: ${USE_SUDO:="false"}

# Http request CLI
KCL_HTTP_REQUEST_CLI=curl

# GitHub Organization and repo name to download release
GITHUB_ORG=kcl-lang
GITHUB_REPO=kcl

# KCL languge server filename
CLI_FILENAME=kcl-language-server

# --- helper functions for logs ---
info() {
    local action="$1"
    local details="$2"
    command printf '\033[1;32m%12s\033[0m %s\n' "$action" "$details" 1>&2
}

warn() {
    command printf '\033[1;33mWarn\033[0m: %s\n' "$1" 1>&2
}

error() {
    command printf '\033[1;31mError\033[0m: %s\n' "$1" 1>&2
}

request() {
    command printf '\033[1m%s\033[0m\n' "$1" 1>&2
}

eprintf() {
    command printf '%s\n' "$1" 1>&2
}

bold() {
    command printf '\033[1m%s\033[0m' "$1"
}

# If file exists, echo it
echo_fexists() {
    [ -f "$1" ] && echo "$1"
}

runAsRoot() {
    local CMD="$*"

    if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
        CMD="sudo $CMD"
    fi

    $CMD
}

checkHttpRequestCLI() {
    if type "curl" > /dev/null; then
        KCL_HTTP_REQUEST_CLI=curl
    elif type "wget" > /dev/null; then
        KCL_HTTP_REQUEST_CLI=wget
    else
        error "Either curl or wget is required"
        exit 1
    fi
}

getLatestRelease() {
    local KCLReleaseUrl="https://api.github.com/repos/${GITHUB_ORG}/${GITHUB_REPO}/releases"
    local latest_release=""

    if [ "$KCL_HTTP_REQUEST_CLI" == "curl" ]; then
        latest_release=$(curl -s $KCLReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    else
        latest_release=$(wget -q --header="Accept: application/json" -O - $KCLReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    fi

    ret_val=$latest_release
}

downloadFile() {
    LATEST_RELEASE_TAG=$1
    OS=$2
    ARCH=$3
    KCL_CLI_UNIX_ARTIFACT="kclvm-${LATEST_RELEASE_TAG}-${OS}-${ARCH}.tar.gz"
    KCL_CLI_WINDOWS_ARTIFACT="kclvm-${LATEST_RELEASE_TAG}-${OS}.zip"
    # Unix tar.gz artifact
    KCL_CLI_ARTIFACT=$KCL_CLI_UNIX_ARTIFACT
    # Windows zip artifact
    if [ "$OS" == "windows" ]; then
        KCL_CLI_ARTIFACT=$KCL_CLI_WINDOWS_ARTIFACT
    fi
    DOWNLOAD_BASE="https://github.com/${GITHUB_ORG}/${GITHUB_REPO}/releases/download"
    DOWNLOAD_URL="${DOWNLOAD_BASE}/${LATEST_RELEASE_TAG}/${KCL_CLI_ARTIFACT}"

    # Create the temp directory
    KCL_TMP_ROOT=$(mktemp -dt kcl-lsp-install-XXXXXX)
    ARTIFACT_TMP_FILE="$KCL_TMP_ROOT/$KCL_CLI_ARTIFACT"

    info "Downloading $DOWNLOAD_URL ..."
    if [ "$KCL_HTTP_REQUEST_CLI" == "curl" ]; then
        curl -SsL "$DOWNLOAD_URL" -o "$ARTIFACT_TMP_FILE"
    else
        wget -q -O "$ARTIFACT_TMP_FILE" "$DOWNLOAD_URL"
    fi

    if [ ! -f "$ARTIFACT_TMP_FILE" ]; then
        error "Failed to download $DOWNLOAD_URL ..."
        exit 1
    else
        info "Scucessful to download $DOWNLOAD_URL"
    fi

    info "Build kcl language server artifact..."

    INSTALL_FOLDER="./kcl-lsp-${OS}-${ARCH}"
    mkdir -p "./bin/$INSTALL_FOLDER"

    tar xf $ARTIFACT_TMP_FILE -C $KCL_TMP_ROOT
    local tmp_kclvm_folder=$KCL_TMP_ROOT/kclvm

    if [ ! -f "$tmp_kclvm_folder/bin/kcl-language-server" ]; then
        error "Failed to unpack KCL language server executable."
        exit 1
    fi

    # Copy kcl-languge-server in the temp folder into the target installation directory.
    info "Copy the kcl language server binary $tmp_kclvm_folder/bin/kcl-language-server into the target installation directory ./bin"
    cp -f $tmp_kclvm_folder/bin/kcl-language-server "./bin/$INSTALL_FOLDER"
    cd ./bin
    if [ "$OS" == "windows" ]; then
        TARBALL="./kcl-lsp-${LATEST_RELEASE_TAG}-${OS}-${ARCH}.zip"
        info "Zip $TARBALL..."
        zip -r $TARBALL $INSTALL_FOLDER
    else
        TARBALL="./kcl-lsp-${LATEST_RELEASE_TAG}-${OS}-${ARCH}.tar.gz"
        info "Tar $TARBALL..."
        tar -zcf $TARBALL $INSTALL_FOLDER
    fi
    cd ..
    info "Build kcl language server artifact successful!"
}

fail_trap() {
    result=$?
    if [ "$result" != "0" ]; then
        error "Failed to install KCL language server"
        info "For support, go to https://kcl-lang.io"
    fi
    cleanup
    exit $result
}

cleanup() {
    if [[ -d "${KCL_TMP_ROOT:-}" ]]; then
        rm -rf "$KCL_TMP_ROOT"
    fi
}

# -----------------------------------------------------------------------------
# main
# -----------------------------------------------------------------------------
trap "fail_trap" EXIT

checkHttpRequestCLI

if [ -z "$1" ]; then
    echo "Getting the latest KCL language server ..."
    getLatestRelease
else
    ret_val=v$1
fi

info "Find the latest KCL language server version $ret_val"

downloadFile $ret_val "darwin" "amd64"
downloadFile $ret_val "darwin" "arm64"
downloadFile $ret_val "linux" "amd64"
downloadFile $ret_val "windows" "amd64"
cleanup

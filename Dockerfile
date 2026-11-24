ARG BUILDER_IMAGE=golang:1.26
ARG BASE_IMAGE=debian:12-slim

# Build the manager binary
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder

COPY / /src
WORKDIR /src

# The TARGETOS and TARGETARCH args are set by docker. We set GOOS and GOARCH to
# these values to ask Go to compile a binary for these architectures. If
# TARGETOS and TARGETOS are different from BUILDPLATFORM, Go will cross compile
# for us (e.g. compile a linux/amd64 binary on a linux/arm64 build machine).
ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build  GOARCH=${TARGETARCH}  GOOS=linux CGO_ENABLED=0 make build

FROM  ${BASE_IMAGE}

COPY --from=builder /src/bin/kcl /usr/local/bin/kcl

# Verify KCL installation and basic functionality
RUN kcl version && \
    echo 'a=1' | kcl run -

# The minimum libbrotli1 version that carries the CVE-2026-111497 fix. Debian
# backports the fix into bookworm-security as 1.0.9-2+deb12u1 (DLA-4855-1) and
# into trixie-security as 1.1.0-2+deb13u1 (DSA-6582-1).
ARG LIBBROTLI_MIN_VERSION=1.0.9-2+deb12u1

# Install git for KCL package management
# Use best practices for apt-get commands
# The upgrade keeps the base image on the latest security updates, so the image
# is never shipped with stale packages.
RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates ssh && \
    apt-get upgrade -y && \
    rm -rf /var/lib/apt/lists/*

# git depends on libcurl3-gnutls, which in turn depends on libbrotli1. The Brotli
# decoder of libbrotli1 is vulnerable to CVE-2026-111497 unless the package is
# taken from the security archive. Fail the build instead of publishing an image
# that ships a vulnerable decoder.
RUN set -eu; \
    brotli_version="$(dpkg-query -W -f='${Version}' libbrotli1 2>/dev/null || true)"; \
    if [ -n "$brotli_version" ] && ! dpkg --compare-versions "$brotli_version" ge "$LIBBROTLI_MIN_VERSION"; then \
        echo "libbrotli1 $brotli_version is affected by CVE-2026-111497, need >= $LIBBROTLI_MIN_VERSION" >&2; \
        exit 1; \
    fi; \
    echo "libbrotli1 version: ${brotli_version:-not installed}"

# Create the temporary directory
RUN mkdir -p /tmp

# Configure KCL runtime environment
# Set temporary directories for write permissions
ENV KCL_LIB_HOME=/tmp \
    KCL_PKG_PATH=/tmp \
    KCL_CACHE_PATH=/tmp \
    LANG=en_US.utf8

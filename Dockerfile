FROM --platform=${BUILDPLATFORM} golang:1.23 AS build
COPY / /src
WORKDIR /src

# The TARGETOS and TARGETARCH args are set by docker. We set GOOS and GOARCH to
# these values to ask Go to compile a binary for these architectures. If
# TARGETOS and TARGETOS are different from BUILDPLATFORM, Go will cross compile
# for us (e.g. compile a linux/amd64 binary on a linux/arm64 build machine).
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build

FROM debian:11-slim AS image

COPY --from=build /src/bin/kcl /usr/local/bin/kcl
# Verify KCL installation and basic functionality
RUN kcl version && \
    echo 'a=1' | kcl run -

# Install git for KCL package management
# Use best practices for apt-get commands
RUN apt-get update && \
    apt-get install -y --no-install-recommends git && \
    rm -rf /var/lib/apt/lists/*

# Configure KCL runtime environment
# Set temporary directories for write permissions
ENV KCL_LIB_HOME=/tmp \
    KCL_PKG_PATH=/tmp \
    KCL_CACHE_PATH=/tmp \
    LANG=en_US.utf8

# Switch to non-root user for security
USER nonroot:nonroot

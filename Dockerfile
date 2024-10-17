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

FROM --platform=${BUILDPLATFORM} ubuntu:22.04 AS base
ENV LANG=en_US.utf8

FROM base

ARG TARGETARCH

COPY --from=build /src/bin/kcl /usr/local/bin/kcl
RUN /usr/local/bin/kcl
RUN apt-get update && apt-get install make gcc git -y && rm -rf /var/lib/apt/lists/*
# The reason for doing this below is to prevent the
# container from not having write permissions.
ENV KCL_LIB_HOME=/tmp
ENV KCL_PKG_PATH=/tmp
ENV KCL_CACHE_PATH=/tmp
# Install the tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-${TARGETARCH} /tini
RUN chmod +x /tini

FROM --platform=${BUILDPLATFORM} golang:1.23 AS build
COPY / /src
WORKDIR /src

# The TARGETOS and TARGETARCH args are set by docker. We set GOOS and GOARCH to
# these values to ask Go to compile a binary for these architectures. If
# TARGETOS and TARGETOS are different from BUILDPLATFORM, Go will cross compile
# for us (e.g. compile a linux/amd64 binary on a linux/arm64 build machine).
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED 0

RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build

FROM gcr.io/distroless/base-debian11 AS image

COPY --from=build /src/bin/kcl /usr/local/bin/kcl
# Show KCL version
RUN /src/bin/kcl version
# Enable kcl works fine
RUN echo 'a=1' | kcl run -
# Install Git Dependency
RUN apt-get update && apt-get install git -y && rm -rf /var/lib/apt/lists/*
# The reason for doing this below is to prevent the
# container from not having write permissions.
ENV KCL_LIB_HOME /tmp
ENV KCL_PKG_PATH /tmp
ENV KCL_CACHE_PATH /tmp
ENV LANG en_US.utf8
USER nonroot:nonroot

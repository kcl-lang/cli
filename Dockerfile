FROM --platform=${BUILDPLATFORM} golang:1.21 AS build
COPY / /src
WORKDIR /src

# The TARGETOS and TARGETARCH args are set by docker. We set GOOS and GOARCH to
# these values to ask Go to compile a binary for these architectures. If
# TARGETOS and TARGETOS are different from BUILDPLATFORM, Go will cross compile
# for us (e.g. compile a linux/amd64 binary on a linux/arm64 build machine).
ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build

FROM --platform=${BUILDPLATFORM} ubuntu:22.04 AS base
ENV LANG=en_US.utf8

FROM base
COPY --from=build /src/bin/kcl /usr/local/bin/kcl
RUN /usr/local/bin/kcl
RUN cp -r /root/go/bin/* /usr/local/bin/
RUN apt-get update
RUN apt-get install gcc git -y
# The reason for doing this below is to prevent the
# container from not having write permissions.
ENV KCL_PKG_PATH=/tmp
ENV KCL_CACHE_PATH=/tmp
# In the image, we can generate a runtime in advance to
# avoid writing files in the image
ENV KCL_GO_DISABLE_INSTALL_ARTIFACT=true
ENV KCL_GO_DISABLE_ARTIFACT_IN_PATH=false
# Install the tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

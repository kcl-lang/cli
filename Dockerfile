FROM golang:1.21 AS build
COPY / /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build make build

FROM ubuntu:22.04 AS base
ENV LANG=en_US.utf8

FROM base
COPY --from=build /src/bin/kcl /usr/local/bin/kcl
RUN /usr/local/bin/kcl
RUN cp -r /root/go/bin/* /usr/local/bin/
RUN apt-get update
RUN apt-get install gcc -y
# The reason for doing this below is to prevent the
# container from not having write permissions.
ENV KCL_GO_DISABLE_ARTIFACT=on
ENV KCL_PKG_PATH=/tmp
ENV KCL_CACHE_PATH=/tmp

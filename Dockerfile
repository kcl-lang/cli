FROM golang:1.19 AS build
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
ENV KCL_GO_DISABLE_ARTIFACT=on

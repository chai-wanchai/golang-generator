FROM golang:1.17-alpine3.14 AS builder

# Define current working directory
WORKDIR /app

RUN apk update; \
    apk add --no-cache \
    git \
    make
RUN go env -w GOSUMDB=off

# turn on go module
ENV GO111MODULE on

# Setup build environments
# CGO_ENABLED=0
# for cross-compilation is because we should use the go
# built-in support of the target platform cross-compilation
# and there should be no reason not to do so.
# GOOS=linux
# use the linux operation system for build
# - GOARCH=arm64
# use the arm64 as the CPU architecture
ENV CGO_ENABLED=0
ENV GOOS=linux
# ENV GOARCH=arm64

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.mod .
COPY go.sum .
COPY . ./
RUN go mod vendor
RUN go install github.com/wandercn/hotbuild@latest
RUN ls -a -l
ENTRYPOINT ["hotbuild","run"]
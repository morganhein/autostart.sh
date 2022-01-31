# COMPILE
FROM golang:1.17 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.2

ADD ./ /build
WORKDIR /build

# RUN
FROM ubuntu:hirsute

EXPOSE 40000

WORKDIR /shoelace
COPY --from=build-env /go/bin/dlv /
COPY ./build/shoelace /shoelace/shoelace
ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/shoelace/shoelace", "--"]
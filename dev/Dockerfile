# COMPILE
FROM golang:1.17 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.2

ADD ../dev /build
WORKDIR /build


RUN env CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /shoelace main.go

# RUN
FROM ubuntu:hirsute

EXPOSE 40000

WORKDIR /shoelace
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /shoelace /shoelace/shoelace
ADD examples/ /shoelace/config/
ADD dev/start.sh /shoelace/start.sh
COPY --from=build-env /go/bin/dlv /
ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/shoelace/shoelace", "--"]
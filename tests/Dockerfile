# COMPILE
FROM golang:1.17 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.2

ADD ./ /build
WORKDIR /build


RUN env CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /autostart main.go

# RUN
FROM ubuntu:hirsute

EXPOSE 40000

WORKDIR /autostart
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /autostart /autostart/autostart
ADD examples/ /autostart/config/
ADD tests/start.sh /autostart/start.sh
COPY --from=build-env /go/bin/dlv /
ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/autostart/autostart", "--"]
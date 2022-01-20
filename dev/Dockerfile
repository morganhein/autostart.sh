# COMPILE
FROM golang:1.17 AS build-env

WORKDIR /build
RUN env CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /shoelace main.go

# RUN
FROM ubuntu:hirsute

EXPOSE 40000

WORKDIR /shoelace
COPY --from=build-env /shoelace /shoelace/shoelace
ADD examples/ /shoelace/config/
ADD dev/start.sh /shoelace/start.sh
ENTRYPOINT ["exec", "/shoelace/shoelace", "--"]
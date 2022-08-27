FROM golang:bullseye
# need to pass -v for a mountpoint when running, mounting the base project

EXPOSE 40000

WORKDIR /
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

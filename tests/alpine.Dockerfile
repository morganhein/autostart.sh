FROM golang:alpine

EXPOSE 40000

WORKDIR /
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN apk update && apk add zsh -y

WORKDIR /app
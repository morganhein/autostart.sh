FROM golang:1.17.3-buster

RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.2

EXPOSE 40000
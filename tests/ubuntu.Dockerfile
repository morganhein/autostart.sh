FROM ubuntu:hirsute

RUN apt-get update -y && apt-get install gcc make golang ca-certificates -y
RUN go version

RUN mkdir /build
WORKDIR /build
ADD go.mod /build/
ADD go.sum /build/
RUN go mod download

RUN mkdir go
WORKDIR /go

ENTRYPOINT ["go", "test"]
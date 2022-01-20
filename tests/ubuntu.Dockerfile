FROM ubuntu:hirsute

RUN apt-get update -y && apt-get install gcc make golang ca-certificates -y
RUN go version
RUN mkdir /go
WORKDIR /go

#ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out", "--tags", "ubuntu"]

ENTRYPOINT ["/bin/bash"]
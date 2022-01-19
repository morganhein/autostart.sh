#!/usr/bin/env bash

# first build
echo "building binary"
env CGO_ENABLED=0 GOOS=linux go build -o ./../build/shoelace ./../main.go
chmod +x ./../build/shoelace

echo "starting docker"
# start the docker container, with this file mounted, and the directory where the config is coming from
docker run --rm -v ${PWD}/../examples:/data -v ${PWD}/../build/:/shoelace --name shoelace -it -w /shoelace ubuntu:hirsute bash
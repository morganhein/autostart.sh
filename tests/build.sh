#!/usr/bin/env bash

# first build
echo "building binary"
env CGO_ENABLED=0 GOOS=linux go build -o ./../build/autostart ./../main.go
chmod +x ./../build/autostart

echo "starting docker"
# start the docker container, with this file mounted, and the directory where the config is coming from
docker run --rm -v ${PWD}/../examples:/data -v ${PWD}/../build/:/autostart --name autostart -it -w /autostart ubuntu:hirsute bash
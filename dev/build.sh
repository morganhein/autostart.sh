#!/usr/bin/env bash

# first build
#echo "building binary"
#env CGO_ENABLED=0 GOOS=linux go build -o ./../build/shoelace ./../main.go
#chmod +x ./../build/shoelace


# TODO: get the absolute path to the above directory and set it as the build path
cd ..

# build dockerfile
docker build -t shoelace:latest -f dev/Dockerfile .

echo "starting docker"
# start the docker container, with this file mounted, and the directory where the config is coming from
docker run --rm -v ${PWD}/configs/examples:/data -v ${PWD}/build/:/shoelace -p 40000:40000 --name shoelace -w /shoelace shoelace:latest

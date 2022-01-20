#!/usr/bin/env bash

# build the image
docker build -t ubuntu-shoelace:latest -f ubuntu.Dockerfile .

# run it
docker run --rm -v ${PWD}/../:/go ubuntu-shoelace:latest
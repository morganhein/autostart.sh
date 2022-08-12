#!/usr/bin/env bash

# get script location, and make sure WORKDIR=PWD/../
SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
cd "${SCRIPTPATH}"/.. || exit
echo "Working in ${PWD}"

# first build
#echo "building binary"
env CGO_ENABLED=0 GOOS=linux go build -o ./build/envy ./main.go
chmod +x ./../build/envy

# build dockerfile
docker build -t envy:latest -f dev/Dockerfile .

echo "starting docker"
# start the docker container, with this file mounted, and the directory where the config is coming from
docker run --rm -v ${PWD}/configs/examples:/data -v ${PWD}/build/:/envy -p 40000:40000 -w /envy envy:latest task gvm --config=/data/personal.toml
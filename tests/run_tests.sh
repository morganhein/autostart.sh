#!/usr/bin/env bash

DISTROS=("debian" "alpine")

for i in "${DISTROS[@]}"
do
   :
   # build it
   echo "Building ${i}"
#   docker build -t envy-"${i}":latest -f tests/"${i}".Dockerfile ./..
   # run it
   echo "Running tests on ${i}"
   docker run --rm -e CGO_ENABLED=0 -v $PWD:/app envy-"${i}" go test -v tests/*.go
done

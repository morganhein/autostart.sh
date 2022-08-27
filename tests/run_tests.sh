#!/usr/bin/env bash

TESTS=( $(grep -Phro "func \K(Test[a-zA-Z0-9\-_]+)" tests/*_test.go) )
DISTROS=("debian" "alpine")

for distro in "${DISTROS[@]}"
do
   :
   # build it
   echo "Building ${distro}"
   #docker build -t envy-"${distro}":latest -f tests/"${distro}".Dockerfile ./..
   echo "Running on ${distro}"
   for test in "${TESTS[@]}"
   do
     # run it
      echo "Running test ${test}"
     docker run --rm -e CGO_ENABLED=0 -v $PWD:/app envy-"${distro}" go test -v --tags=integrated -run "${test}" ./tests/
   done
done

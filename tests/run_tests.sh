#!/usr/bin/env bash

# UBUNTU first

# build the image
docker build -t ubuntu-envy:latest -f ubuntu.Dockerfile ./..

# list all the tests
TESTS=( $(grep -Phro "func \K(Test[a-zA-Z0-9\-_]+)" ubuntu_test.go) )
TAGS="ubuntu"

echo "
Container built, running tests!
"
for i in "${TESTS[@]}"
do
   :
   # run it
   echo "Running test ${i}"
   # need to pass args for `go test` including tags and the test to run
   #echo "docker run --rm -v ${PWD}/../:/go ubuntu-envy:latest --tags=${TAGS} -run ${i} tests/"
   docker run --rm -v "${PWD}"/../:/go ubuntu-envy:latest --tags=${TAGS} -run "${i}" ./tests/
done
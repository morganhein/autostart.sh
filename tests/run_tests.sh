#!/usr/bin/env bash

function help(){
    echo "USAGE: run_tests.sh -B, -D, -o <distro/OS>, -t <testname>, -h"
}

BUILD=false
DELVE=false
DISTROS=("debian" "alpine")
TESTS=( $(grep -Phro "func \K(Test[a-zA-Z0-9\-_]+)" tests/*_test.go) )

# When using dlv, a specific test and os must be specified
# When a specific test is specified, an os must be specified
function checkArgs() {
  echo "Delve: ${DELVE}, Distro/OS: ${DISTRO}, Test: ${TEST}, Build: ${BUILD}"
  # if test was passed, then we also need the -o argument to specify which distro
  if [ "$DELVE" = true ] && [[ -z "$DISTRO" ]]; then
       echo "The \"-o <DISTRO/OS>\" argument is required when using delve."
       exit 1
  fi
  if [[ -n "$TEST" ]] && [[ -z "$DISTRO" ]]; then
    echo "The \"-o <DISTRO/OS>\" argument is required when targeting a specific test."
    exit 1
  fi
  if [[ "$DELVE" = true ]] && [[ -z "$TEST" ]]; then
    echo "The \"-t <testname>\" argument is required when using delve."
    exit 1
  fi
}

function build(){
  if [ "$BUILD" = false ]; then
    return
  fi
  echo "Building ${DISTRO}"
  docker build -t envy-"${DISTRO}":latest -f tests/"${DISTRO}".Dockerfile ./..
}

function run(){
  docker run --rm -e CGO_ENABLED=0 -v $PWD:/app envy-"${DISTRO}" go test -v --tags=integrated -run "${TEST}" ./tests/
}

function runWithDelve(){
  echo "Starting delve session for test: \"${TEST}\""
  docker run --rm -e CGO_ENABLED=0 --security-opt="apparmor=unconfined" --cap-add=SYS_PTRACE -p 40000:40000 -v $PWD:/app -it envy-"${DISTRO}" dlv test --listen=:40000 --headless=true --api-version=2 --accept-multiclient ./tests/*.go -- -test.run ^"${TEST}"$
}

function runAll(){
  # if we've requested a specific distro, set the list to that
  if [[ -n $DISTRO ]]; then
    DISTROS=("$DISTRO")
  fi

  echo "Running all tests with the following target(s): ${DISTROS[*]}!"

  for distro in "${DISTROS[@]}"
  do
     :
     DISTRO="${distro}"
     build
     echo "Running on ${distro}"
     for test in "${TESTS[@]}"
     do
       TEST="${test}"
       run
     done
  done
}

while getopts ":BDt:o:h" ARG; do
  case "$ARG" in
    B) BUILD=true;;
    t) TEST=$OPTARG;;
    D) DELVE=true;;
    o) DISTRO=$OPTARG;;
    h) help ;;
    :) echo "argument missing" ;;
    \?) echo "Something is wrong" ;;
  esac
done

shift "$((OPTIND-1))"

# Sanity check arguments
checkArgs

# If we are running a delve session, start it
if [ "$DELVE" = true ]; then
  build
  runWithDelve
  exit $?
fi

# If we are just running a specific test
if [[ -n "$TEST" ]]; then
  build
  run
  exit $?
fi

# Otherwise run all tests on target(s)
runAll
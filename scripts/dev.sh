#!/usr/bin/env bash

function help(){
    echo "USAGE: dev.sh -B, -o <distro/OS>, -h"
}

while getopts ":Bo:h" ARG; do
  case "$ARG" in
    B) BUILD=true;;
    o) DISTRO=$OPTARG;;
    h) help ;;
    :) echo "argument missing" ;;
    \?) echo "Something is wrong" ;;
  esac
done

# TODO
# Spawn a dlv session with passed parameters passed to application
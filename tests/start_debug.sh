#!/usr/bin/env bash

docker run --rm -it --entrypoint=bash -p 40000:40000 -w /data/tests -v ${PWD}/../:/data shoelace_deb

/usr/bin/bash -c "which apt"

error running command `/usr/bin/bash -c "which brew"`
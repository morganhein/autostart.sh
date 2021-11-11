#!/usr/bin/env bash

env CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /autostart ../main.go

dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /autostart -- --config=basic.toml --sudo=false -v golang
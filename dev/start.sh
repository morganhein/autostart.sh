#!/usr/bin/env bash
cd ..
env CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /envy main.go

dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /envy -- --config=basic.toml --sudo=false -v golang
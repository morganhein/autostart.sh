#!/usr/bin/env bash

cd /app

if (($# > 0))
then
  echo "running tests with dlv"
  dlv test --listen=:40000 --headless=true --api-version=2 --accept-multiclient pkg/io/*.go
else
  echo "running tests without dlv"
  go test ./...
fi
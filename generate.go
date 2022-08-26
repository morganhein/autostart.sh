package main

//go:generate rm -f pkg/out/runner_mock.go
//go:generate rm -f pkg/io/filesystem_mock.go

//go:generate moq -out pkg/io/runner_mock.go pkg/io Runner
//go:generate moq -out pkg/io/filesystem_mock.go pkg/io Filesystem

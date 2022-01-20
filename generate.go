package main

//go:generate moq -out pkg/io/runner_mock.go pkg/io Runner
//go:generate moq -out pkg/io/filesystem_mock.go pkg/io Filesystem

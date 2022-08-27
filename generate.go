package main

//go:generate rm -f pkg/io/shell_mock.go
//go:generate rm -f pkg/io/filesystem_mock.go

//go:generate moq -out pkg/io/shell_mock.go pkg/io Shell
//go:generate moq -out pkg/io/filesystem_mock.go pkg/io Filesystem

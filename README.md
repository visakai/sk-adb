#goadb

[![Build Status](https://travis-ci.org/zach-klippenstein/goadb.svg?branch=master)](https://travis-ci.org/zach-klippenstein/goadb)
[![GoDoc](https://godoc.org/github.com/zach-klippenstein/goadb?status.svg)](https://godoc.org/github.com/zach-klippenstein/goadb)

A Golang library for interacting with the Android Debug Bridge (adb).

See [demo.go](cmd/demo/demo.go) for usage.


For this project, a tool called `stringer` is used to modify some files during `go generate` step. 
Codes will be generated in the following files:
```
devicedescriptortype_string.go
devicestate_string.go
internal/errors/errcode_string.go
```
Please make sure your environment variable `GOBIN` is correctly set, so that `stringer` can be successfully installed.

There is a Makefile at repo root. Just need to run:
> make

It will download all dependencies, install `stringer` locally, generate code, and run the tests.


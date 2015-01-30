#!/bin/bash

export GOPATH=$(pwd)

go install networkmodule

go run src/main.go

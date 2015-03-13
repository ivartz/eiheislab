#!/bin/bash

export GOPATH=$(pwd)

go install network

go run src/main.go -raddr="129.241.187.151:20017" -lport = 20018

#!/bin/bash

#export GOPATH=$HOME/Documents
#test
export GOPATH=$(pwd)

#go install network

go install driver queue states communication network

#go install driver

#go run src/main.go -raddr="129.241.187.153:20021" -lport=20017

#go build src/main.go
go run src/main.go

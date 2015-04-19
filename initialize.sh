#!/bin/bash

#export GOPATH=$HOME/Documents
#test
export GOPATH=$(pwd)

#go install network

go install driver queue states

#go install driver

#go run src/main.go -raddr="129.241.187.153:20021" -lport=20017

go run src/main.go

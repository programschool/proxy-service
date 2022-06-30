#!/bin/bash

cp ./conf.d/config-local.json config.json

export GOOS=linux
export GOARCH=amd64

go mod vendor
go fmt main.go
go build -o middle-proxy main.go

build="docker build . -f Dockerfile -t registry.com:5000/middle-proxy:latest"
$build

push="docker push registry.com:5000/middle-proxy:latest"
$push

rm middle-proxy

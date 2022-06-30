#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build -o entry-proxy main.go

cp conf.d/config-local.json config.json


build="docker build . -f Dockerfile -t registry.com:5000/entry-proxy:latest"
$build

push="docker push registry.com:5000/entry-proxy:latest"
$push

rm entry-proxy

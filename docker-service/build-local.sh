#!/bin/bash

cp ./conf.d/config-local.json config.json

export GOOS=linux
export GOARCH=amd64

go mod vendor
go fmt main.go
go build -o docker-service main.go

build="docker build . -f Dockerfile -t registry.com:5000/docker-service:latest"
$build

push="docker push registry.com:5000/docker-service:latest"
$push

rm docker-service

#!/bin/bash

cp ./conf.d/config-online.json config.json

export GOOS=linux
export GOARCH=amd64

go mod vendor
go fmt main.go
go build -o docker-service main.go

build="docker build . -f Dockerfile -t org-apps.programschool.com/docker-service:latest"
$build

push="docker push org-apps.programschool.com/docker-service:latest"
$push

rm docker-service

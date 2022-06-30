#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build -o entry-proxy main.go

cp conf.d/config-local.json config.json


build="docker build . -f Dockerfile -t org-apps.programschool.com/entry-proxy:latest"
$build

push="docker push org-apps.programschool.com/entry-proxy:latest"
$push

rm entry-proxy

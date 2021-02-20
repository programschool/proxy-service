#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build -o entry-proxy main.go


dev="-dev"

if [[ $1 = '--prod' ]]
then
    dev=""
fi

build="docker build . -f Dockerfile -t registry.cn-wulanchabu.aliyuncs.com/programschool$dev/entry-proxy:latest"
$build

push="docker push registry.cn-wulanchabu.aliyuncs.com/programschool$dev/entry-proxy:latest"
$push

rm entry-proxy
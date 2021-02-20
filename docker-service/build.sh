#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build -o docker-service main.go


dev="-dev"

if [[ $1 = '--prod' ]]
then
    dev=""
fi

build="docker build . -f Dockerfile -t registry.cn-wulanchabu.aliyuncs.com/programschool$dev/docker-service:latest"
$build

push="docker push registry.cn-wulanchabu.aliyuncs.com/programschool$dev/docker-service:latest"
$push

rm docker-service

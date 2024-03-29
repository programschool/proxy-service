#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build -o docker-service main.go


cp conf.d/config-dev.json config.json
dev="-dev"

if [[ $1 = '--prod' ]]
then
  cp conf.d/config.json config.json
  dev=""
fi

build="docker build . -f Dockerfile -t registry.cn-wulanchabu.aliyuncs.com/programschool$dev/docker-service:latest"
$build

push="docker push registry.cn-wulanchabu.aliyuncs.com/programschool$dev/docker-service:latest"
$push

rm docker-service

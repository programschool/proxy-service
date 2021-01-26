#!/bin/bash

export GOOS=linux
export GOARCH=amd64

go fmt main.go
go build main.go
scp main root@192.168.10.103:/home/services/docker-service
rm main

# 123456

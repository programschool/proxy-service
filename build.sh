#!/bin/bash

export GOOS=linux
export GOARCH=amd64

cd entry-proxy
go fmt main.go
go build main.go
scp main root@192.168.10.105:/home/services/entry-proxy
rm main
cd ..

cd middle-proxy
go fmt main.go
go build main.go
scp main root@192.168.10.104:/home/services/entry-proxy
rm main
cd ..

cd docker-service
go fmt main.go
go build main.go
scp main root@192.168.10.103:/home/services/docker-service
rm main
cd ..

# 123456

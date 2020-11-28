#!/bin/bash

export GOOS=linux
export GOARCH=amd64

cd entry-proxy
go fmt main.go
go build -o ../../../node-router main.go
cd ..


cd middle-proxy
go fmt main.go
go build -o ../../../container-node main.go
cd ..

cd docker-service-api
go fmt main.go
go build  -o ../../ main.go
cd ..

#!/bin/bash

export GOOS=linux
export GOARCH=amd64

cd entry-proxy
go build main.go
cp main ../../../node-router
cd ..


cd middle-proxy
go build main.go
cp main ../../../container-node
cd ..

#cd router
#go fmt update_container.go
#go build  -o ../build-program/router/update_container update_container.go
#
#go fmt router.go
#go build -o ../build-program/router/router router.go
#
#cd host-create
#go fmt docker_proxy.go
#go build -o ../../build-program/router/docker_proxy/docker_proxy docker_proxy.go
#
#cd ../../node-proxy
#go fmt node_proxy.go
#go build -o ../build-program/node_proxy/node_proxy node_proxy.go

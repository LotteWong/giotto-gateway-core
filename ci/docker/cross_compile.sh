#!/bin/bash

# set go env
export GO111MODULE=auto
export GOPROXY=https://goproxy.io,direct
go mod tidy

# build binary executable
GOOS=linux GOARCH=amd64 go build -o ../bin/giotto_gateway

# build docker images
commit=echo `git rev-parse --short HEAD`
docker build -f Dockerfile-Management -t giotto-gateway-management:$commit ..
docker build -f Dockerfile-Proxy -t giotto-gateway-proxy:$commit ..

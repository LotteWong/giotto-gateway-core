#!/bin/bash

# set go env
export GO111MODULE=auto
export GOPROXY=https://goproxy.io,direct
go mod tidy

# build binary executable
mkdir -p ../bin/giotto_gateway_core
GOOS=linux GOARCH=amd64 go build -o ../bin/giotto_gateway_core

# build docker images
commit=echo `git rev-parse --short HEAD`
docker build -f Dockerfile -t giotto-gateway-core:$commit ..

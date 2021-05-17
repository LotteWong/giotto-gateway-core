#!/bin/bash

# set go env
export GO111MODULE=auto
export GOPROXY=https://goproxy.io,direct
go mod tidy

# build binary executable
mkdir -p ./bin
go build -o ./bin/giotto_gateway_core

# kill processes already started
pkill -9 giotto_gateway_core

# run management backgroud server
nohup ./bin/giotto_gateway_core -config ./configs/prod/ >> ./logs/core.log 2>&1 &
echo 'nohup ./bin/giotto_gateway_core -config ./configs/prod/ >> ./logs/core.log 2>&1 &'

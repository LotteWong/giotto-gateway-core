#!/bin/bash

# set go env
export GO111MODULE=auto
export GOPROXY=https://goproxy.io,direct
go mod tidy

# build binary executable
go build -o ../bin/giotto_gateway

# kill processes already started
pkill -9 giotto_gateway

# run management backgroud server
nohup ../bin/giotto_gateway -config ../config/prod/ -endpoint management >> logs/dashboard.log 2>&1 &
echo 'nohup ../bin/giotto_gateway -config ./config/prod/ -endpoint management >> logs/management.log 2>&1 &'

# run management backgroud server
nohup ../bin/giotto_gateway -config ../config/prod/ -endpoint proxy >> logs/dashboard.log 2>&1 &
echo 'nohup ../bin/giotto_gateway -config ./config/prod/ -endpoint proxy >> logs/proxy.log 2>&1 &'

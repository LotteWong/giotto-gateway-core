#!/bin/sh

# remove docker images
commit=`git rev-parse --short HEAD`
docker rmi -f giotto-gateway-core:$commit
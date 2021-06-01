#!/bin/sh

# update docker images
commit=`git rev-parse --short HEAD`
kubectl set image deployment/giotto-gateway-core giotto-gateway-core=giotto-gateway-core:$commit
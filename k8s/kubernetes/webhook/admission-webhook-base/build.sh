#!/bin/bash

: ${DOCKER_USER:? required}

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
# build webhook
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o admission-webhook-base
# build docker image
docker build --no-cache -t ${DOCKER_USER}/admission-webhook-base:v1.0.0 .
rm -rf admission-webhook-base

docker push ${DOCKER_USER}/admission-webhook-base:v1.0.0
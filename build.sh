#!/bin/bash

VER="${1:-v1}"
DOCKER_REPO="${2:-10.190.167.198:5000}"

if [[ $VER != "" ]]; then
  rm device-plugin
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o device-plugin cmd/main.go
  docker build --no-cache -t localhost:5000/mf-device-plugin:$VER .
  docker push localhost:5000/mf-device-plugin:$VER
  docker tag localhost:5000/mf-device-plugin:$VER $DOCKER_REPO/mf-device-plugin:$VER
fi

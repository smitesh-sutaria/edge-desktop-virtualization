#!/bin/bash

# Copyright (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

VER="${1:-v1}"
DOCKER_REPO="${2:-127.0.0.1:5000}"

if ! command -v go > /dev/null ; then
  echo "Error: 'go' not found in PATH, is it installed?"
  exit 1
fi

if [[ $VER != "" ]]; then
  rm -f device-plugin
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o device-plugin cmd/main.go
  docker build --no-cache -t localhost:5000/mf-device-plugin:$VER .
  docker push localhost:5000/mf-device-plugin:$VER
  docker tag localhost:5000/mf-device-plugin:$VER $DOCKER_REPO/mf-device-plugin:$VER
fi

#!/bin/bash

# Copyright (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

GIT_REPO=https://github.com/open-edge-platform/edge-microvisor-toolkit.git
DEFAULT_TAG=3.0.20250718
IDV_JSON_PATH=idv.json

TAG=${1:-$DEFAULT_TAG}

# clone the emt repo
git clone $GIT_REPO
cd edge-microvisor-toolkit
# checkout the required TAG
git checkout $TAG

# pre-requisites
sudo ./toolkit/docs/building/prerequisites-ubuntu.sh
sudo ln -vsf /usr/lib/go-1.21/bin/go /usr/bin/go
sudo ln -vsf /usr/lib/go-1.21/bin/gofmt /usr/bin/gofmt
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# build the toolkit
cd toolkit
if [[ -z "$IDV_JSON_PATH" ]]; then
    wget https://raw.githubusercontent.com/open-edge-platform/edge-desktop-virtualization/refs/heads/emt-dv-iso/emt-dv-iso/idv.json
    cp idv.json ./imageconfigs
else
    cp $IDV_JSON_PATH ./imageconfigs/idv.json

sudo make toolchain REBUILD_TOOLS=y VALIDATE_TOOLCHAIN_GPG=n

# build the iso image
sudo make iso -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n VALIDATE_TOOLCHAIN_GPG=n CONFIG_FILE=./imageconfigs/idv.json

# copy the generated iso to same parent folder
cp ../out/images/idv/*.iso ../../.

# cleanup
cd ../../
sudo rm -rf edge-microvisor-toolkit

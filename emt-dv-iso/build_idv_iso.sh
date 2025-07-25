#!/bin/bash

# Copyright (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

set -eE

# Define color variables for readability
RED='\e[31m'
GREEN='\e[32m'
BLUE='\e[34m'
ENDCOLOR='\e[0m' # Reset to default color

# ------------------- Global Variables ----------------------------

# Git repo to build against. This can be any forked repo of EMT as well.
GIT_REPO=https://github.com/open-edge-platform/edge-microvisor-toolkit.git

# Default tag. This will be the latest EMT release tag.
DEFAULT_TAG=3.0.20250718

# Full path of the image config JSON. If this is NULL, default will be fetched.
IDV_JSON_PATH=""

# Register the current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# optional argument 1. To build against a specific tag.
TAG=${1:-$DEFAULT_TAG}

function launch_build() {
    echo -e ${BLUE}"Current working directory is: " ${GREEN}$DIR ${ENDCOLOR}
    # clone the emt repo
    echo -e "${BLUE}Cloning the EMT repo @${GREEN}${GIT_REPO}${ENDCOLOR}"
    git clone $GIT_REPO
    cd edge-microvisor-toolkit

    # checkout the required TAG
    echo -e "${BLUE}Checkout tag : ${GREEN}${TAG}${ENDCOLOR}"
    git checkout $TAG

    # pre-requisites
    echo -e "${BLUE}Installing all the pre-requisites${ENDCOLOR}"
    sudo ./toolkit/docs/building/prerequisites-ubuntu.sh
    sudo ln -vsf /usr/lib/go-1.21/bin/go /usr/bin/go
    sudo ln -vsf /usr/lib/go-1.21/bin/gofmt /usr/bin/gofmt
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER

    # build the toolkit
    cd toolkit
    if [[ -z "$IDV_JSON_PATH" ]]; then
        echo -e "${BLUE}JSON input not provided. Hence download the default one.${ENDCOLOR}"
        wget https://raw.githubusercontent.com/open-edge-platform/edge-desktop-virtualization/refs/heads/emt-dv-iso/emt-dv-iso/idv.json
        cp idv.json ./imageconfigs
    else
        echo -e "${BLUE}JSON input provided is : ${GREEN}.${ENDCOLOR}"
        cp $IDV_JSON_PATH ./imageconfigs/idv.json
    fi
    sudo make toolchain REBUILD_TOOLS=y VALIDATE_TOOLCHAIN_GPG=n

    # build the iso image
    sudo make iso -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n VALIDATE_TOOLCHAIN_GPG=n CONFIG_FILE=./imageconfigs/idv.json

    # copy the generated iso to same parent folder
    cp ../out/images/idv/*.iso ../../.
}

function cleanup() {
    echo -e "${GREEN}Performing cleanup ${ENDCOLOR}"
    cd $DIR
    sudo rm -rf edge-microvisor-toolkit
}

trap cleanup EXIT
trap cleanup ERR

#---------------------- main ------------------------

launch_build
cleanup

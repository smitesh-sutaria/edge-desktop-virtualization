#!/bin/bash

# Copyright (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# These contents may have been developed with support from one or more
# Intel-operated generative artificial intelligence solutions.

username=$(getent passwd 1000 | cut -d: -f1)

sudo chown $username:$username /opt/idv/launcher/start_all_vms.log
sudo chmod +x /opt/idv/launcher/start_all_vms.log

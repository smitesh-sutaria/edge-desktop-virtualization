#!/bin/bash

# Copyright (C) 2025 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

echo "ALL ALL=(ALL) NOPASSWD: /usr/bin/X,/usr/local/bin/idv/init/setup_sriov_vfs.sh,/usr/local/bin/idv/init/setup_display.sh,/usr/local/bin/idv/init/setup_file_permissions.sh,/usr/local/bin/idv/launcher/start_vm.sh,/usr/local/bin/idv/launcher/start_all_vms.sh,/usr/local/bin/idv/launcher/stop_vm.sh,/usr/local/bin/idv/launcher/stop_all_vms.sh" | sudo tee -a /etc/sudoers.d/guest > /dev/null


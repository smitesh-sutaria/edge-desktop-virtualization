#!/bin/bash

# These contents may have been developed with support from one or more
# Intel-operated generative artificial intelligence solutions.

username=$(getent passwd 1000 | cut -d: -f1)

echo "$username ALL=(ALL) NOPASSWD: /opt/idv/launcher/start_vm.sh" | sudo tee /etc/sudoers.d/start_vm > /dev/null
echo "$username ALL=(ALL) NOPASSWD: /opt/idv/launcher/start_all_vms.sh" | sudo tee /etc/sudoers.d/start_all_vms > /dev/null

echo "$username ALL=(ALL) NOPASSWD: /opt/idv/launcher/stop_vm.sh" | sudo tee /etc/sudoers.d/stop_vm > /dev/null
echo "$username ALL=(ALL) NOPASSWD: /opt/idv/launcher/stop_all_vms.sh" | sudo tee /etc/sudoers.d/stop_all_vms > /dev/null

sudo chown $username:$username /opt/idv/launcher/start_all_vms.log
sudo chmod +x /opt/idv/launcher/start_all_vms.log

loginctl enable-linger $username

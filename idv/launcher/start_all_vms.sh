#!/bin/bash

# These contents may have been developed with support from one or more
# Intel-operated generative artificial intelligence solutions.

source vm.conf

declare -A VM_LIST

VM_LIST=()
for ((counter = 1; counter <= ${guest}; counter++)); do
  vm="vm${counter}"
  VM_LIST[${#VM_LIST[@]}]=${vm}
done

trap 'trap " " SIGTERM; kill 0; wait' SIGINT SIGTERM

for vm in "${VM_LIST[@]}"; do
    name="${vm}_name"
    echo "Starting Windows Guest ${!name} ..."
    ram="${vm}_ram"
    cpu="${vm}_cores"
    firmware_file="${vm}_firmware_file"
    qcow2_file="${vm}_qcow2_file"
    connector="${vm}_connector0"
    ssh="${vm}_ssh"
    winrdp="${vm}_winrdp"
    winrm="${vm}_winrm"
    usb="${vm}_usb"

    sudo ./start_vm.sh -m ${!ram}G -c ${!cpu} -n ${!name} -f ${!firmware_file} -d ${!qcow2_file} --display connectors.0=${!connector} -p ssh=${!ssh},winrdp=${!winrdp},winrm=${!winrm} -u ${!usb} &
    
    # Added sleep time of 3 seconds to make sure there is no issue related to swtpm socket
    sleep 3
done

wait

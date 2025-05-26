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

for vm in "${VM_LIST[@]}"; do
    name="${vm}_name"
    sudo ./stop_vm.sh ${!name}
done

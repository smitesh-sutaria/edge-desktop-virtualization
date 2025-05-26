#!/bin/bash

grep_output=$(ps aux | grep qemu | grep -i $1)

if [ -n "$grep_output" ]; then
    pid=$(echo "$grep_output" | awk '{print $2}')
    echo "Stopping VM: $1"
    sudo kill -9 $pid
else
    echo "Could not find QEMU process for $1"
fi

#!/bin/bash

# Abort immediately if any command fails
set -e

usage() {
    cat <<EOF
Usage: $0 [--skip-confirm] <device name> <os image>
    --skip-confirm
        Optional. If specified, bypasses the confirmation prompt.
    <device_name>
        Example: /dev/sda, /dev/nvme0n1
        Device to flash the OS image onto
        !!! All partitions will be overwritten from the OS image !!!
    <os image>
        Example: tiber-readonly-mf-3.0.20250117.1126.raw
        The tiber OS image to flash using bmaptool
EOF
    exit 1
}

# Check if the first argument is --skip-confirm
skip_confirm=false
if [[ "$1" == "--skip-confirm" ]]; then
    skip_confirm=true
    shift # Shift arguments to remove --skip-confirm from the list
fi

if [[ $# -ne 2 ]]; then
    usage
fi

# Argument 1 is the root device name (/dev/sda, /dev/nvme0n1, etc.)
target="$1"
if [[ ! -b "$target" ]]; then
    echo "Target device '$target' does not exist or is not accessible"
    usage
fi

# Argument 2 is the OS image
image="$2"
if [[ ! -f "$image" ]]; then
    echo "Image '$image' does not exist or is not accessible"
    usage
fi

for cmd in bmaptool sgdisk sfdisk e2fsck growpart resize2fs; do
    if ! command -v "$cmd" &> /dev/null; then
        echo "Error: Required command '$cmd' is not installed or not in PATH."
        exit 1
    fi
done

if [[ $UID -ne 0 ]]; then
    echo "This script must be run as root (via sudo, etc.)"
    exit 1
fi

echo "Current partition table of $target:"
sfdisk -l "$target"

# Ask for confirmation unless --skip-confirm is specified
if [[ "$skip_confirm" == false ]]; then
    echo ""
    read -p "Are you sure you want to proceed? This will overwrite all partitions on $target that are listed above. (y/n): " confirm
    if [[ ! "$confirm" =~ ^(Y|y|yes|YES|Yes)$ ]]; then
        echo "Operation canceled."
        exit 1
    fi
fi

echo "Loading image '$image' onto '$target'..."
bmaptool copy --nobmap "$image" "$target"

echo "Restoring backup GPT header to end of $target"
sgdisk -e "$target"

echo "Verifying the partition table on $target"
sfdisk -Vl "$target"

if [[ "$target" =~ nvme ]]; then
    # Target is an nvme drive, with p# naming
    system_partition="${target}p2"
    data_partition="${target}p3"
else
    # Target is not an nvme drive (i.e. sda), partitions are just appended #
    system_partition="${target}2"
    data_partition="${target}3"
fi
if [[ ! -b "$system_partition" ]]; then
    echo "System partition '$system_partition' does not exist or is not accessible, possibly this script can't handle this drive type"
    exit 1
fi
if [[ ! -b "$data_partition" ]]; then
    echo "System partition '$data_partition' does not exist or is not accessible, possibly this script can't handle this drive type"
    exit 1
fi

echo "Performing filesystem check on $data_partition"
e2fsck -f -y -v "$data_partition"

# List the free/unpartitioned space on the target drive.
# Get the total number of free sectors, then divide that by 5. Not much space
# is needed for the root partition, everything else lives on the data partition (/home, /opt, ...)
# The data partition will be moved by this much, then expanded to fill the 2nd half of the free space
# The 2nd partition (root) will then be expanded to fill the first half of the free space
# in case of 1TB NVME (1945138575), shift_sector given is 945138575
shift_sector=$(( $(sfdisk -F "$target" | tail -n 1 | awk '{print $3}') / 5 ))

if [[ "$shift_sector" -eq 0 ]]; then
    echo "Error: Unable to determine free sectors on $target"
    exit 1
fi
if [[ "$shift_sector" -le 2000000 ]]; then
    echo "Error: Insufficient free sectors on $target to move data partition"
    exit 1
fi

# Move the data partition
echo "Moving the data partition to the middle of the available free space"
echo "+${shift_sector}" | sfdisk --move-data --force $target -N 3

echo "Growing the system partition (partition 2)"
growpart "$target" 2

echo "Checking & repairing the expanded system partition"
e2fsck -f -p "$system_partition"

echo "Expanding the system partition's filesystem"
resize2fs "$system_partition"

echo "Growing the data partition (partition 2)"
growpart "$target" 3

echo "Checking & repairing the expanded data partition"
e2fsck -f -p "$data_partition"

echo "Expanding the data partition's filesystem"
resize2fs "$data_partition"

echo "Final partition table of $target:"
sfdisk -l "$target"

echo "Image $image has been successfully loaded onto $target"

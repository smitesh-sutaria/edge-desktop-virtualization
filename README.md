# IDV solution with SR-IOV setup

> These contents may have been developed with support from one or more Intel-operated generative artificial intelligence solutions.

## Table of Contents
- [Pre-requisites](#pre-requisites)
  - [Create a directory to save qcow2 image and firmware files](#create-a-directory-to-save-qcow2-image-and-firmware-files)
  - [SR-IOV virtual functions enumeration](#sr-iov-virtual-functions-enumeration)
  - [Run X server](#run-x-server)
- [Launch Windows11 virtual machines](#launch-windows11-virtual-machines)

## Pre-requisites

### Create a directory to save qcow2 image and firmware files

- Create a '/opt/qcow2' directory and local user should have write access to it.

  ```bash
  sudo mkdir /opt/qcow2
  sudo chmod -R 755 /opt/qcow2/
  sudo chown -R <your_username>:<your_username> /opt/qcow2/

### SR-IOV virtual functions enumeration

- Move to the 'idv' directory using the following command
  
  ```bash
  cd idv
  ```

- Run the `setup_sriov_vfs.sh` script with superuser privileges to enable virtual functions
  
  ```bash
  sudo ./setup_sriov_vfs.sh
  ```
  
  - Ensure SR-IOV is enabled. Run the command `lspci` to verify the presence of 7 virtual functions. 

  ```bash
  lspci | grep VGA
  ```
  
  Expected output:

  ```bash
  00:02.0 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.1 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.2 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.3 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.4 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.5 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.6 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  00:02.7 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
  ```

### Run X server

- Run the following command to start X server
  
  ```bash
  sudo X
  ```

  **Note**: After running the above command, you will see a blank screen. Ensure you have SSH access to the machine for the next steps.

- Grant X11 server access to local users

  SSH into the machine and run the following commands to grant access to all users:
  ```bash
  export DISPLAY=:0
  xhost +
  ```

## Launch Windows11 virtual machines

- Once you have completed all the above steps, move to the working directory. Run the following command

  ```bash
  cd idv
  ```

- The `vm.conf` file contains configuration parameters for the virtual machines. Modify this file to specify the number of VMs to launch and their respective settings. 

  - Set the `guest` variable to the number of VMs to launch
  - Fill in the required configuration parameters for each VM in the right order. If `guest` is set to `2`, modify/set the variables starting with `vm1_*` and `vm2_*`

  Example:

  ```ini
  # Memory in GB
  vm1_ram=3
  # Name of the VM
  vm1_name=windows_vm1
  # Number of CPU cores
  vm1_cores=3
  # Name of the firmware file present in `/opt/qcow2` directory
  vm1_firmware_file=OVMF_VARS_windows1.fd
  # Name of the qcow2 file present in `/opt/qcow2` directory
  vm1_qcow2_file=win1.qcow2
  # Name of the display connector (monitor)
  vm1_connector0=HDMI-1
  # SSH port for the VM
  vm1_ssh=4444
  # WinRDP port for the VM
  vm1_winrdp=3389
  # WinRM port for the VM
  vm1_winrm=5986
  ```

  **Note:** Set unique values for ssh, winrdp and winrm ports to avoid conflicts when launching multiple VMs.

- Run the `start_all_vms` script with superuser privileges to launch the VMs
  
  ```bash
  sudo ./start_all_vms.sh
  ```

  Verify that the VMs are running by checking the process list
  
  ```bash
  ps aux | grep qemu
  ```

## Troubleshooting
  - **Issue:** VMs fail to launch with `start_all_vms.sh`.
  - **Solution:** Check the `vm.conf` file for errors or missing parameters. Ensure all required files (e.g., firmware and qcow2 files) are present in `/opt/qcow2`.


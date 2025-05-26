# Desktop Virtualization solution with graphics SR-IOV

This file contains steps to launch virtual machines using a system service.

## Table of Contents
- [Modify the VM configuration](#modify-the-vm-configuration)
- [Run script to copy necessary files to `/opt` directory](#run-script-to-copy-necessary-files-to-opt-directory)
- [Enable and start `idv-init` service](#enable-and-start-idv-init-service)
- [Enable and start `idv-launcher` service](#enable-and-start-idv-launcher-service)
- [Troubleshooting](#troubleshooting)

## Modify the VM configuration

- The `vm.conf` file in `idv` directory contains configuration parameters for the virtual machines. Modify this file to specify the number of VMs to launch and their respective settings. 

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
  # Path of firmware file
  vm1_firmware_file=OVMF_VARS_windows1.fd
  # Path of qcow2 file
  vm1_qcow2_file=win1.qcow2
  # Name of the display connector (monitor)
  vm1_connector0=HDMI-1
  # Comma separated list of USB devices to attach to the VM in the format: <hostbus>-<hostport>, where hostport is the end port to which the device is attached
  vm1_usb=3-1.1,3-1.2,3-1.3,3-1.4
  # SSH port for the VM
  vm1_ssh=4444
  # WinRDP port for the VM
  vm1_winrdp=3389
  # WinRM port for the VM
  vm1_winrm=5986

## Run script to copy necessary files to `/opt` directory

- Move to the `idv` directory and run the `copy_files.sh` file with superuser privileges using the following command

  ```bash
  cd idv
  sudo ./copy_files.sh
  ```
  This copies all the scripts and services to appropriate directories.

## Enable and start `idv-init` service

  The `idv-init` service initializes the environment by enumerating SR-IOV virtual functions, starting the X server and setting up permissions required to run the scripts to launch VMs. This is a prerequisite for launching the virtual machines.

- Run the following command to enable `idv-init` service
  
  ```bash
  sudo systemctl enable idv-init.service
  ```

- Run the following command to start `idv-init` service
  
  ```bash
  sudo systemctl start idv-init.service
  ```

- Verify that the service is running:

    ```bash
    sudo systemctl status idv-init.service
    ```
   **Note**: After starting the idv-init service, the screen will go blank because X is running. Ensure you have SSH access to the machine for the next steps.

## Enable and start `idv-launcher` service

  The `idv-launcher` service launches the configured virtual machines in their respective monitors.

- Run the following command to enable `idv-launcher` service
  
  ```bash
  systemctl --user enable idv-launcher.service
  ```

- Run the following command to start `idv-launcher` service
  
  ```bash
  systemctl --user start idv-launcher.service
  ```

- Verify that the service is running:

  ```bash
  systemctl --user status idv-launcher.service
  ```
   **Note**: Once the idv-launcher service starts, all the VMs should be launched in respective monitors.

## Troubleshooting

- If the `idv-init` service fails to start, check the service logs using the following command:
  
  ```bash
  sudo journalctl -u idv-init.service
  ```
  Ensure that all required files are present in `/opt/idv`.


- If the VMs do not launch after starting the `idv-launcher` service, check the service logs using the following command:

  ```bash
  journalctl --user -u idv-launcher.service
  ```

  You can also check the `start_all_vms.log` in `/opt/idv/launcher` directory for errors using the command:

  ```bash
  sudo cat /opt/idv/launcher/start_all_vms.log
  ```  
  Ensure that the `vm.conf` file is correctly configured and all required files (e.g., firmware and qcow2 files) are present and the file paths are valid.

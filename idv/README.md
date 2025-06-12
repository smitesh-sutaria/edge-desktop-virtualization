# Desktop Virtualization solution with graphics SR-IOV

This file contains steps to launch virtual machines using a system service.

## Table of Contents
- [Modify the VM configuration](#modify-the-vm-configuration)
- [Run IDV services via an RPM package](#run-idv-services-via-an-rpm-package)
- [Run script to copy necessary files to `/opt` directory](#run-script-to-copy-necessary-files-to-opt-directory)
- [Enable and start `idv-init` service](#enable-and-start-idv-init-service)
- [Enable and start `idv-launcher` service](#enable-and-start-idv-launcher-service)
- [Troubleshooting](#troubleshooting)

## Modify the VM configuration

- The `vm.conf` file in `idv/launcher` directory contains configuration parameters for the virtual machines. Modify this file to specify the number of VMs to launch and their respective settings. 

  - Set the `guest` variable to the number of VMs to launch.
  - Set the `OVMF_CODE_FILE` variable to the path of OVMF_CODE.fd file.
  - Fill in the required configuration parameters for each VM in the right order. If `guest` is set to `2`, modify/set the variables starting with `vm1_*` and `vm2_*`

  Example:

  ```ini
  # Memory in GB
  vm1_ram=3
  # OS to be configured
  vm1_os=windows
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
  # WinRDP port for the VM - Set this only for Windows VM.
  vm1_winrdp=3389
  # WinRM port for the VM - Set this only for Windows VM
  vm1_winrm=5986
  ```

    **Note:** Set unique values for ssh, winrdp and winrm ports to avoid conflicts when launching multiple Windows VMs.

## Run IDV services via an RPM package
- For detailed instructions on running IDV services using an RPM package, follow the instructions in [RPM Packaging Guide](rpm-packaging-guide.md).

## Manual steps to run IDV service

## Run script to copy necessary files to `/opt` directory

- Move to the `idv` directory and run the `copy_files.sh` file with superuser privileges using the following command

  ```bash
  cd idv
  sudo chmod +x copy_files.sh
  sudo ./copy_files.sh
  ```
  This copies all the scripts and services to appropriate directories.

## Enable and start `idv-init` service

  The `idv-init` service initializes the environment by enumerating SR-IOV virtual functions, starting the X server and setting up permissions required to run the scripts to launch VMs. This is a prerequisite for launching the virtual machines.

- Run the following command to enable `idv-init` service
  
  ```bash
  systemctl --user enable idv-init.service
  ```

- Run the following command to start `idv-init` service
  
  ```bash
  systemctl --user start idv-init.service
  ```

- Verify that the service is running:

    ```bash
    systemctl --user status idv-init.service
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
  journalctl --user -xeu idv-init.service
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

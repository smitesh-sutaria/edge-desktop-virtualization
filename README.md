# Desktop Virtualization solution with graphics SR-IOV

## Table of Contents
- [Desktop Virtualization solution with graphics SR-IOV](#desktop-virtualization-solution-with-graphics-sr-iov)
  - [Table of Contents](#table-of-contents)
  - [Pre-requisites](#pre-requisites)
    - [System Requirements](#system-requirements)
      - [Recommended Hardware Configuration](#recommended-hardware-configuration)
    - [Build EMT](#build-emt)
    - [Install EMT](#install-emt)
    - [Generate Virtual Machine qcow2 with required drivers for SR-IOV](#generate-virtual-machine-qcow2-with-required-drivers-for-sr-iov)
    - [SR-IOV virtual functions enumeration](#sr-iov-virtual-functions-enumeration)
    - [Display Setup](#display-setup)
      - [Disable DPMS and screen blanking on the X Window System](#disable-dpms-and-screen-blanking-on-the-x-window-system)
      - [Start X server](#start-x-server)
  - [VM configuration file](#virtual-machine-configuration-file)
  - [Launch one Windows11 virtual machine](#launch-one-windows11-virtual-machine)
  - [Launch one Ubuntu virtual machine](#launch-one-ubuntu-virtual-machine)
  - [Launch multiple virtual machines](#launch-multiple-virtual-machines)
  - [Stop virtual machines](#stop-virtual-machines)
  - [Troubleshooting](#troubleshooting)

## Pre-requisites

### System Requirements

Edge Microvisor Toolkit + Graphics SR-IOV is designed to support all Intel® Core platforms from 12th gen onwards.

This software is validated on below:

|         Core™         |
| ----------------------|
| 12th Gen Intel® Core™ |
| 13th Gen Intel® Core™ |

#### Recommended Hardware Configuration

| Component    | Edge Microvisor Toolkit + graphics SR-IOV|
|--------------|-----------------------------------|
| CPU          | Intel® Core (12th gen and higher) |
| RAM          | 64GB recommended                  |
| Storage      | 500 GB SSD or NVMe minimum        |
| Networking   | 1GbE Ethernet                     |


### Build EMT

Reference to the build steps as mentioned here : [EMT Image build](https://github.com/smitesh-sutaria/edge-microvisor-toolkit/blob/3.0/docs/developer-guide/get-started/building-howto.md)

#### Pre-requisite
- Ubuntu 22.04
- Install the dependencies mentioned [here](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/toolkit/docs/building/prerequisites-ubuntu.md)

#### Image Build Steps 

**Step 1: Clone EMT repo**
```bash
git clone https://github.com/open-edge-platform/edge-microvisor-toolkit.git 
# checkout to the 3.0 tag 
git checkout 3.0.20250411
```
**Step 2: Edit the Chroot env in the go code [toolkit/tools/internal/safechroot/safechroot.go](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0.20250411/toolkit/tools/internal/safechroot/safechroot.go)** 
```go
# add the following lines under "defaultChrootEnv" variable declaration, after the line 102
fmt.Sprintf("https_proxy=%s", os.Getenv("https_proxy")),
fmt.Sprintf("no_proxy=%s", os.Getenv("no_proxy")),
```
It should look something like this 
![safechroot.go](docs/artifacts/proxy-go.png)

**Step 3: Build the toolkit**
```bash
cd edge-microvisor-toolkit/toolkit
sudo -E  make toolchain REBUILD_TOOLS=y
```
**Step 4: Build the image** 
Build EMT image for graphics SR-IOV using the spec [edge-image-mf-dev.json](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0-dev/toolkit/imageconfigs/edge-image-mf-dev.json)
```bash 
sudo -E make image -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n CONFIG_FILE=imageconfigs/edge-image-mf-dev.json
# created image will be available under "edge-microvisor-toolkit/out/images/edge-image-mf-dev"
```
> ⚠️ **Note: Please remove "intel" related proxy from "no_proxy" system env variable before step 3**

### Install EMT

To Flash EMT MF image on a NUC follow [EMT image installation docs](https://github.com/intel-innersource/applications.virtualization.maverickflats-tiberos-itep/blob/vm_sidecar_dev_plugin/tiber/tiber_flash_partition.md) 

To verify checkout [Other methods](https://github.com/smitesh-sutaria/edge-microvisor-toolkit/blob/3.0/docs/developer-guide/get-started/installation-howto.md)

### Generate Virtual Machine qcow2 with required drivers for SR-IOV

Follow the qcow2 creation for windows till post install launch from this readme.

https://github.com/ThunderSoft-SRIOV/sriov/blob/main/docs/deploy-windows-vm.md#microsoft-windows-11-vm

### SR-IOV virtual functions enumeration

- Move to the 'idv' directory using the following command
  
  ```bash
  cd idv/init
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

### Display setup

## Disable DPMS and screen blanking on the X Window System

- Move to the 'idv' directory using the following command
  
  ```bash
  cd idv/init
  ```

- Run the `setup_display.sh` script with superuser privileges to disable DPMS and screen blanking on the X Window System
  
  ```bash
  sudo ./setup_display.sh

## Start X server

- Run the following command to start X server
  
  ```bash
  sudo X
  ```

  **Note**: After running the above command, you will see a blank screen because X is running. To access the console, try control+alt+f3. To return to X, switch back with control-alt-f2. If X is not currently the active display, the VMs will not boot. Ensure you have SSH access to the machine for the next steps.

- Grant X11 server access to local users

  SSH into the machine and run the following commands to grant access to all users:
  ```bash
  export DISPLAY=:0
  xhost +
  ```

## Virtual machine configuration file

- The `vm.conf` file in `idv/launcher` directory consists of the following parameters - 

| Parameter           | Description                                      | Example Value                             |
|---------------------|--------------------------------------------------|-------------------------------------------|
| `vm1_ram`           | Memory allocated to the VM (in GB).              | `3`                                       |
| `vm1_os`            | Operating system of the VM.                      | `windows` or `ubuntu`                     |
| `vm1_name`          | Name of the VM.                                  | `vm1`                                     |
| `vm1_cores`         | Number of CPU cores allocated to the VM.         | `3`                                       |
| `vm1_firmware_file` | Path of firmware (.fd) file.                     | `/opt/qcow2/win1.fd`                      |
| `vm1_qcow2_file`    | Path of qcow2 file.                              | `/opt/qcow2/win1.qcow2`                   |
| `vm1_connector0`    | Display connector for the VM.                    | `HDMI-1`                                  |
| `vm1_usb`           | USB devices to attach (comma-separated)          | `3-1.1,3-1.2,3-1.3,3-1.4`                 |
| `vm1_ssh`           | SSH port for the VM.                             | `4444` (for windows), `2222` (for Ubuntu) |
| `vm1_winrdp`        | WinRDP port (Set this only for Windows VM)       | `3389`                                    |
| `vm1_winrm`         | WinRM port (Set this only for Windows VM)        | `5986`                                    |

- `vm1_usb` should be a comma separated list of USB devices to attach to the VM in the format: <hostbus>-<hostport>, where hostbus is the bus number and hostport is the end port to which the device is attached.

## Launch one Windows11 virtual machine

- Once you have completed all the above steps, move to the working directory using the following command

  ```bash
  cd idv/launcher
  ```

- The `vm.conf` file contains configuration parameters for the virtual machines. Modify the variables starting with `vm1_*` to set the configuration parameters.

  Example:

  ```ini
  vm1_ram=3
  vm1_os=windows
  vm1_name=windows_vm1
  vm1_cores=3
  vm1_firmware_file=/opt/qcow2/OVMF_VARS_windows1.fd
  vm1_qcow2_file=/opt/qcow2/win1.qcow2
  vm1_connector0=HDMI-1
  vm1_usb=3-1.1,3-1.2,3-1.3,3-1.4
  vm1_ssh=4444
  vm1_winrdp=3389
  vm1_winrm=5986
  ```

  - Set the `OVMF_CODE_FILE` variable to the path of OVMF_CODE.fd file.

- Run the `start_vm` script with superuser privileges to launch the VM
  
  ```bash
  sudo ./start_vm.sh
  ```

  Verify that the VM is running by checking the process list
  
  ```bash
  ps aux | grep qemu
  ```

## Launch one Ubuntu virtual machine

- Move to the working directory using the following command

  ```bash
  cd idv/launcher
  ```

- Modify the variables in `vm.conf` file starting with `vm1_*` to set the configuration parameters.

  Example:

  ```ini
  vm1_ram=3
  vm1_os=ubuntu
  vm1_name=ubuntu_vm1
  vm1_cores=3
  vm1_firmware_file=/opt/qcow2/ubuntu.fd
  vm1_qcow2_file=/opt/qcow2/ubuntu.qcow2
  vm1_connector0=HDMI-1
  vm1_usb=3-3.1,3-3.2,3-3.3,3-3.4
  vm1_ssh=2222
  ```

  - Set the `OVMF_CODE_FILE` variable to the path of OVMF_CODE.fd file.

- Run the `start_vm` script with superuser privileges to launch the VM
  
  ```bash
  sudo ./start_vm.sh
  ```

  Verify that the VM is running by checking the process list
  
  ```bash
  ps aux | grep qemu
  ```

## Launch multiple virtual machines

- Move to the working directory using the following command

  ```bash
  cd idv/launcher
  ```

- Modify the `vm.conf` file to specify the number of VMs to launch and their respective settings. 

  - Set the `guest` variable to the number of VMs to launch
  - Fill in the required configuration parameters for each VM in the right order. If `guest` is set to `2`, modify/set the variables starting with `vm1_*` and `vm2_*`

  Example:

  ```ini
  # Windows VM
  vm1_ram=3
  vm1_os=windows
  vm1_name=windows_vm1
  vm1_cores=3
  vm1_firmware_file=/opt/qcow2/OVMF_VARS_windows1.fd
  vm1_qcow2_file=/opt/qcow2/win1.qcow2
  vm1_connector0=HDMI-1
  vm1_usb=3-1.1,3-1.3,3-1.3,3-1.4
  vm1_ssh=4444
  vm1_winrdp=3389
  vm1_winrm=5986

  # Ubuntu VM
  vm2_ram=3
  vm2_os=ubuntu
  vm2_name=ubuntu_vm1
  vm2_cores=3
  vm2_firmware_file=/opt/qcow2/ubuntu.fd
  vm2_qcow2_file=/opt/qcow2/ubuntu.qcow2
  vm2_connector0=HDMI-1
  vm2_usb=3-3.1,3-3.2,3-3.3,3-3.4
  vm2_ssh=2222
  ```

  **Note:** Set unique values for ssh, winrdp and winrm ports to avoid conflicts when launching multiple VMs.

- Run the `start_all_vms` script with superuser privileges to launch the VMs
  
  ```bash
  sudo ./start_all_vms.sh
  ```

  **Note:** Using the above configuration, a combination of Ubuntu and Windows VMs can be launched.

  Verify that the VMs are running by checking the process list
  
  ```bash
  ps aux | grep qemu
  ```

## Stop virtual machines

- To stop all VMs, open a new terminal and run the following commands
  
  ```bash
  cd idv/launcher
  sudo ./stop_all_vms.sh
  ```

## Troubleshooting
  - **Issue:** VMs fail to launch with `start_all_vms.sh`.
  - **Solution:** Check the `vm.conf` file for errors or missing parameters. Ensure all required files (e.g., firmware and qcow2 file) are present and file paths are valid.

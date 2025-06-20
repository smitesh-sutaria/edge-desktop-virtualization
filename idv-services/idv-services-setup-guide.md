# This document contains steps to enable and start idv services on an EMT image.

# Table of Contents

1. [Modify VM Configuration](#modify-vm-configuration)
2. [Reload System Daemon](#reload-system-daemon)
3. [Enable IDV Services](#enable-idv-services)
4. [Start `idv-init.service`](#start-idv-init-service)
5. [Start `idv-launcher.service`](#start-idv-launcher-service)
6. [Troubleshooting](#troubleshooting)

## Modify VM configuration

- Refer to the [Modify VM configuration file](modify-vm-config-file.md) for details on how to modify the VM configuration file.

## Reload system daemon

- Run the following command to reload system-daemon
    
  ```bash
  systemctl --user daemon-reload
  ```

## Enable idv services

- Run the following command to enable `idv-init.service` and `idv-launcher.service`
  
  ```bash
  systemctl --user preset-all
  ```

## Start `idv-init` service

  The `idv-init` service initializes the environment by enumerating SR-IOV virtual functions, starting the X server. This is a prerequisite for launching the virtual machines.

- Run the following command to start `idv-init` service
    
    ```bash
    systemctl --user start idv-init.service
    ```

- Verify that the service is running:

  ```bash
  systemctl --user status idv-init.service
  ```
  **Note**: After starting the idv-init service, the screen will go blank because X is running. Ensure you have SSH access to the machine for the next steps.

## Start `idv-launcher` service

  The `idv-launcher` service launches the configured virtual machines in their respective monitors.

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

- If the VMs do not launch after starting the `idv-launcher` service, check the service logs using the following command:

  ```bash
  sudo journalctl -t idv-services
  ```  
  Ensure that the `vm.conf` file is correctly configured and all required files (e.g., firmware and qcow2 files) are present and the file paths are valid.
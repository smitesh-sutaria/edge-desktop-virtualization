# Run IDV Services via RPM

This guide provides step-by-step instructions to package and run IDV services using RPM.

## Table of Contents
1. [Steps to Package and Install IDV Services](#steps-to-package-and-install-idv-services)
   - [Step 1: Create a Tar File](#step-1-create-a-tar-file)
   - [Step 2: Copy Files to RPM SOURCES](#step-2-copy-files-to-rpm-sources)
   - [Step 3: Setup Permissions for Running Scripts](#step-3-setup-permissions-for-running-scripts)
   - [Step 4: Setup RPM Environment](#step-4-setup-rpm-environment)
     - [Step 4.1: Install RPM](#step-41-install-rpm)
     - [Step 4.2: Install RPM-Build](#step-42-install-rpm-build)
     - [Step 4.3: Create the RPM Build Environment](#step-43-create-the-rpm-build-environment)
     - [Step 4.4: Configure the RPM Build Environment](#step-44-configure-the-rpm-build-environment)
   - [Step 5: Build the RPM Package](#step-5-build-the-rpm-package)
   - [Step 6: Install the RPM](#step-6-install-the-rpm-package)
2. [Uninstall the RPM Package](#uninstall-the-rpm-package)
3. [Post-Reboot Instructions](#post-reboot-instructions)
4. [Troubleshooting](#troubleshooting)

### Steps to Package and Install IDV Services

## Step 1: Create a Tar File
- Run the following commands to create a tar file for the IDV solution:

```bash
sudo tar --transform='s,^,idv-solution-1.0/,' -czf idv-solution-1.0.tar.gz launcher init
sudo chmod +x idv-solution-1.0.tar.gz
```

## Step 2: Copy Files to RPM SOURCES
- Run the script to copy files to the RPM SOURCES directory:

```bash
./setup_rpm_source.sh
```
> **Note**: **Do not** run the above with sudo.

## Step 3: Setup Permissions for Running Scripts
- Run the following command to set up permissions for running scripts:

```bash
sudo ./setup_permissions.sh
```

## Step 4: Setup RPM Environment
  ## Step 4.1: Install RPM

  ```bash
  sudo -E dnf install rpm
  ```

  ## Step 4.2: Install RPM-Build

  ```bash
  sudo -E dnf install rpm-build
  ```

  ## Step 4.3: Create the RPM Build Environment

  ```bash
  mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
  ```

  ## Step 4.4: Configure the RPM Build Environment

  ```bash
  echo '%_topdir %(echo $HOME)/rpmbuild' > ~/.rpmmacros
  ```

## Step 5: Build the RPM Package

  ```bash
  rpmbuild -ba ~/rpmbuild/SPECS/idv-solution.spec
  ```
  - If successful, the RPM will be created in the ~/rpmbuild/RPMS/noarch/ directory.
  > **Note**: If the build fails, check the .spec file for syntax errors or missing dependencies. 

## Step 6: Install the RPM package

  ```bash
  sudo rpm -ivh ~/rpmbuild/RPMS/noarch/idv-solution-1.0-1.emt3.noarch.rpm
  ```
  
  > **Note**: If the output of above command may contain a message something similar to -

    ```ini
    Reload daemon failed: Transport endpoint is not connected
    Failed to start jobs: Transport endpoint is not connected
    ```
  
    Please verify the status of the service using the commands below. If the services have started, VMs will be launched in their respective monitors.

  - After installation, verify that the services are installed, enabled, and running:

    ```bash
    systemctl status --user idv-int.service
    systemctl status --user idv-launcher.service
    ```

### Uninstall the RPM Package

- To uninstall the RPM package, run:

  ```bash
  sudo rpm -e idv-solution
  ```

### Post-Reboot Instructions

- If the machine is rebooted, navigate to the `idv/` directory and run the following command to reset permissions:

  ```bash
  sudo ./setup_permissions.sh
  ```
  - Once this script is executed, the IDV services  (idv-int.service and idv-launcher.service) should start automatically. Verify their status using:
  
  ```bash
  systemctl status --user idv-int.service
  systemctl status --user  idv-launcher.service
  ```

### Troubleshooting

- If the `idv-init` service fails to start, check the service logs using the following command:

  ```bash
  journalctl --user -u idv-init.service
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

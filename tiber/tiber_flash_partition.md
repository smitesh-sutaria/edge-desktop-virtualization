# Flashing TiberOS 3.0 SR-IOV Image and creating partitions on Host System

This guide provides instructions to bring up the TiberOS 3.0 SR-IOV image onto a host system. 
The steps include 
1. Booting Host system with Ubuntu Live CD.
2. Setting up the environment to flash and create partitions. 
3. Flashing TiberOS onto NVMe/SSD and creating partitions.

    > [!Note]
    > The TiberOS image used here is **mutable**

## 1. Boot Ubuntu Live Disk and Setup Flashing Environment

1. Boot from Ubuntu Live Disk:
   - Insert the Ubuntu Live Disk into the host system and boot from it.
   - Select "Try Ubuntu" to boot into the live environment.
   - Once it boots, open Terminal and continue with below commands

2. Set HTTPS proxy:
    ```sh
    export http_proxy=http://proxy-dmz.intel.com:911
    export https_proxy=http://proxy-dmz.intel.com:912
    ```

3. Set Proxy for APT:
    ```sh
    sudo nano /etc/apt/apt.conf
    ```
    Add below lines
    ```
    Acquire::http::Proxy "http://proxy-dmz.intel.com:911";
    Acquire::https::Proxy "http://proxy-dmz.intel.com:912";
    ```

4. Update the Package List:
    ```sh
    sudo apt update
    ```

5. Download and Install the `bmap-tools` package:
    ```sh
    wget http://ports.ubuntu.com/pool/universe/b/bmap-tools/bmap-tools_3.5-2_all.deb
    sudo dpkg -i bmap-tools_3.5-2_all.deb
    ```

6. Set the password for the `ubuntu` user:
    ```sh
    sudo passwd ubuntu
    ```

7. Install necessary utilities:
    ```sh
    sudo apt install -y cloud-utils net-tools openssh-server gdisk util-linux e2fsprogs cloud-guest-utils
    ```

**Now the Ubuntu Live Image is updated with the tools needed to flash the TiberOS Core image to NVME/SSD and create the partition..**

8. Search for Gparted tool in Ubuntu Search bar
    - Look at the various drives installed in the system to confirm the one you want to load the TiberOS image onto. Proceeding will overwrite all other partitions & data on that drive, so you want to be sure.
    - No changes need to be made with this tool beyond validating the correct drive (for example, /dev/nvme0n1 or /dev/sda)

9. Flash TiberOS RAW image to NVME/SSD:
    - TiberOS SR-IOV images are hosted in [link](https://af01p-png.devtools.intel.com/artifactory/tiberos-png-local/mf/3.0/) 
    - Download any latest image by navigating to folders *Ex.: 20250321.0803/edge-readonly-mf-dev-3.0.20250321.0803.raw.gz*.
    - Use the script to load that onto the desired drive.
      - /dev/nvme0n1 for NVME
      - /dev/sda for SSD
    - *Ex. for NVME*
      ```sh
      sudo ./tiber/scripts/partition_image.sh /dev/nvme0n1 ~/Downloads/edge-readonly-mf-dev-3.0.20250321.0803.raw.gz
      ```

**TiberOS 3.0 SR-IOV Image is now flashed onto NVME/SSD, reboot the system**

10. Unmount the Ubuntu Disk, and Reboot:
    ```sh
    sudo reboot
    ```
**Now you should see TiberOS booting on your Host machine**

## 2. Verify Graphics SR-IOV and partition on Host machine booted with TiberOS Image

1. SSH to TiberOS, default credentials
   - Username : guest
   - Password : intel@123

2. Ensure SR-IOV is enabled
    ```sh
    sudo dmesg | grep -i i915
    ```
    - Output:
        ```sh
        [6.246297] i915 0000:00:02.0: [drm] Found ALDERLAKE_P/RPL-P (device ID a7a0) display version 13.00 stepping E0
        [6.246313] i915 0000:00:02.0: Running in SR-IOV PF mode
        [6.246876] i915 0000:00:02.0: [drm] VT-d active for gfx access
        [6.247005] i915 0000:00:02.0: vgaarb: deactivate vga console
        [6.247045] i915 0000:00:02.0: [drm] Using Transparent Hugepages
        [6.247542] i915 0000:00:02.0: vgaarb: VGA decodes changed: olddecodes=io+mem,decodes=io+mem:owns=io+mem
        [6.251233] i915 0000:00:02.0: [drm] Finished loading DMC firmware i915/adlp_dmc.bin (v2.20)
        [6.257232] i915 0000:00:02.0: [drm] GT0: GuC firmware i915/adlp_guc_70.bin version 70.36.0
        [6.257236] i915 0000:00:02.0: [drm] GT0: HuC firmware i915/tgl_huc.bin version 7.9.3
        [6.270238] i915 0000:00:02.0: [drm] GT0: HuC: authenticated for all workloads
        [6.270998] i915 0000:00:02.0: [drm] GT0: GUC: submission enabled
        [6.270999] i915 0000:00:02.0: [drm] GT0: GUC: SLPC enabled
        [6.271507] i915 0000:00:02.0: [drm] GT0: GUC: RC enabled
        [6.326610] [drm] Initialized i915 1.6.0 for 0000:00:02.0 on minor 0
        [6.385387] fbcon: i915drmfb (fb0) is primary device
        [6.544458] i915 0000:00:02.0: [drm] fb0: i915drmfb frame buffer device
        [6.573481] i915 0000:00:02.0: 7 VFs could be associated with this PF
        ```

3.  Verify the partitions, *Below Ex. is for 1TB NVME storage*
    ```sh
    lsblk
    ```
    - Output
        ```sh
        NAME        MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
        nvme0n1     259:0    0 931.5G  0 disk
        ├─nvme0n1p1 259:1    0   299M  0 part /boot/efi
        ├─nvme0n1p2 259:2    0 453.8G  0 part /
        └─nvme0n1p3 259:3    0 477.4G  0 part /var/edge-node/pua
                                              /etc/intel_manageability.conf_bak
                                              /etc/intel_manageability.conf
                                              /etc/dispatcher.environment
                                              /var/log/inbm-update-log.log
                                              /var/log/inbm-update-status.log
                                              /var/lib/dispatcher
                                              /etc/intel-manageability
                                              /var/cache/manageability
                                              /var/intel-manageability
                                              /var/lib/rancher
                                              /etc/default
                                              /etc/lvm/backup
                                              /etc/lvm/archive
                                              /etc/kubernetes
                                              /etc/cni
                                              /etc/netplan
                                              /etc/rancher
                                              /etc/sysconfig
                                              /etc/cloud
                                              /etc/udev
                                              /etc/systemd
                                              /etc/ssh
                                              /etc/pki
                                              /etc/machine-id
                                              /etc/intel_edge_node
                                              /etc/hosts
                                              /etc/environment
                                              /etc/fstab
                                              /home
                                              /opt
        ```

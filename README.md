# applications.virtualization.maverickflats-tiberos-itep
> Host OS (Tiber) should support Graphics SRIOV.

Clone all the submodules with git command: 
```
git submodule update --init --recursive
```

# Hardware
1. Intel NUC 12 Pro Kit NUC12WSHi7 Mini PC with 12th Gen Core i7-1260P Processor (12 Cores 16 Threads 4.70GHz 18MB Cache Intel Iris Xe Graphics) with 64GB DDR4 RAM, 1TB M.2 SSD, 2.5GbE LAN, Wi-Fi 6E, Bluetooth 5.3, 2x Thunderbolt 4 ports
2. Intel NUC 12 Pro Kit NUC12WSHi5 Mini PC with 12th Gen Core i5-1240P Processor (12 Cores 16 Threads 4.40GHz 12MB Cache Intel Iris Xe Graphics) with 32GB DDR4 RAM, 1TB M.2 SSD, 2.5GbE LAN, Wi-Fi 6E, Bluetooth 5.3
3. ASUS ExpertCenter PN64-E1, Ultra-compact mini PC with 13th Gen Intel® Core™ mobile processor with Intel® Iris Xe Graphics, supports quad displays and 4K resolution, 2x PCIe® Gen4 x4 M.2 NVMe® SSD, 2.5 Gb LAN, WiFi 6E

# Steps:
  1. System Bring Up
     - [Flashing SR-IOV Enabled TiberOS and Creating Partitions](tiber/tiber_flash_partition.md)
     - [Userspace setup for Maverick Flats](tiber/tiber_mf_setup.md)
  2. [Installing Kubernetes](tiber/kubevirt_gfx_guide.md)
  3. [Install Intel custom Kubevirt (includes SR-IOV patched QEMU)](tiber/kubevirt_gfx_guide.md#13-build-and-install-customized-kubevirt-for-maverick-flats)
  4. [Install Device Plugin](https://github.com/intel-innersource/applications.virtualization.maverickflats-deviceplugin-itep?tab=readme-ov-file#deployment)
  5. Deployment Package - includes Sidecar and VM deployment Helm charts.
     - [Discrete VM Helm chart to deployment](deployment/discrete/discrete.md)
     - Single Helm deployment - In Progress

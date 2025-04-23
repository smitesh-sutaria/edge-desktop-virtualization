# applications.virtualization.maverickflats-tiberos-itep
> Host OS (Tiber) should support Graphics SRIOV.

Clone all the submodules with git command: 
```
git submodule update --init --recursive
```

# Steps:
  1. System Bring Up
     - [Flashing SR-IOV Enabled TiberOS and Creating Partitions](tiber/tiber_flash_partition.md)
     - [Userspace setup for Maverick Flats](tiber/tiber_mf_setup.md)
  2. [Installing Kubernetes](tiber/kubevirt_gfx_guide.md)
  3. [Install Intel custom Kubevirt (includes SR-IOV patched QEMU)](tiber/kubevirt_gfx_guide.md#install-intel-built-kubevirt)
  4. [Install Device Plugin](https://github.com/intel-innersource/applications.virtualization.maverickflats-deviceplugin-itep?tab=readme-ov-file#deployment)
  5. Deployment Package - includes Sidecar and VM deployment Helm charts.
     - [Discrete VM Helm chart to deployment](deployment/discrete/discrete.md)
     - Single Helm deployment - In Progress

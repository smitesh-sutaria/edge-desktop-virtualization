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
  2. [Installing Kubernetes](Link)
  3. [Install Intel Kubevirt (includes SR-IOV patched QEMU)](https://github.com/intel-innersource/applications.virtualization.maverickflats-kubevirt-itep)
  4. [Install Device Plugin](https://github.com/intel-innersource/applications.virtualization.maverickflats-deviceplugin-itep?tab=readme-ov-file#deployment)
  5. Deployment Package - includes Sidecar and VM deployment Helm charts.
     - [Discrete VM Helm chart to deployment](deployment/discrete/discrete.md)
     - Single Helm deployment - In Progress

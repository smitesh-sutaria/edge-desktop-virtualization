# applications.virtualization.maverickflats-tiberos-itep
> Host OS (Tiber) should support gfx-SRIOV.

Clone all the submodules with git command: 
```
git submodule update --init --recursive
```

# Steps to bring:
  1. [System Bring Up](tiber/tiber_flash_partition.md)
  2. [Installing Kubernetes](Link)
  3. [Install Intel Kubevirt (includes patched qemu)](Link)
  4. [Install Device Plugin](https://github.com/intel-innersource/applications.virtualization.maverickflats-deviceplugin-itep?tab=readme-ov-file#deployment)
  5. Deployment Package - includes Sidecar and VM deployment Helm charts.
     1. [Discrete VM Helm chart to deployment](deployment/discrete/discrete.md)
     2. [Single Helm deployment](Link)

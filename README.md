# applications.virtualization.maverickflats-tiberos-itep
> Host OS (Tiber) should support gfx-SRIOV.

Clone all the submodules with git command: 
```
git submodule update --init --recursive
```

## Different components:
  1. custom intel kubevirt. (includes patched qemu)
  2. device plugins.
  3. sidecar hook to pass required qemu parameters.
  4. deployment package.

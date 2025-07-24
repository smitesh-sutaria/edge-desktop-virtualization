<!-- Copyright (C) 2025 Intel Corporation -->
# Building EMT ISO with Desktop Virtualization (graphics SR-IOV)

## Using Standalone Build Script (Automated)

> The pre-reuisite : Ubuntu 22.04 or Ubuntu 24.04

### Run the script as sudo
```sh
sudo ./build_idv_iso.sh
```
### ISO file will be generated in the same path

<p align="left">
<img align="center" width=50% height=50% src="docs/emt-idv-iso-out.png" >
</p>
<p align="center">
<em></em>
</p>

### Refer the demo below


## Manual Steps

The image configuration is part of this repo [here](./idv.json)

### Pre-requisite

[Build Requirements](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/toolkit/docs/building/prerequisites-ubuntu.md#build-requirements-on-ubuntu)

> It is recommended to built against a stable/release tag.

#### Step 1: clone the EMT repo
```sh
git clone https://github.com/open-edge-platform/edge-microvisor-toolkit
```
#### Step 2: Checkout the tag
```sh
git checkout tags/<tag_name>
```
#### Step 3: Copy the idv.json to edge-microvisor-toolkit/toolkit/imageconfigs/
```sh
cp idv.json edge-microvisor-toolkit/toolkit/imageconfigs/
```
#### Step 4: Build the tools
```sh
cd edge-microvisor-toolkit/toolkit
sudo make toolchain REBUILD_TOOLS=y
```
#### Step 5: Build the ISO for desktop virtualization (IDV) 
```sh
sudo make iso -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n CONFIG_FILE=./imageconfigs/idv.json
```

### Troubleshoot

### Clean build

> For re-building with any other tags, its recommended to start clean and repeat above Steps 1 to 5.
> ```sh
> cd edge-microvisor-toolkit
> sudo make -C toolkit clean
> ```

#### Working with Proxies

> If you are behind proxies and have them set, use -E option with all make commands
> For ex :
> ```
> sudo -E make toolchain REBUILD_TOOLS=y
> sudo -E make iso -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n CONFIG_FILE=./imageconfigs/idv.json
> ```

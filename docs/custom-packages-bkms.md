# To create custom packages in EMT here are some indentified BKMs

## Build Package 

Build package to be included as a part of EMT image 

**Step 1: create a package directory under SPECS directory**
```bash
cd edge-microvisor-toolkit
mkdir -p SPECS/<custom-package-name>
```
**Step 2: Create a spec file for the package**
```bash
touch <custom-package-name>.spec 
```
Here is a link to [how to create a spec file](https://www.redhat.com/en/blog/create-rpm-package)

**Step 3: Create a signature for source files**
Create a signature for source files that needs to be include during the creation of spec file 
it could be a tar file, scripts, binaries, etc. keep all of the required sources under 
SPEC/custom-package-name
```bash 
#get the sha256sum 
sha256sum <source0file>
sha256sum <source1file>
.
.
```
Update the sha256sum in the json file
```json
# name of file :<custom-package-name.signature>.json
{
  "source0file": "source0file sha256sum",
  "source1file": "source1file sha256sum"
}
```

**Step 3: build the package**
```bash 
cd toolkit
sudo -E make build-packages SRPM_PACK_LIST=<custom-package-name> SRPM_FILE_SIGNATURE_HANDLING=update

#this will create <custom-package-name>.rpm which the image will fetch during the build 
```

## Build Image 

After the rpm is create now we can include to the build time of the image 

**Step 1: create a custom package list**
create a custom package list under edge-microvisor-toolkit/toolkit/imageconfigs/packagelists
```json
#name of file: <custom-packages>.json
{
    "packages": [
        "<custom-package>"
    ]
}
# you can include other necessary dependencies here if needed
```
**Step 2: Edit the edge-image-mf-dev.json file. Add the custom packagelist**
```json
...
"PackageLists": [
  "packagelists/core-packages-image-systemd-boot.json",
  "packagelists/ssh-server.json",
  "packagelists/virtualization-host-packages.json",
  "packagelists/agents-packages.json",
  "packagelists/tools-tinker.json",
  "packagelists/persistent-mount-package.json",
  "packagelists/fde-verity-package.json",
  "packagelists/selinux-full.json",
  "packagelists/intel-gpu-base.json",
  "packagelists/os-ab-update.json",
  "packagelists/<custom-packages>.json" #here, a new entry
],
...
```
**Step 3: Build the image**
```bash
sudo -E make image -j8 REBUILD_TOOLS=y REBUILD_PACKAGES=n CONFIG_FILE=imageconfigs/edge-image-mf-dev.json
```

## References
- [building-howto](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/docs/developer-guide/get-started/building-howto.md#example-2-adding-a-new-rpm-package)
- [building options](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/toolkit/docs/building/building.md)
# Building Kubevirt and patching QEMU

This document provides steps, related to 
- Patching QEMU with Intel GPU SR-IOV patches and replacing the version of QEMU in Kubevirt
- Enabling Kubevirt to local support with GTK

## Overview

The following will be captured in this document:

- Steps to patch QEMU source code on Ubuntu host using the SR-IOV patches
- Steps to build QEMU within a Centos 9 image container
- Steps to patch qemu.conf in Kubevirt source code
- Steps to patch the Kubevirt 
- Steps to build the Kubevirt images and manifests using custom QEMU binary

> [!Note]
> OS and QEMU version provided in default Kubevirt virt-launcher image is

```shell
cat /etc/os-release
NAME="CentOS Stream" VERSION="9"

/usr/libexec/qemu-kvm --version
QEMU emulator version 9.0.0 (qemu-kvm-9.0.0-10.el9)
```

## 1. Build System Setup

> [!Note]
> Build system where the Kubevirt and QEMU build performed is Ubuntu 22.04 LTS

1. Install podman and setup registry
    ```sh
    sudo apt -y install podman

    podman run -d -p 5000:5000 --name local-registry registry:2
    ```


## 2. Enabling Kubevirt with local GTK display

1. Clone the repo:

    ```sh
    mkdir ~/workspace

    cd ~/workspace

    git clone https://github.com/kubevirt/kubevirt.git

    cd kubevirt
    ```

2. Add the following packages to `centos_main` section of `hack/rpm_deps.sh`
   
    ```diff
    @ -42,6 +42,63 @@ centos_main="
    acl
    curl-minimal
    vim-minimal
    + gtk3-devel
    + libjpeg-turbo
    + openjpeg2
    + libjpeg-turbo-devel
    + libproxy
    + libproxy-webkitgtk4
    + xdg-dbus-proxy
    + SDL2
    + SDL2-devel
    + libxdp-devel
    + mesa-libgbm
    + mesa-libgbm-devel
    + gdk-pixbuf2
    + gdk-pixbuf2-modules
    + gdk-pixbuf2-devel
    + cairo
    + vulkan-tools
    + vulkan-loader
    + cairo-gobject
    + cairo-devel
    + cairo-gobject-devel
    + gdk-pixbuf2
    + gdk-pixbuf2-modules
    + gdk-pixbuf2-devel
    + vte-profile
    + vte291
    + vte291-devel
    + libX11-xcb
    + xorg-x11-fonts-ISO8859-1-100dpi
    + libX11-common
    + libX11
    + xorg-x11-proto-devel
    + libX11-devel
    + sound-theme-freedesktop
    + alsa-lib
    + pulseaudio-libs
    + pulseaudio-libs-glib2
    + pipewire-pulseaudio
    + pulseaudio-libs-devel
    + pipewire-libs
    + pipewire
    + pipewire-jack-audio-connection-kit-libs
    + pipewire-jack-audio-connection-kit
    + pipewire-alsa
    + pipewire-pulseaudio
    + brlapi
    + brlapi-devel
    + fuse3-libs
    + fuse3-devel
    + libiscsi
    + libblkio
    + libblkio-devel
    + librbd1
    + librbd-devel
    + librados2
    + librados-devel
    + libradospp-devel
    "
    centos_extra="
    coreutils-single
    ```

3. Export the location of the docker registry and build tag (local docker registry in this case)

    ```sh
    export DOCKER_PREFIX=localhost:5000
    export DOCKER_TAG=mybuild
    ```

4. Build Kubevirt
   
   ```sh
   make rpm-deps
   ```

**Once the build completes successfully, Kubevirt is has been enabled with GTK dependent libraries**


## 3. Patching QEMU

The SRIOV patches to QEMU are based on QEMU 8

1. Download the the SR-IOV patches for QEMU

    ```sh
    cd ~/workspace

    wget -N --no-check-certificate https://download.01.org/intel-linux-overlay/ubuntu/pool/main/q/qemu/qemu_8.2.1+ppa1-noble9.debian.tar.xz

    tar -xvf qemu_8.2.1+ppa1-noble9.debian.tar.xz
    ```

2. Download QEMU 8:

    ```sh
    wget -N --no-check-certificate https://download.qemu.org/qemu-8.0.0.tar.xz

    tar -xvf qemu-8.0.0.tar.xz

    cd qemu-8.0.0
    ```

3. Copy the patches from `qemu_8.2.1+ppa1-noble9.debian` to qemu-8.0.0 directory:

    ```sh
    cp  -r ../qemu_8.2.1+ppa1-noble9.debian/debian/patches/sriov/ .
    ```

4. Apply patches:

    ```sh
    git apply ./sriov/*.patch

    cd ~/workspace
    ```

### 3.1 Creating CentOS 9 containerized environment

QEMU is built in Centos 9 container environment to ensure compatible with the Centos 9 based container image for virt-launcher.

The original idea to build within the Centos container comes from this [link](https://github.com/alicefr/kubevirt-debug/tree/main/build-with-custom-files#build-kubevirt-with-qemu-from-source-code)

1. Generate the Centos 9 image to be used for QEMU build environment

    ```sh
    cd qemu-8.0.0

    ./tests/lcitool/libvirt-ci/bin/lcitool --data-dir ./tests/lcitool dockerfile centos-stream-9 qemu > Dockerfile.centos-stream9
    ```

2. Patch `Dockerfile.centos-stream9` to include missing dependencies

    ```sh
    vim Dockerfile.centos-stream9
    ```

    ```diff
            zlib-devel \
            zlib-static \
    -        zstd &&
    +        zstd \
    +        libslirp-devel \
    +        liburing-devel \
    +        libbpf-devel \
    +        libblkio-devel && \
        dnf autoremove -y && \
        dnf clean all -y && \
        rpm -qa | sort > /packages.txt && \
    ```
    > [!Note] 
    > All expected dependencies for QEMU build in this example are there after the patch - if QEMU will be built with other flags it is possible that some dependencies may be missing and will need to be added to this Dockerfile.

3. Build the Centos 9 image

    ```sh
    podman build -t qemu_build:centos-stream9 -f Dockerfile.centos-stream9 .
    ```

4. Starting Centos 9 build environment

    Start the Centos 9 environment container from `qemu-8.0.0` directory. This allows the whole content of the QEMU source to be inside the environment.

    ```sh
    cd qemu-8.0.0

    podman run -ti \
        -v $(pwd):/src:Z \
        -w /src  \
        --security-opt label=disable \
        qemu_build:centos-stream9
    ```

### 3.2 Building QEMU inside the Centos 9 container

1. Configure QEMU build

    ```sh
    [root@<container> src]# ./configure --prefix=/usr --enable-kvm --disable-xen --enable-libusb --enable-debug-info --enable-debug  --enable-sdl --enable-vhost-net --enable-spice --disable-debug-tcg  --enable-opengl  --enable-gtk  --enable-virtfs --target-list=x86_64-softmmu --audio-drv-list=pa --firmwarepath=/usr/share/qemu-firmware:/usr/share/ipxe/qemu:/usr/share/seavgabios:/usr/share/seabios:/usr/share/qemu-kvm/ --disable-spice
    ```

    > [!Note] 
    > Certain flags were added to the default configure command to enable **GTK on local display**:
    > * spice is disabled with --disable-spice (spice is not natively supported in RHEL9 image)
    > * Firmware paths necessary for the operation of VM within the virt-launcher pods are added `--firmwarepath=/usr/share/qemu-firmware:/usr/share/ipxe/qemu:/usr/share/seavgabios:/usr/share/seabios:/usr/share/qemu-kvm/`


2. Build QEMU

    ```sh
    [root@<container> src]# cd build
    [root@<container> src]# ninja
    [root@<container> src]# ninja install
    ```

3. Exit the CentOS build container environment - the built artifacts will be present in the `qemu-8.0.0` directory on host.

    ```sh
    [root@<container> src]# exit
    ```
    On build system:
    ```sh
    ls -la build/qemu-system-x86_64
    -rwxr-xr-x 1 user user 47M Feb 13 14:27 build/qemu-system-x86_64

    sha256sum build/qemu-system-x86_64
    13c2760bf012a8011ddbe0c595ec3dca24249debe32bc4d1e338ec8538ad7453 build/qemu-system-x86_64
    ```

## 4. Patch Kubevirt with SR-IOV patched QEMU binary

1. Create a build directory to store new QEMU binary and directory for shared libraries

    ```sh
    cd ~/workspace/kubevirt

    mkdir build
    ```

2. Copy the QEMU binary built from the `qemu-8.0.0` directory into `kubevirt`'s directory `build`.

    ```sh
    cp ../qemu-8.0.0/build/qemu-system-x86_64 build/qemu-system-x86_64
    ```

3. Obtain the `SHA` hash number of the QEMU binary

    ```sh
    sha256sum ./build/qemu-system-x86_64
    13c2760bf012a8011ddbe0c595ec3dca24249debe32bc4d1e338ec8538ad7453 ./build/qemu-system-x86_64
    ```

5. Patch the top level `WORKSPACE` file in top level `kubevirt` directory.
   Replace the `sha256` with the one of the new QEMU binary.
   This will point the build to where the local binary is.

    ```sh
    vim WORKSPACE
    ```

    ```diff
    @@ -152,6 +152,15 @@ http_file(
        ],
    )

    +http_file(
    +    name = "custom-qemu",
    +    downloaded_file_path = "qemu-kvm",
    +    sha256 = "13c2760bf012a8011ddbe0c595ec3dca24249debe32bc4d1e338ec8538ad7453",
    +    urls = [
    +        "file:///root/go/src/kubevirt.io/kubevirt/build/qemu-system-x86_64",
    +    ],
    +)
    +
    http_archive(
        name = "bazeldnf",
        sha256 = "fb24d80ad9edad0f7bd3000e8cffcfbba89cc07e495c47a7d3b1f803bd527a40",
    ```

5. Patch the `cmd/virt-launcher/BUILD.bazel` file.
   This will point to where the new QEMU binary in virt-launcher container image build

    ```sh
    vim cmd/virt-launcher/BUILD.bazel
    ```

    ```diff
    --- a/cmd/virt-launcher/BUILD.bazel
    +++ b/cmd/virt-launcher/BUILD.bazel
    @@ -165,6 +178,15 @@ pkg_tar(
        owner = "0.0",
    )

    +pkg_tar(
    +    name = "custom-qemu-build",
    +    srcs = ["@custom-qemu//file"],
    +    mode = "0755",
    +    owner = "0.0",
    +    package_dir = "/usr/libexec",
    +    visibility = ["//visibility:public"],
    +)
    +
    container_image(
        name = "version-container",
        directory = "/",
    @@ -189,6 +211,8 @@ container_image(
                ":passwd-tar",
                ":nsswitch-tar",
                ":qemu-kvm-modules-dir-tar",
    +            ":custom-qemu-build",
                "//rpm:launcherbase_x86_64",
            ],
        }),
    ```

6. Patch the `qemu.conf` file

    To patch the custom qemu.conf edit the local configuration at `cmd/virt-launcher/qemu.conf`
    Any additional qemu configurations should be added to this file.

    ```sh
    vim cmd/virt-launcher/qemu.conf
    ```

    ```diff
    --- a/cmd/virt-launcher/qemu.conf
    +++ b/cmd/virt-launcher/qemu.conf
    @@ -8,3 +8,7 @@ dynamic_ownership = 1
    remember_owner = 0
    namespaces = [ ]
    cgroup_controllers = [ ]
    +security_default_confined = 0
    +seccomp_sandbox = 0
    +cgroup_device_acl = ["/dev/null", "/dev/full", "/dev/zero","/dev/random", "/dev/urandom", "/dev/ptmx", "/dev/kvm", "/dev/dri/card0"]
    +security_driver = []
    ```

## Build Kubevirt images and manifests

1. Export the location of the docker registry and build tag (local docker registry in this case)

    ```sh
    export DOCKER_PREFIX=localhost:5000
    export DOCKER_TAG=mybuild
    ```

2. Build the Kubevirt images from `kubevirt` top level directory

    ```sh
    make all
    make bazel-build-images
    ```

3. Push the images to the local Docker registry

    ```sh
    make push
    ***
    2025/02/13 15:43:38 Successfully pushed Docker image to localhost:5000/network-passt-binding:mybuild
    BUILD_ARCH= DOCKER_PREFIX=localhost:5000 DOCKER_TAG=mybuild hack/push-container-manifest.sh
    ```

4. Build manifests referencing the image locations

    ```sh
    make manifests
    ```

5. To install Kubevirt
    ```sh
    kubectl apply -f _out/manifests/release/kubevirt-operator.yaml
    kubectl apply -f _out/manifests/release/kubevirt-cr.yaml
    ```

6.  Verify Deployment
    ```sh
    kubectl get all -n kubevirt

    NAME                                   READY   STATUS    RESTARTS      AGE
    pod/virt-api-999875d56-4dvsc           1/1     Running   6 (18d ago)   19d
    pod/virt-controller-546cb985cd-f4zns   1/1     Running   5 (18d ago)   19d
    pod/virt-controller-546cb985cd-kxmsr   1/1     Running   5 (18d ago)   19d
    pod/virt-handler-s4m9j                 1/1     Running   7 (15d ago)   19d
    pod/virt-operator-6459bcf8c6-vxbqx     1/1     Running   6 (18d ago)   19d
    pod/virt-operator-6459bcf8c6-xhktx     1/1     Running   6 (18d ago)   19d

    NAME                                  TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
    service/kubevirt-operator-webhook     ClusterIP   10.43.86.170   <none>        443/TCP   19d
    service/kubevirt-prometheus-metrics   ClusterIP   None           <none>        443/TCP   19d
    service/virt-api                      ClusterIP   10.43.68.37    <none>        443/TCP   19d
    service/virt-exportproxy              ClusterIP   10.43.189.94   <none>        443/TCP   19d

    NAME                          DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
    daemonset.apps/virt-handler   1         1         1       1            1           kubernetes.io/os=linux   19d

    NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/virt-api          1/1     1            1           19d
    deployment.apps/virt-controller   2/2     2            2           19d
    deployment.apps/virt-operator     2/2     2            2           19d

    NAME                                         DESIRED   CURRENT   READY   AGE
    replicaset.apps/virt-api-6676df49cc          0         0         0       19d
    replicaset.apps/virt-api-999875d56           1         1         1       19d
    replicaset.apps/virt-controller-546cb985cd   2         2         2       19d
    replicaset.apps/virt-controller-54c7869f6c   0         0         0       19d
    replicaset.apps/virt-operator-6459bcf8c6     2         2         2       19d

    NAME                            AGE   PHASE
    kubevirt.kubevirt.io/kubevirt   19d   Deployed
    ```

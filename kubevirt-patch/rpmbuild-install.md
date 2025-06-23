# Building RPM and Installing Kubevirt

### Prerequisites
System should be installed with rpm build environment.

## Build

1. Download and copy the [Kubevirt TAR](https://github.com/open-edge-platform/edge-desktop-virtualization/releases/download/pre-release-v0.1/intel-idv-kubevirt-v0.1.tar.gz) and [Device-Plugin TAR](https://github.com/open-edge-platform/edge-desktop-virtualization/releases/download/pre-release-v0.1/intel-idv-device-plugin-v0.1.tar.gz) to `SOURCES`
2. Copy the desktop-virtualization-k3s.spec to `SPECS`
3. Build the RPM
   ```sh
   rpmbuild -ba desktop-virtualization-k3s.spec
   ```
4. Once build completes, RPM file will be present in `RPMS/x86_64` and `SRPMS`
5. Copy the file `RPMS/x86_64/desktop-virtualization-k3s-v1.0-1.x86_64.rpm` to host system

## Install
1. Install the RPM on host
   ```sh
    sudo rpm -i desktop-virtualization-k3s-v1.0-1.x86_64.rpm
   ```
2. Once after install completes, Kubevirt and Device-plugin tar files will be present in `/usr/share/desktop-virtualization-k3s/`
   ```sh
   ls -la /usr/share/desktop-virtualization-k3s/
   total 624520
   drwxr-xr-x.   2 root root      4096 Jun 12 06:21 .
   drwxr-xr-x. 107 root root      4096 Jun 12 06:21 ..
   -rw-rw-r--.   1 root root  14397311 Jun 10 09:51 dv-device-plugin.tar.gz
   -rw-rw-r--.   1 root root 625096727 Jun 10 09:51 kubevirt.tar.gz
   ```

### Load Kubevirt and Device-Plugin containers to k3s containerd

3.  Extract TAR files
    ```sh
    mkdir -p ~/display-virtualization

    tar -xzvf kubevirt.tar.gz ~/display-virtualization
    tar -xzvf dv-device-plugin.tar.gz ~/display-virtualization

    cd ~/display-virtualization
    ```
4.  Import the images into the container runtime
    ```sh
    sudo ctr -n k8s.io images import sidecar-shim.tar 
    sudo ctr -n k8s.io images import virt-api.tar
    sudo ctr -n k8s.io images import virt-controller.tar
    sudo ctr -n k8s.io images import virt-handler.tar
    sudo ctr -n k8s.io images import virt-launcher.tar
    sudo ctr -n k8s.io images import virt-operator.tar

    sudo ctr -n k8s.io images import device-plugin.tar
    sudo ctr -n k8s.io images import busybox.tar
    ```
5.  Verify the images are imported correctly
    ```sh
    sudo crictl images | grep localhost

    localhost:5000/sidecar-shim                           v1.5.0_DV           c48d79a700926       51.5MB
    localhost:5000/virt-api                               v1.5.0_DV           025a39d7f7504       28.6MB
    localhost:5000/virt-controller                        v1.5.0_DV           d1cb23d032aa0       27.9MB
    localhost:5000/virt-handler                           v1.5.0_DV           a9bd1a37e2e0c       90.7MB
    localhost:5000/virt-launcher                          v1.5.0_DV           c69ddc6b90387       403MB
    localhost:5000/virt-operator                          v1.5.0_DV           99462ddb3a866       39.8MB
    localhost:5000/mf-device-plugin                       v1                  156ba1fcaf549       21.3MB
    localhost:5000/busybox                                latest              ff7a7936e9306       2.21MB
    ```
6.  Deploy Kubevirt and Device Plugin
    ```sh
    kubectl apply -f kubevirt-operator.yaml
    kubectl apply -f kubevirt-cr.yaml
    kubectl apply -f device-plugin.yaml
    ```
7.  Verify Deployment
    ```sh
    kubectl get all -A

    NAMESPACE     NAME                                          READY   STATUS    RESTARTS      AGE
    .
    .
    kube-system   pod/device-plugin-q2c2n                       1/1     Running   0             10d
    kubevirt      pod/virt-api-6c66767447-tvqwz                 1/1     Running   0             8d
    kubevirt      pod/virt-controller-599f9b4d86-ffv2b          1/1     Running   0             8d
    kubevirt      pod/virt-controller-599f9b4d86-pt5rn          1/1     Running   0             8d
    kubevirt      pod/virt-handler-hbtsj                        1/1     Running   0             8d
    kubevirt      pod/virt-operator-69cb894b4c-djrzh            1/1     Running   0             8d
    kubevirt      pod/virt-operator-69cb894b4c-jc8sk            1/1     Running   0             8d
    .
    .
    .
    NAMESPACE   NAME                            AGE   PHASE
    kubevirt    kubevirt.kubevirt.io/kubevirt   9d    Deployed
    .
    .
    ```
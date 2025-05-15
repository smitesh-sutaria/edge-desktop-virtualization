# Installing Kubernetes and Enabling Graphics SR-IOV for Kubevirt

Clone the repo to Tiber Host system
```sh
git clone https://github.com/intel-innersource/applications.virtualization.maverickflats-tiberos-itep.git
```

[Intel® Graphics SR-IOV Enablement Toolkit](https://github.com/intel/kubevirt-gfx-sriov) is used to enable Graphics SR-IOV and is cloned to folder `kubevirt-gfx-sriov`, credits to all contibutors of it.\
Minor changes has been made to it to work on TiberOS and to support Intel's custom Kubevirt

## 1. Installation
Below steps are customized for TiberOS and derived from [Manual Install](../kubevirt-gfx-sriov/docs/manual-install.md) of kubevirt-gfx-sriov, for more details follow the same link

### 1.1 Install K3s

This step will setup a single node cluster where the host function as both the server/control plane and the worker node.\
This step is only required if you don't already have a Kubernetes cluster setup that you can use

> [!Note]
> K3s is a lightweight Kubernetes distribution suitable for Edge and IoT use cases.

```sh
cd kubevirt-gfx-sriov

./scripts/setuptools.sh -ik
```

### 1.2 Enable Graphics VFs on boot

Add systemd service unit file to enable graphics VFs on boot.
```sh
sudo mkdir -p /opt/vm/scripts

sudo cp scripts/configvfs.sh /opt/vm/scripts/

sudo chmod +x /opt/vm/scripts/configvfs.sh

sudo cp systemd/gfx-virtual-func.service /etc/systemd/system/

sudo systemctl daemon-reload

sudo systemctl start gfx-virtual-func.service

sudo systemctl enable gfx-virtual-func.service

sudo reboot
```

#### 1.2.1 Check the `configvfs.sh` log and `gfx-virtual-func.service` daemon status for any error
```sh
sudo systemctl status gfx-virtual-func.service
```
-   Output:
    ```sh
    ● gfx-virtual-func.service - Intel Graphics SR-IOV Virtual Function Manager
        Loaded: loaded (/etc/systemd/system/gfx-virtual-func.service; enabled; preset: disabled)
        Drop-In: /usr/lib/systemd/system/service.d
                └─10-timeout-abort.conf
        Active: active (exited) since Fri 2025-04-04 16:49:47 UTC; 2 weeks 4 days ago
    Main PID: 1833 (code=exited, status=0/SUCCESS)
            CPU: 205ms

    Apr 04 16:49:47 EdgeMicrovisorToolkit systemd[1]: Starting gfx-virtual-func.service - Intel Graphics SR-IOV Virtual Function Manager...
    Apr 04 16:49:47 EdgeMicrovisorToolkit bash[1833]: Device: /sys/bus/pci/devices/0000:00:02.0
    Apr 04 16:49:47 EdgeMicrovisorToolkit bash[1833]: Total VF: 7
    Apr 04 16:49:47 EdgeMicrovisorToolkit bash[1833]: ID: 0x8086 0xa7a0
    Apr 04 16:49:47 EdgeMicrovisorToolkit bash[1833]: VF enabled: 7
    Apr 04 16:49:47 EdgeMicrovisorToolkit systemd[1]: Finished gfx-virtual-func.service - Intel Graphics SR-IOV Virtual Function Manager.
    ```
-   Verify VFs enumarated on PCI bus
    ```sh
    lspci

    00:02.0 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.1 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.2 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.3 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.4 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.5 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.6 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    00:02.7 VGA compatible controller: Intel Corporation Raptor Lake-P [Iris Xe Graphics] (rev 04)
    ```

### 1.3 Install customized Kubevirt for Maverick-Flats

To build and install Kubevirt [refer](./kubevirt_dv_build_guide.md/#steps-to-build-intel-cutomized-kubevirt)

#### 1.3.1 Setup Quick Install from one-intel-edge-sandbox repository (for Intel Internal use only)

1.  Update `registry.yaml` to pull from `one-intel-edge-sandbox` registry
    ```sh
    sudo vi /etc/rancher/k3s/registries.yaml
    ```
    Add
    ```sh
    mirrors:
    "amr-registry-pre.caas.intel.com":
        endpoint:
        - "https://amr-registry-pre.caas.intel.com"
    ```

2.  Update Proxy for K3S
    ```sh
    sudo vi /etc/systemd/system/k3s.service.env
    ```
    Add `.intel.com` IP in `NO_PROXY`
    ```sh
    HTTPS_PROXY="http://proxy-dmz.intel.com:912"
    HTTP_PROXY="http://proxy-dmz.intel.com:911"
    NO_PROXY="localhost,::1,127.0.0.1,.intel.com"
    ```

3.  Copy the Certificate
    ```sh
    sudo cp tiber/certificates/harbor.crt /etc/pki/ca-trust/source/anchors/

    sudo update-ca-trust
    ```

4.  Restart K3S
    ```sh
    sudo systemctl restart k3s
    ```

#### 1.3.2 Install Kubevirt from one-intel-edge-sandbox (for Intel Internal use only)

1.  Install Kubevirt
    ```sh
    kubectl apply -f tiber/kubevirt/manifests/release/kubevirt-operator.yaml
    kubectl apply -f tiber/kubevirt/manifests/release/kubevirt-cr.yaml
    ```

    -   Verify Kubevirt Deployment
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

### 1.4 Install CDI
```sh
cd kubevirt-gfx-sriov

./scripts/setuptools.sh -iv
```

### 1.5 Enable Virt-Handler to discover Graphics VFs
Update KubeVirt custom resource configuration to enable virt-handler to discover graphics VFs on the host. All discovered VFs will be published as *allocatable* resource

**Update Graphics Device ID in `kubevirt-gfx-sriov/manifests/kubevirt-cr-gfx-sriov.yaml` if not found**
  - Read the Device ID of Intel Graphics Card from Host, Ex: for RPL
    ```sh
    $ cat /sys/devices/pci0000\:00/0000\:00\:02.0/device

    0xa7a0
    ```
  - Add the Device ID in `pciHostDevices` section
    ```yaml
    - pciVendorSelector: "8086:a7a0"
      resourceName: "intel.com/sriov-gpudevice"
      externalResourceProvider: false
    ```

Apply the YAML changes

```sh
cd kubevirt-gfx-sriov

kubectl apply -f manifests/kubevirt-cr-gfx-sriov.yaml
```

#### 1.5.1 Check for presence of `intel.com/sriov-gpudevices` resource

```sh
kubectl describe nodes
```
Output:
```sh
Capacity:
    intel.com/sriov-gpudevice:     7
Allocatable:
    intel.com/sriov-gpudevice:     7
Allocated resources:
    Resource                       Requests     Limits
    --------                       --------     ------
    intel.com/sriov-gpudevice      0            0
```
> [!Note] 
> Please wait for all virt-handler pods to complete restarts\
> The value of **Requests** and **Limits** will increase upon successful resource allocation to running pods/VMs

### 1.6 Install Device-Plugin

To build and install Device-plugin [refer](./kubevirt_dv_build_guide.md/#steps-to-build-device-plugin)

#### For Quick Install from one-intel-edge-sandbox repository (for Intel Internal use only)

Ensure [these steps](./kubevirt_gfx_guide.md/#131-setup-quick-install-from-one-intel-edge-sandbox-repository) are completed.

1. Install Device-Plugin from one-intel-edge-sandbox repository

    ```sh
    kubectl apply -f tiber/device-plugin/manifests/maverikflats-device-plugin.yaml
    ```
    or
    ```sh
    # Helm deployment

    cd tiber/device-plugin/helm/

    helm install device-plugin .
    ```

    -   To verify 
    ```sh
    kubectl describe nodes
    ```
    -   Output
        ```sh
        Capacity:
            intel.com/igpu:                 1k
            intel.com/udma:                 1k
            intel.com/usb:                  1k
            intel.com/vfio:                 1k
            intel.com/x11:                  1k
        Allocatable:
            intel.com/igpu:                 1k
            intel.com/udma:                 1k
            intel.com/usb:                  1k
            intel.com/vfio:                 1k
            intel.com/x11:                  1k
        Allocated resources:
            Resource                       Requests          Limits
            --------                       --------          ------
            intel.com/igpu                   0                 0
            intel.com/udma                   0                 0
            intel.com/usb                    0                 0
            intel.com/vfio                   0                 0
            intel.com/x11                    0                 0
        ```

## 2. Create Windows-10/11 Image

Refer [Installation](../kubevirt-gfx-sriov/docs/deploy-windows-vm.md#installation) and [Deployment](../kubevirt-gfx-sriov/docs/deploy-windows-vm.md#deployment) section

> [!Note]
Once after steps are complete, search and copy the `disk.qcow2` or `disk.img` from `/opt/vm/images/win10` which can be later used to boot VM from CDI

## 3. Uninstall

1. To uninstall all the components you can run the command below or you can specify which component to uninstall.

   *Note: Get help on `setuptools.sh` by running `setupstool.sh -h`*
   ```sh
   cd kubevirt-gfx-sriov

   ./scripts/setuptools.sh -u kvw
   ``` 
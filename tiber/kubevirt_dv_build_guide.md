# Build setup
>[!Note]
Kubevirt build setup is based on `Ubuntu 22.04 LTS`.

1.  Setting up local registry on Ubuntu build system which will be used as registry to pull Kubevirt for Tiber Host
    ```sh
    sudo apt-get -y install podman

    podman run -d -p 5000:5000 --name local-registry registry:2
    ```
    or
    ```sh
    docker run -d -p 5000:5000 --name registry registry:2.7
    ```

# Settings on Tiber Host to pull Kubevirt or Device-Plugin from build system registry

1.  Update the Registry for K3S to pull from build system

    Ex: If Localserver registry IP is 10.223.97.134:5000 `kubevirt-operator.yaml`, update that in `registries.yaml` and in `NO_PROXY` of `k3s.service.env`
    ```sh
    sudo vi /etc/rancher/k3s/registries.yaml
    ```
    Add
    ```sh
    mirrors:
    "10.190.167.198:5000":
        endpoint:
        - "http://10.190.167.198:5000"
    "10.223.97.134:5000":
        endpoint:
        - "http://10.223.97.134:5000"
    ```

2.  Update Proxy for K3S
    ```sh
    sudo vi /etc/systemd/system/k3s.service.env
    ```
    Add server IP in `NO_PROXY`
    ```sh
    HTTPS_PROXY="http://proxy-dmz.intel.com:912"
    HTTP_PROXY="http://proxy-dmz.intel.com:911"
    NO_PROXY="localhost,::1,127.0.0.1,.intel.com,10.190.167.198,10.223.97.134"
    ```

3.  Restart K3S
    ```sh
    sudo systemctl restart k3s
    ```

# Steps to build Intel cutomized Kubevirt

[Maverick-Flats-Kubevirt](https://github.com/intel-innersource/applications.virtualization.maverickflats-kubevirt-itep) version hosted in Intel-Innersource

1.  Clone the repo, build the Kubevirt, for detailed build steps [refer](https://github.com/intel-innersource/applications.virtualization.maverickflats-kubevirt-itep/blob/v1.5.0/docs/build-the-builder.md).
    ```sh
    git clone https://github.com/intel-innersource/applications.virtualization.maverickflats-kubevirt-itep.git

    cd applications.virtualization.maverickflats-kubevirt-itep

    export DOCKER_PREFIX=<build_system_ip>:5000   #Ex. localhost:5000
    export DOCKER_TAG=latest

    make all
    make bazel-build-images
    make push
    make manifests
    ```
## Install Kubevirt on Tiber Host
1.  To install Kubevirt
    Copy the `kubevirt-operator.yaml` and `kubevirt-cr.yaml` from `_out/manifests/release` of build system to TiberOS host machine.
    ```sh
    kubectl apply -f kubevirt-operator.yaml
    kubectl apply -f kubevirt-cr.yaml
    ```

2.  Verify Deployment
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

# Steps to build Device-plugin

1.  Clone the repo
    ```sh
    git clone https://github.com/intel-innersource/applications.virtualization.maverickflats-deviceplugin-itep.git

    cd applications.virtualization.maverickflats-deviceplugin-itep
    ```

2.  Build
    ```sh
    make docker-build

    make docker-push
    ```

## Install Device-Plugin on Tiber Host

1.  Copy `deploy` folder from build system to Tiber host
2.  Replace `localhost:5000` in `deploy/helm/values.yaml` and `deploy/manifests/maverikflats-device-plugin.yaml` with your Ubuntu build server IP
    ```sh
    kubectl apply -f tiber/device-plugin/manifests/maverikflats-device-plugin.yaml
    ```
    or
    ```sh
    # Helm deployment

    cd tiber/device-plugin/helm/

    helm install device-plugin .
    ```
    Verify deployment 
    ```sh
    kubectl describe nodes

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

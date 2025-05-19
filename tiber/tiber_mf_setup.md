# 1. Setup Hugepages
Set up Hugepages with pagesize 2MB, for 4 VMs with each VM RAM set to 12GB (48 GB total)

Create a service to set up these hugepages at boot time
```sh
sudo vi /etc/systemd/system/hugepages.service
```
Add
```
[Unit]
Description=Configure Hugepages
Before=k3s.service

[Service]
Type=oneshot
ExecStart=/bin/sh -c 'echo $(( 6 * 1024 * 4 )) | sudo tee /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages'

[Install]
WantedBy=multi-user.target
```

Enable and start the service, which will configure the hugepages and exit.
```
sudo systemctl daemon-reload
sudo systemctl enable hugepages.service
sudo systemctl start hugepages.service
```

Check that hugepages were configured
```
sudo cat /proc/meminfo | grep -i hugepages
```
```
HugePages_Total:   24576
HugePages_Free:    24576
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
```

# 2. Set USB permissions
To use USB peripherals connected to Host machine with Virtual machines, create a udev rule that will automatically give the user qemu access to them
```
sudo vi /etc/udev/rules.d/99-usb-qemu.rules
```
Add
```
SUBSYSTEM=="usb", MODE="0664", GROUP="qemu"
```
Apply changes
```
sudo udevadm control --reload-rules
```
Unplug and re-plug the USB devices you plan to attach to VMs, then check the permissions are set correctly:
```
ls -alR /dev/bus/usb/
```
```
...
crw-rw-r--. 1 root qemu 189, 8 May 15 22:14 009
...
```

# 3. Display setup for TiberOS

TiberOS boots with no GUI and prompts for user login, login using default credentials\
XSERVER is installed by default, make the below settings before starting X server

## 3.1 Disable DPMS and screen blanking on the X Window System

-   DPMS Disable
    ```sh
    sudo vi /usr/share/X11/xorg.conf.d/10-extensions.conf
    ```
    Add
    ```conf
    Section "Extensions"
        Option "DPMS" "false"
    EndSection
    ```

-   Disable Screen Blanking and Timeouts
    ```sh
    sudo vi /usr/share/X11/xorg.conf.d/10-serverflags.conf
    ```
    Add
    ```conf
    Section "ServerFlags"
        Option "StandbyTime" "0"
        Option "SuspendTime" "0"
        Option "OffTime"     "0"
        Option "BlankTime"   "0"
    EndSection
    ```


## 3.2 Create a service to autostart the X server
```sh
sudo vi /etc/systemd/system/x.service
```
Add
```
[Unit]
Description=Launch X server at startup
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/X

[Install]
WantedBy=graphical.target
```

Enable and start the service. You should now see a black screen on the monitors.
**NOTE: When you reboot the machine, you will end up with a black screen, because X is running. To access the console, try `control+alt+f3`.** To return to X, switch back with `control-alt-f2`. If X is not currently the active display, the VMs will not boot, and will error with "SyncVMI failed".
```
sudo systemctl daemon-reload
sudo systemctl enable x.service
sudo systemctl start x.service
```

### 3.2.1 Check Monitor's resolution and names of connected displays
Open SSH session to the Host system
```sh
DISPLAY=:0 xrandr
```
-   Output:
    ```sh
    Screen 0: minimum 320 x 200, current 7680 x 1080, maximum 16384 x 16384
    HDMI-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 521mm x 293mm
    1920x1080     60.00*+  50.00    59.94
    1600x1200     60.00
    1680x1050     59.88
    1400x1050     59.95
    1600x900      60.00
    1280x1024     75.02    60.02
    1440x900      59.90
    1280x960      60.00
    1280x800      59.91
    1152x864      75.00
    1280x720      60.00    50.00    59.94
    1024x768      75.03    70.07    60.00
    832x624       74.55
    800x600       72.19    75.00    60.32    56.25
    720x576       50.00
    720x480       60.00    59.94
    640x480       75.00    72.81    66.67    60.00    59.94
    720x400       70.08
    HDMI-2 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 527mm x 296mm
    1920x1080     60.00*+  50.00    59.94
    1680x1050     59.88
    1600x900      60.00
    1280x1024     75.02    60.02
    1280x800      59.91
    1152x864      75.00
    1280x720      60.00    50.00    59.94
    1024x768      75.03    60.00
    832x624       74.55
    800x600       75.00    60.32
    720x576       50.00
    720x480       60.00    59.94
    640x480       75.00    60.00    59.94
    720x400       70.08
    DP-1 connected 1920x1080+3840+0 (normal left inverted right x axis y axis) 521mm x 293mm
    1920x1080     60.00*+  74.92    50.00    59.94
    1600x1200     60.00
    1680x1050     59.95
    1400x1050     59.98
    1280x1024     75.02    60.02
    1440x900      59.89
    1280x960      60.00
    1280x800      59.81
    1152x864      75.00
    1280x720      60.00    50.00    59.94
    1440x576      50.00
    1024x768      75.03    70.07    60.00
    1440x480      60.00    59.94
    832x624       74.55
    800x600       72.19    75.00    60.32    56.25
    720x576       50.00
    720x480       60.00    59.94
    640x480       75.00    72.81    66.67    60.00    59.94
    720x400       70.08
    DP-2 disconnected (normal left inverted right x axis y axis)
    DP-3 connected 1920x1080+5760+0 (normal left inverted right x axis y axis) 521mm x 293mm
    1920x1080     60.00*+  74.99    50.00    59.94
    1600x1200     60.00
    1680x1050     59.88
    1400x1050     59.95
    1280x1024     75.02    60.02
    1440x900      59.90
    1280x960      60.00
    1280x800      59.91
    1152x864      75.00
    1280x720      60.00    50.00    59.94
    1440x576      50.00
    1024x768      75.03    70.07    60.00
    1440x480      60.00    59.94
    832x624       74.55
    800x600       72.19    75.00    60.32    56.25
    720x576       50.00
    720x480       60.00    59.94
    640x480       75.00    72.81    66.67    60.00    59.94
    720x400       70.08
    DP-4 disconnected (normal left inverted right x axis y axis)
    ```

## 3.3 Setup XDOTOOL in a container to scale applications to full-screen
> [!Note]
> This is needed if VM doesn't scale to full-screen after launching

### 3.3.1 Install Docker in MF image:
```sh
sudo dnf install -y moby-engine moby-cli ca-certificates
```
### 3.3.2 Create a http-proxy.conf file with below contents to let docker work in Intel network with proxy.
```sh
sudo mkdir -p /etc/systemd/system/docker.service.d/

sudo vi /etc/systemd/system/docker.service.d/http-proxy.conf
```
Add
```conf
[Service]
Environment="HTTPS_PROXY=http://proxy-dmz.intel.com:912"
Environment="HTTP_PROXY=http://proxy-dmz.intel.com:911"
Environment="NO_PROXY=localhost,127.0.0.0/8,172.16.0.0/20,192.168.0.0/16,10.0.0.0/8,.intel.com,intel.com"
```
Restart Docker
```sh
sudo systemctl enable docker.service

sudo systemctl daemon-reload

sudo systemctl start docker.service
```
### 3.3.3 Create alpine container and install XDOTOOL
```sh
sudo docker network create tiber-bridge-network

sudo docker run -it --rm --name tiber-resize-container \
    --privileged \
    -e DISPLAY=:0 \
    --ulimit nofile=65535:65535 \
    -v /tmp/.X11-unix:/tmp/.X11-unix \
    --network tiber-bridge-network \
    alpine:latest
```
#### 3.3.3.1 From container set proxy, install xdotool and rescale VM:
```sh
export https_proxy=http://proxy-dmz.intel.com:912

apk add xdotool
```
-   Output
    ```sh
    fetch https://dl-cdn.alpinelinux.org/alpine/v3.21/main/x86_64/APKINDEX.tar.gz
    fetch https://dl-cdn.alpinelinux.org/alpine/v3.21/community/x86_64/APKINDEX.tar.gz
    (1/14) Installing libxau (1.0.11-r4)
    (2/14) Installing libmd (1.1.0-r0)
    (3/14) Installing libbsd (0.12.2-r0)
    (4/14) Installing libxdmcp (1.1.5-r1)
    (5/14) Installing libxcb (1.16.1-r0)
    (6/14) Installing libx11 (1.8.10-r0)
    (7/14) Installing libxext (1.3.6-r2)
    (8/14) Installing libxinerama (1.1.5-r4)
    (9/14) Installing libxtst (1.2.5-r0)
    (10/14) Installing xkeyboard-config (2.43-r0)
    (11/14) Installing xz-libs (5.6.3-r0)
    (12/14) Installing libxml2 (2.13.4-r5)
    (13/14) Installing libxkbcommon (1.7.0-r1)
    (14/14) Installing xdotool (3.20211022.1-r1)
    Executing busybox-1.37.0-r12.trigger
    OK: 15 MiB in 29 packages
    ```
-   To resize Virtual Machine Window
    ```sh
    xdotool search --onlyvisible --name "vm1" windowsize 1920 1080
    ```

# 4. Install packages

```sh
sudo dnf install helm git
```

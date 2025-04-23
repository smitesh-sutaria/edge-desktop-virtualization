# Setup Hugepages
To setup Hugepages with pagesize 2048M, for 4 VMs with each VM RAM set to 12GB
```
sudo su

echo 24576 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
```

# Set USB permissions
To use USB peripherals connected to Host machine with Virtual machines, set the USB devices permission to user `qemu`
```
sudo chown -R qemu:root /dev/bus/usb/
```
> [!Note]
> This has to be done everytime USB device is hot-plugged

# Display setup for TiberOS

TiberOS boots with no GUI and prompts for user login, login using default credentials\
XSERVER is installed by default, make the below settings before starting X server

## Disable DPMS and screen blanking on the X Window System

-   DPMS Disable
    ```sh
    sudo vi /usr/share/X11/xorg.conf.d/10-extensions.conf
    ```
    Add
    ```
    Section "Extensions"
        Option "DPMS" "false"
    EndSection
    ```

-   Disable Screen Blanking and Timeouts
    ```sh
    sudo vi /usr/share/X11/xorg.conf.d/10-serverflags.conf
    ```
    Add
    ```
    Section "ServerFlags"
        Option "StandbyTime" "0"
        Option "SuspendTime" "0"
        Option "OffTime"     "0"
        Option "BlankTime"   "0"
    EndSection
    ```


## Start X Server
```
sudo X
```
- You can now see Black screen on monitors

### Check Monitor's resolution
Open SSH session to the Host system
```
DISPLAY=:0 xrandr
```
-   Output:
    ```
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

## Setup XDOTOOL in a container to scale applications to full-screen
> [!Note]
> This is needed if VM doesn't scale to full-screen after launching

### Install Docker in MF image:
```
sudo tdnf install -y moby-engine moby-cli ca-certificates
```
### Create a http-proxy.conf file with below contents to let docker work in Intel network with proxy.
```
sudo mkdir -p /etc/systemd/system/docker.service.d/

sudo vi /etc/systemd/system/docker.service.d/http-proxy.conf
```
Add
```
[Service]
Environment="HTTPS_PROXY=http://proxy-dmz.intel.com:912"
Environment="HTTP_PROXY=http://proxy-dmz.intel.com:911"
Environment="NO_PROXY=localhost,127.0.0.0/8,172.16.0.0/20,192.168.0.0/16,10.0.0.0/8,.intel.com,intel.com"
```
Restart Docker
```
sudo systemctl enable docker.service

sudo systemctl daemon-reload

sudo systemctl start docker.service
```
### Create alpine container and install XDOTOOL
```
sudo docker network create tiber-bridge-network

sudo docker run -it --rm --name tiber-resize-container \
    --privileged \
    -e DISPLAY=:0 \
    --ulimit nofile=65535:65535 \
    -v /tmp/.X11-unix:/tmp/.X11-unix \
    --network tiber-bridge-network \
    alpine:latest
```
#### From container set proxy, install xdotool and rescale VM:
```
export https_proxy=http://proxy-dmz.intel.com:912

apk add xdotool
```
-   To resize Virtual Machine Window
    ```
    xdotool search --onlyvisible --name "vm1" windowsize 1920 1080
    ```

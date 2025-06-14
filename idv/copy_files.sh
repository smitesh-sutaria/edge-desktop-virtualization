#!/bin/bash

if [ ! -d /opt/idv ]; then
    echo "idv directory not found. Creating directory"
    sudo mkdir /opt/idv
    sudo mkdir /opt/idv/launcher
    sudo mkdir /opt/idv/init
fi

username=$(getent passwd 1000 | cut -d: -f1)

# create/copy launcher files
sudo cp -r launcher /opt/idv/
sudo touch /opt/idv/launcher/start_all_vms.log
sudo chmod +x /opt/idv/launcher/start_all_vms.log
sudo chown $username:$username /opt/idv/launcher/start_all_vms.log

# create/copy init files
sudo cp -r init /opt/idv/
sudo touch /opt/idv/init/setup_sriov_vfs.log

# copy service files
sudo cp etc/systemd/user/idv-init.service /etc/systemd/user/idv-init.service
sudo cp etc/systemd/user/idv-launcher.service /etc/systemd/user/idv-launcher.service

# allow scripts to be run without password
echo "$username ALL=(ALL) NOPASSWD: /usr/bin/X,/opt/idv/init/setup_sriov_vfs.sh,/opt/idv/init/setup_display.sh,/opt/idv/init/setup_file_permissions.sh,/opt/idv/launcher/start_vm.sh,/opt/idv/launcher/start_all_vms.sh,/opt/idv/launcher/stop_vm.sh,/opt/idv/launcher/stop_all_vms.sh" | sudo tee -a /etc/sudoers.d/guest > /dev/null

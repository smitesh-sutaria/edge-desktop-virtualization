#!/bin/bash

sudo cp autologin.conf $HOME/rpmbuild/SOURCES
sudo cp etc/systemd/user/* $HOME/rpmbuild/SOURCES
sudo cp setup_permissions.sh $HOME/rpmbuild/SOURCES
sudo cp idv-solution-1.0.tar.gz $HOME/rpmbuild/SOURCES

sudo cp idv-solution.spec $HOME/rpmbuild/SPECS/idv-solution.spec 

#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

echo "Uninstalling MedHash Tools"
sudo pkgutil --forget com.ghifari160.medhash-tools

if [ -f /usr/local/bin/medhash-gen ]; then
    echo "Removing medhash-gen symlink"
    sudo rm -f /usr/local/bin/medhash-gen
fi

if [ -f /usr/local/bin/medhash-chk ]; then
    echo "Removing medhash-chk symlink"
    sudo rm -f /usr/local/bin/medhash-chk
fi

if [ -f /usr/local/bin/medhash-upgrade ]; then
    echo "Removing medhash-upgrade symlink"
    sudo rm -f /usr/local/bin/medhash-upgrade
fi

echo "Removing MedHash Tools"
installDir=$(dirname "$0")
sudo rm -rf $installDir

echo "Done"

exit 0

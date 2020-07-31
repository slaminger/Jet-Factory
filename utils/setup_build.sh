#!/bin/bash

#Ensure nothing happens outside the parent directory
cd "$(dirname "$0")"
cd ..

SCRIPT_DIRECTORY=$(pwd)

apt update -y && apt upgrade -y && apt install -y software-properties-common
add-apt-repository ppa:longsleep/golang-backports -y
apt update -y && apt install -y qemu qemu-user-static arch-install-scripts linux-image-generic golang-go libguestfs-tools libguestfs-dev
mkdir -p /root/linux
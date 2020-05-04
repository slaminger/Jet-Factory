# Jet Factory

AIO universal L4T distribution builder

## Scripts options

```
Usage: create-rootfs.sh [options] distribution-name
Options:
 -f, --force             Download setup files anyway
 --hekate                Build for Hekate
 -n, --no-docker         Build without Docker
 -s, --staging           Install built local packages
 --distro <name>         Select a distro to install
 -h, --help              Show this help text
```

## Building

On a Ubuntu host :

```sh
sudo apt-get install git tar wget p7zip unzip parted xz-utils dosfstools lvm2 qemu qemu-user-static proot
git clone https://github.com/Azkali/jet-factory
sudo ./jet-factory/create-rootfs.sh
```

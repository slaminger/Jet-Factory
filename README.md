# Jet Factory

AIO universal L4T distribution builder

## Scripts options

```
Usage: create-rootfs.sh [options] <distribution-name>
Options:
  -d, --docker             Build with Docker
  -f, --force              Download setup files anyway
  --hekate                 Build for Hekate
  -s, --staging            Install built local packages
  -h, --help               Show this help text
```

## Building

On a Ubuntu host :

```sh
sudo apt-get install git tar wget p7zip unzip parted xz-utils dosfstools lvm2 qemu qemu-user-static arch-install-scripts
```

```sh
git clone https://github.com/Azkali/jet-factory
```

```sh
sudo ./jet-factory/create-rootfs.sh
```

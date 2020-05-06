# Jet Factory

AIO universal L4T distribution builder

## Scripts options

```
Usage: create-rootfs.sh [options] <distribution-name>
Options:
  -d, --docker             Build with Docker
  --hekate                 Build for Hekate
  -k, --keep               Keep downloaded files
  -s, --staging            Install built local packages
  -h, --help               Show this help text
```

### Dependencies

**The following steps would consider your host as a Debian/Ubuntu based distribution, adapt if necessary**

*using Docker* :

```sh
sudo apt-get install lvm2 multipath-tools
```

*without using docker* :

```sh
sudo apt-get install git dtrx wget p7zip lvm2 qemu dosfstools qemu-user-static arch-install-scripts multipath-tools
```

### Build

```sh
git clone https://github.com/Azkali/jet-factory
```

```sh
sudo ./jet-factory/helpers/create-rootfs.sh
```

## Credits

@GavinDarkglider, @CTCaer, @ByLaws, @ave
For their amazing work, support and contribution to the scene

@Kitsumi, @parkerlreed
For their contribution to this project
# Jet Factory

Create live and flashable linux distribution root filesystem images.

## Scripts options

```txt
Usage: entrypoint.sh [options] <dir>
Options:
  -h, --hekate             Build an hekate installable filesystem
  -u, --usage              Show script usage
```

## Dependencies

**Install the dependencies only if you build without using Docker !**

Ubuntu :

```txt
sudo apt-get install qemu qemu-user-static arch-install-scripts linux-image-generic libguestfs-tools wget p7zip-full xz-utils
```

## Build example

- First, create a directory for the build :

```sh
mkdir -p ./linux
```

Then, choose one of the two methods for building :

- Option 1 Build without Docker :

```sh
sudo ./src/entrypoint.sh linux/
```

Or

- Option 2 - Build with Docker :

```sh
sudo docker run --privileged --rm -it -v "$PWD"/linux:/root/linux alizkan/jet-factory:latest
```

## Credits

### Special mentions

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden, @andrebraga \
For their awesome work, support and contribution to this project

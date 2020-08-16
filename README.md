# Jet Factory

Create live and flashable linux distribution root filesystem images.

## Dependencies

```txt
sudo apt-get install qemu qemu-user-static binfmt-support arch-install-scripts libguestfs-tools wget p7zip-full xz-utils
```

## Usage

```txt
Usage: entrypoint.sh <dir>
```

```txt
Variables:
    DISTRO=ARCH       Set target build distribution using file found in `configs/` directory
    HEKATE=true       Build hekate flashable image
```

## Build example

- First, create a directory for the build :

```sh
mkdir -p ./linux
```

Then, choose one of the two methods for building :

- Option 1 Build without Docker :

```sh
export DISTRO=ARCH
sudo ./src/entrypoint.sh linux/
```

Or

- Option 2 - Build with Docker :

```sh
sudo docker run --privileged --rm -e DISTRO=ARCH -it -v "$PWD"/linux:/root/linux alizkan/jet-factory:latest
```

## Credits

### Special mentions

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden, @andrebraga \
For their awesome work, support and contribution to this project

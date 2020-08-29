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
    DEVICE=tegra210   Device as set in the configs directory.
    DISTRO=arch       Target distribution using file found in DEVICE folder.
    HEKATE=true       Build hekate flashable image.
```

## Build example

- First, create a directory for the build :

```sh
mkdir -p ./linux
```

Then, choose one of the two methods for building :

- Option 1 Build without Docker :

```sh
export DISTRO=arch
export DEVICE=tegra210
sudo ./src/entrypoint.sh linux/
```

Or

- Option 2 - Build with Docker :

```sh
sudo docker run --privileged --rm -it -e DISTRO=arch -e DEVICE=tegra210 -v "$PWD"/linux:/out alizkan/jet-factory:latest
```

*You can override the workdir used in the docker, to use your own changes, without rebuilding the image by adding this repository directory as a volume to the docker command above.*

```sh
-v $(pwd):build/
```

## Credits

### Special mentions

@gavin_darkglider, @CTCaer, @ByLaws, @ave \
For their various work and contributions to switchroot.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden, @andrebraga1 \
For their work, support and direct contribution to this project.

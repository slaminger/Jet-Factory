# Jet Factory

Create live and flashable linux distribution root filesystem images.

## Scripts options

```
Usage: jetfactory [options]
Options:
  -arch               Platform build architecture; default aarch64
  -distro             Distribution to build: ubuntu, fedora, opensuse(leap, tumbleweed), slackware, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)
  -hekate             Build an hekate installable filesystem
  -force              Force to redownload files
  -help               Show this help text
```

## Build

To build arch :

```sh
mkdir -p ./linux
docker run --name jet --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --privileged --rm -it -e DISTRO=arch -v "$PWD"/linux:/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
```

## Credits

### Special mentions

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden \
For their awesome work, support and contribution to this project

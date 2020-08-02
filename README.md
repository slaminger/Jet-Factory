# Jet Factory

Create live and flashable linux distribution root filesystem images.

## Scripts options

```txt
Usage: jetfactory [options]
Options:
  --arch               Platform build architecture; default aarch64
  --distro             Distribution to build: ubuntu, fedora, opensuse(leap, tumbleweed), slackware, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)
  --hekate             Build an hekate installable filesystem
  --force              Force to redownload files
  --help               Show this help text
```

## Build example

To build Arch linux for hekate:

```sh
mkdir -p ./linux
docker run --privileged --rm -it -v "$PWD"/linux:/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:latest --distro arch --hekate
```

## Credits

### Special mentions

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden, @andrebraga \
For their awesome work, support and contribution to this project

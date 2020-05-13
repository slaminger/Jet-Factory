# Jet Factory

AIO universal L4T distribution builder

## Scripts options

```
Usage: jet-factory [options]
Options:
  -arch               Platform build architecture; default aarch64
  -distro             Distribution to build; ubuntu, fedora, opensuse(leap, tumbleweed), slackware, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)
  -hekate             Build an hekate installable filesystem
  -path               Output path; defaults to current directory
  -staging            Install built local packages
  -help               Show this help text
```

## Build

```sh
git clone https://github.com/Azkali/jet-factory
```

```sh
cd jet-factory
```

```sh
docker image build -t azkali/jet-factory:1.0.0 .
```

```sh
docker run --privileged --rm -it -e DISTRO=fedora -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD":/root azkali/jet-factory:1.0.0
```

## Credits

### Indirect contributors

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Direct contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden \
For their awesome work, support and contribution to this project

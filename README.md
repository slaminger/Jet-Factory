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

## Dependencies

**The following steps would consider your host as a Debian/Ubuntu based distribution, adapt if necessary**

```sh
sudo apt-get install lvm2 multipath-tools
```

## Build

```sh
git clone https://github.com/Azkali/jet-factory
```

```sh
cd jet-factory
```

```sh
go run factory.go
```

## Credits

### Indirect contributors

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Direct contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden \
For their awesome work, support and contribution to this project

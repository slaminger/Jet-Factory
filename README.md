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
docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO=fedora -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
```

## Credits

### Special mentions

@GavinDarkglider, @CTCaer, @ByLaws, @ave \
For their work and contributions.

### Contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden \
For their awesome work, support and contribution to this project

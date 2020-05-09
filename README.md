# Jet Factory

AIO universal L4T distribution builder

## Scripts options

```
Usage: jet-factory [options] <distribution-name>
Options:
  hekate             Build an hekate installable filesystem
  staging            Install built local packages
  help               Show this help text
```

## Dependencies

**The following steps would consider your host as a Debian/Ubuntu based distribution, adapt if necessary**

*using Docker* :

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

@GavinDarkglider, @CTCaer, @ByLaws, @ave
For their work and contributions.

### Direct contributors

@Stary2001, @Kitsumi, @parkerlreed, @AD2076, @PabloZaiden
For their awesome work, support and contribution to this project

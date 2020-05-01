package main

import (
	"log"

	"github.com/kless/gotool/flagplus"
)

type (
	// Common symbols
	pkg                string
	distributionName   string
	option             string
	desktopEnvironment map[distributionName]Desktops

	// Rootfs :
	Rootfs struct {
		distribution Distribution
		version      string
	}
	// Desktops :
	Desktops struct {
		distributionName
		list []string
	}
	// Distribution :
	Distribution struct {
		root map[distributionName]Packages
	}
	// Packages :
	Packages struct {
		desktopEnvironment
		pkgs []pkg
	}
)

// Consts :
const (
	arch     = "http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"
	fedora   = "http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version-number}/Server/aarch64/images/Fedora-${version}.raw.xz"
	opensuse = "http://download.opensuse.org/ports/aarch64/distribution/${version-number}/appliances/openSUSE-${version}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"
)

var (
	latest                 = ""
	avalaibleDistributions = []distributionName{
		"arch",
		"fedora",
		"opensuse",
		"ubuntu",
	}

	options = map[option]flagplus.Command{
		"coreboot": {
			"Compile latest coreboot",
			[]*flagplus.Subcommand{},
			[]string{""},
		},
		"linux": {
			"Compile latest Linux L4T",
			[]*flagplus.Subcommand{},
		},
	}
)

// Prepare :
func Prepare(cmd *flagplus.Subcommand, name distributionName, version string) *Rootfs {
	for _, ok := range avalaibleDistributions {
		if ok == name {
			switch desktop {
			case desktops[name][desktop]:
				desktop = desktops[name][desktop]
			default:
				log.Println("Unknown DE: %s, avalaible DE : %s", desktop, desktops[name])
				break
			}
			switch version {
			case "latest":
				version = latest
			default:
				log.Println("Using latest version number %s", latest)
				version = latest
			}
			return &Rootfs{distribution, version}
		}
	}
	panic(log.Println("Unknown distribution:", name))
}

// Coreboot : Compile latest coreboot
func Coreboot() flagplus.Subcommand {
	log.Println("Coreboot function is not implemented yet. Skipping...")
}

// Kernel : Compile latest Linux L4T
func Kernel() {
	log.Println("Kernel function is not implemented yet. Skipping...")
}

// Build :
func Build(Rootfs) interface{} {
	Prepare()
}

func main() {
	Build()
}

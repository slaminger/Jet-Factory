package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kless/gotool/flagplus"
)

type (
	// Common symbols
	distributionName            string
	desktopsList                map[distributionName][]string
	avalaibleDesktopEnvironment map[distributionName]desktopsList
	distribution                map[distributionName]Packages

	// Rootfs :
	Rootfs struct {
		distribution
		version string
	}
	// Packages :
	Packages struct {
		avalaibleDesktopEnvironment
		pkgs []string
	}
)

var (
	latest = ""

	arch = map[string][]string{
		"urls":   {"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
		"de":     {"xfce4", "lxde", "plasma"},
		"extras": {},
	}
	fedora = map[string][]string{
		"urls":   {"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version-number}/Server/aarch64/images/Fedora-${version}.raw.xz"},
		"de":     {"XFCE Desktop", "LXDE Desktop"},
		"extras": {},
	}
	opensuse = map[string][]string{
		"url":    {"http://download.opensuse.org/ports/aarch64/distribution/${version-number}/appliances/openSUSE-${version}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
		"de":     {"LXDE", "KDE", "XFCE"},
		"extras": {},
	}

	avalaibleDistributions = []distributionName{
		"arch",
		"fedora",
		"opensuse",
		"ubuntu",
	}

	avalaibleDesktops = avalaibleDesktopEnvironment{
		// TODO : for avalaibleDistributions items attach corresponding DE's
	}

	options = map[option]flagplus.Command{
		"coreboot": {
			"Compile latest coreboot",
			[]*flagplus.Subcommand{
				Run:       Coreboot,
				UsageLine: "",
			},
			[]string{""},
		},
		"linux": {
			"Compile latest Linux L4T",
			[]*flagplus.Subcommand{
				Run:       Kernel,
				UsageLine: "",
			},
		},
	}
)

func MatchDe(name distributionName, desktop string) (bool, string) {
	for _, desktops := range avalaibleDesktops[name] {
		for i := 0; i < len(desktops); i++ {
			if desktop != "" {
				if !strings.Contains(desktops[i], desktop) {
					log.Println("Unknown DE: %s, avalaible DE : %s", desktop, avalaibleDesktops[name])
					log.Println("Using default : XFCE")
				}
				desktop = desktops[i]
				log.Println("Found Desktop environment: %s", desktop)
			}
		}
	}
}

// FetchLatest :
func FetchLatest(name distributionName) {
	switch name {
	case "arch":
	case "fedora":
	}
}

// Prepare :
func Prepare(name distributionName, version, desktop string) *Rootfs {
	// Check if name martches a known distribution
	for _, ok := range avalaibleDistributions {
		if ok == name {
			if !MatchDe(name, desktop) {
				MatchDe(name, "XFCE")
			}
			switch version {
			case "latest":
				version = latest
			default:
				log.Println("Using latest version number %s", latest)
				version = latest
			}
			return &Rootfs{distribution{name, avalaibleDesktopEnvironment{name, desktop}}, version}
		}
	}
	log.Printf("Unknown distribution: %s", name)
	return nil
}

// Coreboot : Compile latest coreboot
func Coreboot() interface{} {
	log.Println("Coreboot function is not implemented yet. Skipping...")
	return nil
}

// Kernel : Compile latest Linux L4T
func Kernel() interface{} {
	log.Println("Kernel function is not implemented yet. Skipping...")
	return nil
}

// Build :
func Build(Rootfs) interface{} {
	Prepare()
	// TODO :
	// Create empty build with ${distributionName}-${version}-aarch64-${date} format directory as docker volume attached to /root/l4t/{
	if proc, err := SpawnProcess("mkdir", "-p", ""); err == nil {
	}
	// Wget shell script + Dockerfile from github to volume
	// Start docker
	// Create docker image
	// Start container build process
	// 7z direcotry with L4S-${distributionName}-${version}-aarch64-${date}.7z format
	return nil
}

// Wrappers
// Spawn Process
func SpawnProcess(cmd string, args ...string) (p *os.Process, err error) {
	if cmd, err = exec.LookPath(cmd); err == nil {
		fmt.Printf("> %s\n", args[0:])
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin,
			os.Stdout, os.Stderr}
		p, err := os.StartProcess(cmd, args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}

func main() {
	go Build()
	go Coreboot()
	go Kernel()

}

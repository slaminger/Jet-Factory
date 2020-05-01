package JetFactory

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type (
	// Base : Represent a distribution conatining a name, version, desktop environment and an optional list packages
	Base struct {
		name, version, desktop string
		pkgs, configs          []string
	}
)

var (
	// dir = ${distributionName}-${version}-aarch64-${date}
	dir, url string

	avalaibleDistributions = map[string]map[string][]string{
		"arch": {
			"urls":    {"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
			"de":      {"xfce4", "lxde", "plasma"},
			"pkgs":    {},
			"configs": {},
		},
		"fedora": {
			"urls":    {"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version}/Server/aarch64/images/Fedora-${version-release}.raw.xz"},
			"de":      {"XFCE Desktop", "LXDE Desktop"},
			"pkgs":    {},
			"configs": {},
		},
		"opensuse": {
			"url":     {"http://download.opensuse.org/ports/aarch64/distribution/${version}/appliances/openSUSE-${version-release}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
			"de":      {"LXDE", "KDE", "XFCE"},
			"pkgs":    {},
			"configs": {},
		},
		"ubuntu": {
			"url":     {""},
			"de":      {"LXDE", "KDE", "XFCE"},
			"pkgs":    {},
			"configs": {},
		},
	}
)

// SpawnProcess :
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

// MatchDe :
func MatchDe(name, desktop string) string {
	desktops := avalaibleDistributions[name]["de"]
	for i := 0; i < len(desktops); i++ {
		if desktop != "" {
			if !strings.Contains(desktops[i], desktop) {
				desktop = MatchDe(name, "XFCE")
				log.Println("Unknown DE: %s, avalaible DE : %s", desktop, desktops)
				log.Println("Using default : XFCE")
			}
			desktop = desktops[i]
			log.Println("Found Desktop environment: %s", desktop)
		}
	}
	return desktop
}

// BaseSetup :
func BaseSetup(name, version, desktop string) (r *Base, err error) {
	// Check if name match a known distribution
	for avalaible := range avalaibleDistributions {
		if !(name == avalaible) {
			log.Printf("Unknown distribution: %s", name)
			return nil, err
		}
		MatchDe(name, desktop)
		if !(version == "latest" || version == "" || name == "arch") {
			func() {
				// HTTP Query version find and match url
			}()
		}

		log.Println("Using latest version number: ")
		BaseSetup(name, "", desktop)
		r = &Base{name, version, desktop, nil, nil}
	}
	return r, nil
}

// BaseBuild :
func BaseBuild(name, version, desktop string) (p *os.Process, err error) {
	if root, err := BaseSetup(name, version, desktop); err != nil {
		return nil, err
	}
	// TODO :
	// Create empty build with ${distributionName}-${version}-aarch64-${date} format directory as docker volume attached to /root/l4t/
	_, mkdir := SpawnProcess("mkdir", "-p", dir)
	// Wget Dockerfile from github to volume
	_, get := SpawnProcess("wget", "https://raw.githubusercontent.com/Azkali/Jet-Factory/master/Dockerfile", "-P", dir)
	// Wget Dockerfile from github to volume
	_, getBase := SpawnProcess("wget", url, "-P", dir)
	// Start docker
	_, start := SpawnProcess("systemctl", "start", "docker.service", "docker.socket")
	// Create image
	_, img := SpawnProcess("docker", "image build -t", "opensusel4tbuild:1.0", dir)
	// Run container build process
	_, run := SpawnProcess("docker", "run --privileged --cap-add=SYS_ADMIN --rm -it", "-v "+dir+":/root/l4t/ l4tbuild:1.0", "/root/l4t/create-rootfs.sh")
	if !(mkdir == nil || get == nil || getBase == nil || start == nil || img == nil || run == nil) {
		return nil, err
	} else {
		// 7z directory with L4S-${distributionName}-${version}-aarch64-${date}.7z format
		proc, err := SpawnProcess("7z", "a", "L4S"+dir+".7z")
		return proc, err
	}
}

// JetFactory :
func JetFactory() {
	// BaseBuild()
}

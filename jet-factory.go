// TODO : Make a single function that handles parsing and replacing urls
// TODO 2 : Make function to check url avalaibility
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
			"urls":    {"http://download.opensuse.org/ports/aarch64/distribution/${version}/appliances/openSUSE-${version-release}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
			"de":      {"LXDE", "KDE", "XFCE"},
			"pkgs":    {},
			"configs": {},
		},
		"ubuntu": {
			"urls":    {""},
			"de":      {"LXDE", "KDE", "XFCE"},
			"pkgs":    {},
			"configs": {},
		},
	}
	variants = map[string][]string{
		"arch": {
			"blackarch",
			"arch-bang",
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
				log.Println("Using XFCE default")
			}
			desktop = desktops[i]
			log.Println("Found Desktop environment: %s", desktop)
		}
	}
	return desktop
}

// BaseSetup :
func Setup(name, version, desktop string) (r *Base, err error) {
	// Check if name match a known distribution
	for avalaible := range avalaibleDistributions {
		if !(name == avalaible) {
			log.Printf("Unknown distribution: %s", name)
			return nil, err
		}
		MatchDe(name, desktop)
		if !(version == "latest" || version == "" || name == "arch") {
			func() {
				// HTTP Query version find, match and construct url
			}()
		}
		if version == "" {
			if name == "arch" {
				log.Println("Using latest for arch anyway !")
				url = avalaibleDistributions[name]["url"][0]
			}
			func() {
				// HTTP Query latest find, match and construct url
			}()
		}

		log.Println("Using latest version number: ")
		Setup(name, "", desktop)
		r = &Base{name, version, desktop, nil, nil}
	}
	return r, nil
}

// BaseBuild :
func Build(name, version, desktop string) (p *os.Process, err error) {
	if root, err := Setup(name, version, desktop); err != nil {
		return nil, err
	}
	// TODO :
	// Create and got to dir
	// dir format - ${distributionName}-${version}-aarch64-${date}
	_, mkdir := SpawnProcess("mkdir", "-p", dir)
	// Wget Dockerfile from github to volume dir
	_, get := SpawnProcess("wget", "https://raw.githubusercontent.com/Azkali/Jet-Factory/master/Dockerfile", "-P", dir)
	// Replace URL variable in Dockerfile
	_, sed := SpawnProcess("sed", "-i \"s/${URL}\"/"+url+"/g Dockerfile")
	// Start docker
	_, start := SpawnProcess("systemctl", "start", "docker.service", "docker.socket")
	// Create image
	_, img := SpawnProcess("docker", "image build -t", "opensusel4tbuild:1.0", dir)
	// Run container build process attach buildir as volume to container
	_, run := SpawnProcess("docker", "run --privileged --cap-add=SYS_ADMIN --rm -it", "-v "+dir+":/root/l4t/ l4tbuild:1.0", "/root/l4t/create-rootfs.sh")
	if !(mkdir == nil || get == nil || start == nil || img == nil || run == nil || sed == nil) {
		return nil, err
	} else {
		// 7z directory with L4S-${distributionName}-${version}-aarch64-${date}.7z format
		proc, err := SpawnProcess("7z", "a", "L4S"+dir+".7z", dir+"/*")
		return proc, err
	}
}

// JetFactory :
func JetFactory() {
	// BaseBuild()
}

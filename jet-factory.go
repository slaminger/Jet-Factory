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
	// Base : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Base struct {
		basename, version, desktop string
		pkgs, configs              []string
	}
	// Variant : Represent a distribution variant
	Variant struct {
		Base
		variantname string
	}
)

var (
	// dir : ${distributionName}-${version}-aarch64-${date}
	basename, dir, url, variantname string
	isVariant                       = true
	root                            interface{}

	distributionsMap = map[string]map[string][]string{
		"arch": {
			"urls":    {"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
			"de":      {"xfce4", "lxde", "plasma"},
			"configs": {},
			"pkgs":    {},
		},
		"fedora": {
			"urls":    {"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version}/Server/aarch64/images/Fedora-${version-release}.raw.xz"},
			"de":      {"XFCE Desktop", "LXDE Desktop"},
			"configs": {},
			"pkgs":    {},
		},
		"opensuse": {
			"urls":    {"http://download.opensuse.org/ports/aarch64/distribution/${version}/appliances/openSUSE-${version-release}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
			"de":      {"LXDE", "KDE", "XFCE"},
			"configs": {},
			"pkgs":    {},
		},
		"ubuntu": {
			"urls":    {""},
			"de":      {"LXDE", "KDE", "XFCE"},
			"configs": {},
			"pkgs":    {},
		},
	}
	variantsMap = map[string][]string{
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

func IsDistro(name string) {
	// Check if name match a known distribution
	for avalaible := range distributionsMap {
		for _, variant := range variantsMap[avalaible] {
			if !(name == variant) {
				isVariant = false
				if !(name == avalaible) {
					log.Printf("Unknown distribution: %s", name)
					// Exit script
				}
			}
			variantname = variant
			// Assign to key corresponding to variant value
			// basename =
		}
		basename = name
	}
}

// GenDesktopEntry :
func GenDesktopEntry(name, desktop string) string {
	desktops := distributionsMap[name]["de"]
	for i := 0; i < len(desktops); i++ {
		if desktop != "" && !strings.Contains(desktops[i], desktop) {
			log.Println("Unknown DE: %s, avalaible DE : %s", desktop, desktops)
			log.Println("Using XFCE default")
			desktop = GenDesktopEntry(name, "XFCE")
		}
		log.Println("Found Desktop environment: %s", desktop)
		desktop = desktops[i]
	}
	return desktop
}

// GenVersionTag :
func GenVersionTag(name, version, desktop string) string {
	if !(version == "latest" || version == "" || name == "arch") {
		func() {
			// HTTP Query version find until match and construct url
			// If no match is found then use latest
			GenVersionTag(name, "", desktop)
		}()
	}

	if version == "" {
		if name == "arch" {
			log.Println("Using latest for arch anyway !")
			url = distributionsMap[name]["url"][0]
		}
		func() {
			// HTTP Query latest and construct url
		}()
		log.Println("Using latest version number :")
	}
	return version
}

// GenConfigs :
func GenConfigs(configs []string) []string {
	return configs
}

// GenPackagesList :
func GenPackagesList(pkgs []string) []string {
	return pkgs
}

// Build :
func Build(name, version, desktop string, configs, pkgs []string) (p *os.Process, err error) {
	IsDistro(name)
	desktop = GenDesktopEntry(name, desktop)
	version = GenVersionTag(name, version, desktop)
	root = &Base{name, version, desktop, GenConfigs(configs), GenPackagesList(pkgs)}
	if isVariant {
		root = &Variant{Base{name, version, desktop, GenConfigs(configs), GenPackagesList(pkgs)}, variantname}
	}
	// TODO :
	// Create dir - dir format : ${distributionName}-${version}-aarch64-${date}
	_, mkdir := SpawnProcess("mkdir", "-p", dir)
	// Wget Dockerfile from github to volume dir
	_, get := SpawnProcess("wget", "https://raw.githubusercontent.com/Azkali/Jet-Factory/master/Dockerfile", "-P", dir)
	// Replace variables in Dockerfile
	_, sed := SpawnProcess("sed", "-i \"s/${URL}\"/"+url+"/g", dir+"Dockerfile")
	// Start docker
	_, start := SpawnProcess("systemctl", "start", "docker.service", "docker.socket")
	// Create image
	_, img := SpawnProcess("docker", "image build -t", "opensusel4tbuild:1.0", dir)
	// Run container build process attach buildir as volume to container
	_, run := SpawnProcess("docker", "run --privileged --cap-add=SYS_ADMIN --rm -it", "-v", dir+":/root/l4t/", "l4tbuild:1.0", "/root/l4t/create-rootfs.sh")

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
	// Build()
}

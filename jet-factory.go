// TODO : Make a single function that handles parsing and replacing urls
// TODO 2 : Make function to check url avalaibility
package JetFactory

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	// Base : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Base struct {
		name, version, desktop string
		pkgs, configs          []string
	}
	// Variant : Represent a distribution variant
	Variant struct {
		Base
		variantName string
	}
)

var (
	// outputDir : ${distributionName}-${version}-aarch64-${date}
	baseName, variantName, outputDir, selectedMirror string
	isVariant                                        = true
	root                                             interface{}

	// TODO : Replace this by parsing config.yaml file
	distributionsMap = map[string]map[string][]string{
		"arch": {
			"urls":    {"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
			"de":      {"xfce4", "lxde", "plasma"},
			"configs": {},
			"pkgs":    {},
		},
		"fedora": {
			"urls":    {"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version}/Server/aarch64/images/Fedora-${version-release}.raw.xz"},
			"de":      {"Xfce Desktop", "LXDE Desktop"},
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

func IsDistro(name string) (find string, ok bool) {
	// Check if name match a known distribution
	for base := range distributionsMap {
		for _, variant := range variantsMap[base] {
			if isVariant != (name == variant) {
				isVariant = false
				if !(name == base) {
					return "", false
				}
				name = base
			}
			name, baseName = variant, base
		}
	}
	return name, true
}

// GenDesktopEntry :
func GenDesktopEntry(desktop string) string {
	desktops := distributionsMap[baseName]["de"]
	for i := 0; i < len(desktops); i++ {
		if desktop != "" && !strings.Contains(desktops[i], desktop) {
			log.Println("Unknown DE: %s, avalaible DE : %s\nUsing XFCE default", desktop, desktops)
			desktop = GenDesktopEntry("XFCE")
		}
		log.Println("Found Desktop environment: %s", desktop)
		desktop = desktops[i]
	}
	return desktop
}

func GenUrl(version string) {
	// HTTP Query version find until match and construct url
	for i := 0; i < len(distributionsMap[baseName]["urls"]); i++ {
		for _, avalaibleMirror := range distributionsMap[baseName]["urls"] {
			// Replace name and version in url for avalaibleMirror
			u, err := url.ParseRequestURI(avalaibleMirror)
			if err != nil {
				log.Println("Mirror URL : %s not available... Skipping to the next one", avalaibleMirror)
			}
			log.Println("Mirror URL: %s", avalaibleMirror)
			selectedMirror = u.String()
		}
	}
	log.Panicln("Couldn't found any valid urls... Exiting")
}

// GenVersionTag :
func GenVersionTag(version, desktop string) string {
	if version == "" || baseName == "arch" {
		GenVersionTag("latest", desktop)
	} else if version == "latest" {
		log.Println("Using latest version avalaible !")
		GenVersionTag(version, desktop)
	} else {
		// TODO :
		// Try to query version number
		// if !(version) {
		// If fails default to latest
		// }
	}
	GenUrl(version)
	return version
}

// TODO
// GenConfigs :
func GenConfigs(configs []string) []string {
	return configs
}

// TODO
// GenPackagesList :
func GenPackagesList(pkgs []string) []string {
	return pkgs
}

// JetFactory :
func JetFactory(name, version, desktop string, configs, pkgs []string) (p *os.Process, err error) {
	var dirName string

	if dirName, ok := IsDistro(name); !ok {
		log.Panicln("No distribution found for : %s", dirName)
	}

	desktop = GenDesktopEntry(desktop)
	version = GenVersionTag(version, desktop)
	configs = GenConfigs(configs)
	pkgs = GenPackagesList(pkgs)

	outputDir = dirName + "-" + "-" + version + "-" + "aarch64" + "-" + time.Now().String()

	root = &Base{baseName, version, desktop, configs, pkgs}

	if isVariant {
		dirName = variantName
		root = &Variant{Base{baseName, version, desktop, configs, pkgs}, variantName}
	}
	// Create dir - dir format : ${distributionName}-${version}-aarch64-${date} && // Wget Dockerfile from github to volume dir && // Replace variables in Dockerfile
	_, mkdir := SpawnProcess("mkdir", "-p", outputDir)
	_, get := SpawnProcess("wget", "https://raw.githubusercontent.com/Azkali/Jet-Factory/master/Dockerfile", "-P", outputDir)
	_, sed := SpawnProcess("sed", "-i", "'s/URL/"+selectedMirror+"/g;'", "'s/NAME/"+name+"/g;'", outputDir+"create-rootfs.sh")

	// Start docker && // Create image && // Run container build process attach buildir as volume to container
	_, start := SpawnProcess("systemctl", "start", "docker.service", "docker.socket")
	_, img := SpawnProcess("docker", "image build -t", "l4tbuild:1.0", outputDir)

	if !(mkdir == nil && get == nil && sed == nil && start == nil && img == nil) {
		panic(err)
	} else {
		// 7z directory with L4S-${distributionName}-${version}-aarch64-${date}.7z format
		proc, err := SpawnProcess("docker", "run --privileged --cap-add=SYS_ADMIN --rm -it", "-v", outputDir+":/root/builder/", "l4tbuild:1.0", "/root/builder/create-rootfs.sh")
		return proc, err
	}
}

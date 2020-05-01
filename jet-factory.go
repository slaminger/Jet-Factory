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
		pkgs, configs          interface{}
	}
)

var (
	latest string

	avalaibleDistributions = map[string]map[string][]string{
		"arch": map[string][]string{
			"urls":    []string{"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
			"de":      []string{"xfce4", "lxde", "plasma"},
			"extras":  []string{},
			"configs": []string{},
		},
		"fedora": map[string][]string{
			"urls":    []string{"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version-number}/Server/aarch64/images/Fedora-${version}.raw.xz"},
			"de":      []string{"XFCE Desktop", "LXDE Desktop"},
			"extras":  []string{},
			"configs": []string{},
		},
		"opensuse": map[string][]string{
			"url":     []string{"http://download.opensuse.org/ports/aarch64/distribution/${version-number}/appliances/openSUSE-${version}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
			"de":      []string{"LXDE", "KDE", "XFCE"},
			"extras":  []string{},
			"configs": []string{},
		},
		"ubuntu": {
			"url":     []string{""},
			"de":      []string{"LXDE", "KDE", "XFCE"},
			"extras":  []string{},
			"configs": []string{},
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

// Prepare :
func PrepareBase(name, version, desktop string) (r *Base, err error) {
	// Check if name match a known distribution
	for avalaible := range avalaibleDistributions {
		if name == avalaible {
			MatchDe(name, desktop)
			switch version {
			case "latest":
				version = latest
			default:
				log.Println("Using latest version number %s", latest)
				version = latest
			}
			r = &Base{name, version, desktop, nil, nil}
			return r, nil
		}
	}
	log.Printf("Unknown distribution: %s", name)
	return nil, err
}

// Build :
func BuildBase(name, version, desktop string) (p *os.Process, err error) {
	// dir = ${distributionName}-${version}-aarch64-${date}
	var dir string
	var url string

	if root, err := PrepareBase(name, version, desktop); err == nil {
		dir = ""
		url = ""
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
		// 7z direcotry with L4S-${distributionName}-${version}-aarch64-${date}.7z format
		proc, err := SpawnProcess("7z", "a", "L4S"+dir+".7z")
		return proc, err
	}
}

func JetFactory() {
	// BuildBase()
}

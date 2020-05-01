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
	distributionName string

	// Base :
	Base struct {
		name             distributionName
		version, desktop string
		pkgs             interface{}
	}
)

var (
	latest string

	avalaibleDistributions = map[distributionName]map[string][]string{
		"arch": map[string][]string{
			"urls":   []string{"http://os.archlinuxarm.org/os/ArchLinuxARM-aarch64-latest.tar.gz"},
			"de":     []string{"xfce4", "lxde", "plasma"},
			"extras": []string{},
		},
		"fedora": map[string][]string{
			"urls":   []string{"http://mirrors.ircam.fr/pub/fedora/linux/version-releases/${version-number}/Server/aarch64/images/Fedora-${version}.raw.xz"},
			"de":     []string{"XFCE Desktop", "LXDE Desktop"},
			"extras": []string{},
		},
		"opensuse": map[string][]string{
			"url":    []string{"http://download.opensuse.org/ports/aarch64/distribution/${version-number}/appliances/openSUSE-${version}-ARM-${desktop}.aarch64-rootfs.aarch64.tar.xz"},
			"de":     []string{"LXDE", "KDE", "XFCE"},
			"extras": []string{},
		},
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

// FetchLatest :
func FetchLatest(name distributionName) {
	switch name {
	case "arch":
	case "fedora":
	}
}

func MatchDe(name distributionName, desktop string) string {
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
func Prepare(name distributionName, version, desktop string) (r *Base, err error) {
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
			r = &Base{name, version, desktop, nil}
			return r, nil
		}
	}
	log.Printf("Unknown distribution: %s", name)
	return nil, err
}

// Build :
func Build(name distributionName, version, desktop string) (p *os.Process, err error) {
	// dir = ${distributionName}-${version}-aarch64-${date}
	var dir string
	var url string

	if root, err := Prepare(name, version, desktop); err == nil {
		dir = ""
		url = ""
	}
	// TODO :
	// Create empty build with ${distributionName}-${version}-aarch64-${date} format directory as docker volume attached to /root/l4t/{
	if proc, err := SpawnProcess("mkdir", "-p", dir); err == nil {
		log.Println(proc)
	}
	// Wget shell script + Dockerfile from github to volume
	if proc, err := SpawnProcess("wget", url, "-P", dir); err == nil {
		log.Println(proc)
	}
	// Start docker
	if proc, err := SpawnProcess("systemctl", "start", "docker", "docker.socket"); err == nil {
		log.Println(proc)
	}
	// Create docker image
	if proc, err := SpawnProcess("docker", "image build -t", "opensusel4tbuild:1.0", dir); err == nil {
		log.Println(proc)
	}
	// Start container build process
	if proc, err := SpawnProcess("docker", "run --privileged --cap-add=SYS_ADMIN --rm -it -v", dir+":/root/l4t/ l4tbuild:1.0", "/root/l4t/create-rootfs.sh"); err == nil {
		log.Println(proc)
	}
	// 7z direcotry with L4S-${distributionName}-${version}-aarch64-${date}.7z format
	if proc, err := SpawnProcess("7z", "a", "L4S"+dir+".7z"); err == nil {
		log.Println(proc)
		return proc, nil
	}
	return nil, err
}

func main() {
	go Build()
}

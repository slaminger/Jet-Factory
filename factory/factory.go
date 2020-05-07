package main

// TODO : Make a single function that handles parsing and replacing urls

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
	isVariant, isAndroid                             = true, false
	root                                             interface{}
	dockerImageName                                  = "docker.io/library/ubuntu:18.04"

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
		"lineage": {
			"icosa",
			"foster",
			"foster_tab",
		},
	}
)

// SpawnProcess : Spawns a shell subprocess
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

// SpawnContainer : Spawns a container based on ubunutu:18.04
func SpawnContainer(cmd []string) error {
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: dockerImageName,
		Cmd:   cmd,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, baseName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if _, err := cli.ContainerWait(ctx, resp.ID); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

// IsDistro : Chechks if a distribution is avalaible in the config files
func IsDistro(name string) *string {
	// Check if name match a known distribution
	for base := range distributionsMap {
		for _, variant := range variantsMap[base] {
			if name == variant {
				isVariant = true
				variantName, baseName = variant, base
				IsAndroid(baseName)
				return &name
			}
		}
		if name == base {
			IsAndroid(name)
			return &name
		}
	}
	return nil
}

// IsAndroid : Checks if the produced build will be Android
func IsAndroid(name string) {
	if name == "lineage" {
		isAndroid = true
	}
}

// GenDesktopEntry : Retrieves a Desktop Environment from the list of avalaibles for the distribution; returns the one find or the first avalaible
func GenDesktopEntry(desktop string) *string {
	desktops := distributionsMap[baseName]["de"]
	for i := 0; i < len(desktops); i++ {
		if strings.Contains(strings.ToLower(desktops[i]), strings.ToLower(desktop)) {
			log.Println("Found Desktop environment: ", desktop)
			desktop = desktops[i]
			return &desktop
		}
		log.Println("Unknown DE:", desktop, "avalaible DE: ", desktops, "\nUsing: ", desktops[0], "as default")
	}
	return &desktops[0]
}

// QueryMirrorAvalaibility : Get the url demanded for the distribution version if specified
func QueryMirrorAvalaibility(version string) *string {
	// HTTP Query version find until match and construct url
	for i := 0; i < len(distributionsMap[baseName]["urls"]); i++ {
		for _, avalaibleMirror := range distributionsMap[baseName]["urls"] {
			// Replace name and version in url for avalaibleMirror
			u, err := url.ParseRequestURI(avalaibleMirror)
			if err != nil {
				log.Println("Mirror: ", avalaibleMirror, "not available... Skipping to the next one")
			}
			log.Println("Mirror URL selected : ", avalaibleMirror)
			selectedMirror = u.String()
			return &selectedMirror
		}
	}
	log.Panicln("Couldn't found any valid urls... Exiting")
	return nil
}

// GenerateVersionTag : Retrieve a URL for a distribution based on the version
func GenerateVersionTag(version, desktop string) *string {
	if version != "latest" && len(version) > 0 {
		if len(variantName) > 0 {
			// base, name := baseName, variantName
		} else {
			// name := baseName
		}
	} else {

		log.Println("Using latest version avalaible !")
	}
	return &version
}

// RunInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func RunInChrootEnv(configs []string) error {
	for _, config := range configs {
		if err := SpawnContainer([]string{"arch-chroot", config}); err != nil {
			return err
		}
	}
	return nil
}

// GeneratePackagesList : Installs packages list; Returns nil if successful
func GeneratePackagesList(pkgs []string) *error {
	// TODO : Assign to distribution package manager
	var pkgManager string

	if _, err := SpawnProcess(pkgManager, pkgs...); err != nil {
		return &err
	}
	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro, version, desktop string, configs, pkgs []string, basepath string) (p *os.Process, err error) {
	var dirName string
	var basePath = "."

	if len(basepath) > 0 {
		basePath = basepath
	}

	if dirName := IsDistro(distro); dirName == nil {
		log.Println("No distribution found : ", dirName)
		return nil, err
	}

	log.Println("Building: ", distro)

	desktop = *GenDesktopEntry(desktop)
	version = *GenerateVersionTag(version, desktop)
	mirror := *QueryMirrorAvalaibility(version)
	// configs = GenConfigs(configs)

	if isVariant {
		dirName = variantName
		root = &Variant{Base{baseName, version, desktop, configs, pkgs}, variantName}
	} else {
		dirName = baseName
		root = &Base{baseName, version, desktop, configs, pkgs}
	}

	fmt.Println("dirName:", dirName)

	if !isAndroid {
		// Create dir - dir format : ${distributionName}-${version}-aarch64-${date} && // Wget Dockerfile from github to volume dir && // Replace variables in Dockerfile
		outputDir = basePath + "/" + baseName + "-" + "-" + version + "-" + "aarch64" + "-" + time.Now().String()

		_, mkdir := SpawnProcess("mkdir", "-p", outputDir)
		_, getroot := SpawnProcess("wget", "-nc -q --show-progress", mirror, "-P", outputDir)
		_, img := SpawnProcess("docker", "image build -t", "l4tbuild:1.0", outputDir)
		// pkgs = GenPackagesList(pkgs)

		if !(mkdir == nil && img == nil && getroot == nil) {
			return nil, err
		}

		// 7z directory with L4S-${distributionName}-${version}-aarch64-${date}.7z format
		// TODO : Proper docker api call to run a container
		return nil, err
	}
	// Create dir - dir format : ${distributionName}-${version}-aarch64-${date} && // Wget Dockerfile from github to volume dir && // Replace variables in Dockerfile
	outputDir = basePath + "/" + baseName
	_, mkdir := SpawnProcess("mkdir", "-p", outputDir)

	if !(mkdir == nil) {
		return nil, err
	}

	proc, err := SpawnProcess("docker", "run --rm -ti -e ROM_NAME="+dirName, "--volume", basePath+":/root/android pablozaiden/switchroot-android-build")
	return proc, err
}

func main() {
	var distro string
	flag.StringVar(&distro, "distro", "arch", "the distro you want to build: ubuntu, fedora, gentoo, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)")

	// is string ok?
	var version string
	flag.StringVar(&version, "version", "", "Version to build")

	var desktop string
	flag.StringVar(&desktop, "desktop", "", "DE environment")

	var basepath string
	flag.StringVar(&basepath, "basepath", ".", "Path to use as Docker storage, can be a mounted external device")

	//TODO get arrays args

	flag.Parse()

	// set default for android build
	if distro == "lineage" {
		distro = "icosa"
	}

	// TODO : Run a bare ubuntu:18.04 container that downloads and extracts all bits then stop the container and check for integrity
	// TODO-2 : Check if it's an image file, check if there is an LVM2 partition
	// TODO-2.a: Either way create the destination ext4 disk image
	// TODO-2.b: If there's no LVM2 partition, run another container and extract the ext2,3,4 to the container
	// TODO-2.c: If there is an LVM2 partition, mount the paritition on the host, run another container and copy the content of the part to the container
	// TODO-2.d: If copy was successful then Commit container to a new Image
	// TODO-3: Clone the parition and attach this Clone to a container, then chroot
	// TODO-3.b: If chroot went well then replace Original by Clone, Commit container to new Image
	// TODO-4: Run another container and create the final partition.

	// Factory(distro, version, desktop, []string{}, []string{}, basepath)
	if err := SpawnContainer([]string{"cat", "/etc/os-release"}); err != nil {
		log.Println(err)
	}
}

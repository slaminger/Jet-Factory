package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type (
	// Base : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Base map[string]map[string][]string
	// Variant : Represent a distribution variant
	Variant struct {
		Base
		Extras map[string]map[string][]string
	}
)

var (
	// outputDir : ${distributionName}-${version}-aarch64-${date}
	baseName, variantName        string
	selectedMirror, dockerOutput string
	hekate, staging              = "--hekate", "--staging"
	isVariant, isAndroid         = false, false
	dockerImageName              = "docker.io/library/ubuntu:18.04"

	baseJSON, _    = ioutil.ReadFile("../configs/base.json")
	variantJSON, _ = ioutil.ReadFile("../configs/variants.json")

	baseDistros    = Base{}
	variantDistros = Variant{Base: baseDistros}

	_ = json.Unmarshal([]byte(baseJSON), &baseDistros)
	_ = json.Unmarshal([]byte(variantJSON), &variantDistros)
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

// SpawnContainer : Spawns a container based on dockerImageName
func SpawnContainer(cmd []string, volumes map[string]struct{}) error {
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
		Image:   dockerImageName,
		Cmd:     cmd,
		Volumes: volumes,
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

	//  dockerOutput =
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

// IsDistro : Chechks if a distribution is avalaible in the config files
func IsDistro(name string) (err error) {
	// Check if name match a known distribution
	for base := range baseDistros {
		baseName = base
		if name == base {
			return nil
		}
		for _, variant := range baseDistros[base]["variants"] {
			if name == variant {
				isVariant = true
				variantName = variant
				return nil
			}
		}
	}
	return err
}

// GenerateURLfromVersionTag : Retrieve a URL for a distribution based on a version
func GenerateURLfromVersionTag(name string) *string {
	// TODO : HTTP Query version find until match and construct url
	for _, avalaibleMirror := range baseDistros[name]["mirror_urls"] {
		constructedURL := ""
		if _, err := url.ParseRequestURI(avalaibleMirror); err != nil {
			log.Println("Couldn't found mirror:", avalaibleMirror)
		}
		log.Println("Mirror URL selected : ", avalaibleMirror)
		return &constructedURL
	}
	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(configs []string) error {
	if err := SpawnContainer([]string{"/bin/bash", "preChroot.sh"}, nil); err != nil {
		return err
	}

	for _, config := range configs {
		if err := SpawnContainer([]string{"arch-chroot", config}, nil); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"/bin/bash", "postChroot.sh"}, nil); err != nil {
		return err
	}

	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(pkgs []string) error {
	var pkgManager string

	if err := SpawnContainer([]string{"/bin/bash", "preChroot.sh"}, nil); err != nil {
		return err
	}

	// TODO : Retrieve the returned shell variable
	if err := SpawnContainer([]string{"arch-chroot", "/bin/bash", "findPackageManager.sh"}, nil); err != nil {
		// pkgManager =
		return err
	}

	// TODO : rule for staging packages

	if err := SpawnContainer([]string{"arch-chroot", pkgManager, strings.Join(pkgs, ",")}, nil); err != nil {
		return err
	}

	if err := SpawnContainer([]string{"/bin/bash", "postChroot.sh"}, nil); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, outDir string) error {
	var pkgs, configs []string

	basePath := outDir
	if !(len(outDir) > 0) {
		basePath = "."
	}

	if err := IsDistro(distro); err != nil {
		log.Println("No distribution found for: ", err)
		return err
	}

	if isVariant {
		image := variantDistros.Extras[distro]
		pkgs = image["packages"]
		configs = image["configs"]
	} else {
		image := baseDistros[distro]
		pkgs = image["packages"]
		configs = image["configs"]
	}

	mirror := *GenerateURLfromVersionTag(distro)

	basePath = basePath + "/" + baseName
	log.Println("Building: ", distro, "in dir: ", basePath)
	if err := os.MkdirAll(basePath, 755); err != nil {
		return err
	}

	if !isAndroid {
		// TODO : Retrieve the returned shell variable
		if _, err := SpawnProcess("/bin/bash", "prepare.sh", mirror, basePath); err != nil {
			return err
		}

		if err := InstallPackagesInChrootEnv(pkgs); err != nil {
			return err
		}

		if err := ApplyConfigsInChrootEnv(configs); err != nil {
			return err
		}

		if err := SpawnContainer([]string{"/bin/bash", "createImage.sh", hekate, baseName, basePath}, nil); err != nil {
			return err
		}
	} else {
		dockerImageName = "pablozaiden/switchroot-android-build:v1"
		if err := SpawnContainer([]string{"ROM_NAME=" + distro, "--volume", basePath + ":/root/android"}, nil); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var distro, basepath string

	flag.StringVar(&distro, "distro", "arch", "the distro you want to build: ubuntu, fedora, gentoo, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)")
	flag.StringVar(&basepath, "basepath", ".", "Path to use as Docker storage, can be a mounted external device")
	if ok := flag.Bool("hekate", false, "Build an hekate installable filesystem"); !*ok {
		hekate = ""
	}
	if ok := flag.Bool("staging", false, "Install built local packages"); !*ok {
		staging = ""
	}
	flag.Parse()

	// Sets default for android build
	if distro == "lineage" {
		distro = "icosa"
		isAndroid = true
	}

	Factory(distro, basepath)
}

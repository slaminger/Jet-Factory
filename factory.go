package main

// TODO-3 (x2 - line181 line224) : Retrieve the returned shell variable
// TODO-4 : Handle Staging packages installation

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
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/manifoldco/promptui"
)

type (
	// Base : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Base struct {
		Params map[string]map[string][]string
		Urls   map[string]map[string][]string `json:"buildarch"`
	}

	// Variant : Represent a distribution variant
	Variant struct {
		Base
		Extras map[string]map[string][]string `json:"variations"`
	}
)

var (
	isVariant, isAndroid    = false, false
	baseName, variantName   string
	dockerOutput, buildarch string
	hekate, staging         = "--hekate", "--staging"
	dockerImageName         = "docker.io/library/ubuntu:18.04"

	baseJSON, _ = ioutil.ReadFile("../configs/base.json")

	baseDistros    = Base{}
	variantDistros = Variant{Base: baseDistros}
	_              = json.Unmarshal([]byte(baseJSON), &baseDistros)
	_              = json.Unmarshal([]byte(baseJSON), &variantDistros)

	architectures = [...]string{"aarch64", "amd64", "i386", "arm"}
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
func SpawnContainer(cmd []string, volume *[2]string) error {
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
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volume[0],
				Target: volume[1],
			},
		},
	}, nil, baseName)
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
	for base := range baseDistros.Params {
		baseName = base
		if name == base {
			return nil
		}
		for _, variant := range baseDistros.Params[base]["variants"] {
			if name == variant {
				isVariant = true
				variantName = variant
				return nil
			}
		}
	}
	return err
}

// CliSelector : Select an item in a menu froim cli
func CliSelector(label string, items []string) *string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	return &result
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(name string, path [2]string) *error {
	var constructedURL string
	var versions []string

	for _, avalaibleMirror := range baseDistros.Urls[buildarch]["mirror_urls"] {
		if !(name == "arch" || name == "slackware" || name == "tumbleweed") && !(name == "ubuntu" && buildarch == "aarch64") {
			if name == "fedora" {
				// versions := []string{}
			} else if name == "leap" {
				// versions := []string{}
			}
			if ok := CliSelector(name, versions); ok == nil {
				return nil
			}
			log.Println("Mirror URL selected : ", constructedURL)
		} else {
			constructedURL := avalaibleMirror
		}
		if _, err := url.ParseRequestURI(constructedURL); err != nil {
			log.Println("Couldn't found mirror:", constructedURL)
		}
	}

	// TODO-3
	if _, err := SpawnProcess("/bin/bash", "prepare.sh", constructedURL, path[0]); err != nil {
		return &err
	}

	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path [2]string) error {
	if err := SpawnContainer([]string{"/bin/bash", "preChroot.sh"}, &path); err != nil {
		return err
	}

	if isVariant {
		for _, config := range baseDistros.Params[baseName]["configs"] {
			if err := SpawnContainer([]string{"arch-chroot", config}, &path); err != nil {
				return err
			}
		}
	}

	for _, config := range variantDistros.Extras[variantName]["configs"] {
		if err := SpawnContainer([]string{"arch-chroot", config}, &path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"/bin/bash", "postChroot.sh"}, &path); err != nil {
		return err
	}

	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(path [2]string) error {
	var pkgManager string

	if err := SpawnContainer([]string{"/bin/bash", "preChroot.sh"}, &path); err != nil {
		return err
	}

	// TODO-3
	if err := SpawnContainer([]string{"arch-chroot", "/bin/bash", "findPackageManager.sh"}, &path); err != nil {
		// pkgManager =
		return err
	}

	// TODO-4
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", pkgManager, strings.Join(variantDistros.Extras[variantName]["packages"], ",")}, &path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", pkgManager, strings.Join(baseDistros.Params[baseName]["packages"], ",")}, &path); err != nil {
		return err
	}

	if err := SpawnContainer([]string{"/bin/bash", "postChroot.sh"}, &path); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, outDir string) error {
	if !(len(outDir) > 0) {
		outDir = "."
	}

	if err := IsDistro(distro); err != nil {
		flag.Usage()
		return err
	}

	basePath := outDir + "/" + distro
	path := [2]string{basePath, "/root/" + distro}

	log.Println("Building: ", distro, "in dir: ", basePath)
	if err := os.MkdirAll(basePath, 755); err != nil {
		return err
	}

	if !isAndroid {
		if err := *DownloadURLfromTags(distro, path); err != nil {
			return err
		}

		if err := InstallPackagesInChrootEnv(path); err != nil {
			return err
		}

		if err := ApplyConfigsInChrootEnv(path); err != nil {
			return err
		}

		if err := SpawnContainer([]string{"/bin/bash", "createImage.sh", hekate, baseName, "."}, &path); err != nil {
			return err
		}
	} else {
		dockerImageName = "pablozaiden/switchroot-android-build:1.0.0"
		if err := SpawnContainer([]string{"-e ROM_NAME=" + distro}, &path); err != nil {
			return err
		}
	}
	log.Println("Done!")
	return nil
}

func main() {
	var distro, basepath string
	flag.StringVar(&distro, "distro", "", "the distro you want to build: ubuntu, fedora, gentoo, arch(blackarch, arch-bang), lineage(icosa, foster, foster_tab)")
	flag.StringVar(&basepath, "basepath", ".", "Path to use as Docker storage, can be a mounted external device")
	flag.StringVar(&buildarch, "arch", "aarch64", "Set the platform build architecture.")
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

	if distro == "opensuse" {
		distro = "leap"
	}

	for _, arch := range architectures {
		if buildarch == arch {
			buildarch = arch
		}
	}

	if distro == "" {
		flag.Usage()
	} else {
		Factory(distro, basepath)
	}
}

package main

// TODO-4 : Handle Staging packages installation

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/manifoldco/promptui"
)

type (
	// Distribution : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Distribution struct {
		Name          string              `json:"name"`
		Configs       []string            `json:"configs"`
		Packages      []string            `json:"packages"`
		Architectures map[string][]string `json:"buildarch"`
		Variants      []Variant           `json:"variants"`
	}

	// Variant : Represent a distribution variant
	Variant struct {
		Name     string   `json:"name"`
		Configs  []string `json:"configs"`
		Packages []string `json:"packages"`
	}
)

var (
	distribution        Distribution
	variant             Variant
	baseName, buildarch string

	isVariant, isAndroid = false, false
	hekate, staging      = "--hekate", "--staging"
	dockerImageName      = "docker.io/library/ubuntu:18.04"

	baseJSON, _ = ioutil.ReadFile("./setup/base.json")
	basesDistro = []Distribution{}
	_           = json.Unmarshal([]byte(baseJSON), &basesDistro)
)

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

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

// IsDistro : Checks if a distribution is avalaible in the config files
func IsDistro(name string) (err error) {
	// Check if name match a known distribution
	for i := 0; i < len(basesDistro); i++ {
		if name == basesDistro[i].Name {
			baseName = basesDistro[i].Name
			distribution = Distribution{Name: basesDistro[i].Name, Architectures: basesDistro[i].Architectures, Configs: basesDistro[i].Configs, Packages: basesDistro[i].Packages}
			return nil
		}
		for j := 0; j < len(basesDistro[i].Variants); j++ {
			if name == basesDistro[i].Variants[j].Name {
				isVariant = true
				variant = Variant{Name: basesDistro[i].Variants[j].Name}
				return nil
			}
		}
	}
	return err
}

// IsValidArchitecture : Check if the inputed architecture can be found for the distribution
func IsValidArchitecture() (archi *string) {
	for archis := range distribution.Architectures {
		log.Println(archis)
		if buildarch == archis {
			log.Println("Found valid architecture: ", buildarch)
			return &buildarch
		}
	}
	log.Println(buildarch, "is not a valid architecture !")
	return nil
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

// WalkURL :
func WalkURL(source, regex string) []string {
	resp, err := http.Get(source)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK) {
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	log.Println(strip.StripTags(string(bodyBytes)))
	search, _ := regexp.Compile(regex)
	found := search.FindStringSubmatch(strip.StripTags(string(bodyBytes)))

	return found
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(path [2]string) (err error) {
	var constructedURL string
	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		if strings.Contains(avalaibleMirror, "{VERSION}") || strings.Contains(avalaibleMirror, "{BUILDARCH}") {

			avalaibleMirror = strings.Replace(avalaibleMirror, "{BUILDARCH}", buildarch, 1)

			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			regexURL := WalkURL(constructedURL, "[[:digit:]](.*?)/")

			version := CliSelector("Select a version: ", regexURL)
			if version == nil {
				return err
			}

			avalaibleMirror = strings.Replace(avalaibleMirror, "{VERSION}", *version, 1)

			regexURL = WalkURL(avalaibleMirror, "[[:digit:]](.*?)/")
			imageFile := CliSelector("Select an image file: ", regexURL)

			if imageFile == nil {
				return err
			}

		} else {
			constructedURL = avalaibleMirror
		}

		if _, err := url.ParseRequestURI(constructedURL); err != nil {
			log.Println("Couldn't found mirror:", constructedURL)
			return err
		}
	}
	log.Println("Mirror URL selected : ", constructedURL)

	err = exec.Command("/bin/bash", "prepare.sh", constructedURL, path[0]).Run()
	if err != nil {
		return err
	}

	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path [2]string) error {
	if err := SpawnContainer([]string{"/bin/bash", "preChroot.sh"}, &path); err != nil {
		return err
	}

	if isVariant {
		for _, config := range variant.Configs {
			if err := SpawnContainer([]string{"arch-chroot", config}, &path); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Configs {
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

	// TODO-4
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", "`/bin/bash findPackageManager.sh`", strings.Join(variant.Packages, ",")}, &path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", pkgManager, strings.Join(distribution.Packages, ",")}, &path); err != nil {
		return err
	}

	if err := SpawnContainer([]string{"/bin/bash", "postChroot.sh"}, &path); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, outDir string) error {
	err := IsDistro(distro)
	if err != nil {
		flag.Usage()
		return err
	}

	basePath := outDir + "/" + distro
	path := [2]string{basePath, "/root/" + distro}

	log.Println("Building: ", distro, "in dir: ", basePath)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	if !isAndroid {
		archi := IsValidArchitecture()
		if archi == nil {
			return err
		}

		if err := DownloadURLfromTags(path); err != nil {
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

	// Sets default for opensuse build
	if distro == "opensuse" {
		distro = "leap"
	}

	if distro == "" {
		flag.Usage()
	} else {
		Factory(distro, basepath)
	}
}

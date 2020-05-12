package main

// TODO : Handle Staging packages installation

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	strip "github.com/grokify/html-strip-tags-go"
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
	hekate, staging      = "hekate", "staging"
	dockerImageName      = "docker.io/library/ubuntu:18.04"

	baseJSON, _ = ioutil.ReadFile("./setup/base.json")
	basesDistro = []Distribution{}
	_           = json.Unmarshal([]byte(baseJSON), &basesDistro)
)

// SpawnContainer : Spawns a container based on dockerImageName
func SpawnContainer(cmd, env []string, volume *[2]string) error {
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	// cli.ImageBuild(ctx)

	reader, err := cli.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: dockerImageName,
		Cmd:   cmd,
		Env:   env,
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
		if buildarch == archis {
			fmt.Println("Found valid architecture: ", buildarch)
			return &buildarch
		}
	}
	fmt.Println(buildarch, "is not a valid architecture !")
	return nil
}

// CliSelector : Select an item in a menu froim cli
func CliSelector(label string, items []string) string {
	var inputValue string
	prompt := &survey.Select{
		Message: label,
		Options: items,
	}
	survey.AskOne(prompt, &inputValue)
	return inputValue
}

// WalkURL :
func WalkURL(source, regex string) []string {
	resp, err := http.Get(source)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK) {
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sanitizedHTML := strip.StripTags(string(bodyBytes))
	search, _ := regexp.Compile(regex)
	query := search.FindAllString(sanitizedHTML, -1)

	return query
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(path [2]string) (err error) {
	var constructedURL string

	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		if strings.Contains(avalaibleMirror, "{VERSION}") || strings.Contains(avalaibleMirror, "{BUILDARCH}") {
			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			regexURL := WalkURL(constructedURL, "(?m)^([[:digit:]]{1,3}.[[:digit:]]+|[[:digit:]]+)(?:/)")
			version := CliSelector("Select a version: ", regexURL)

			if version == "" {
				return err
			}

			constructedURL = strings.Replace(avalaibleMirror, "{VERSION}/", version, 1)
			regexURL = WalkURL(constructedURL, ".*.raw.xz")

			imageFile := CliSelector("Select an image file: ", regexURL)
			if imageFile == "" {
				return err
			}

			constructedURL = constructedURL + imageFile

		} else {
			constructedURL = avalaibleMirror
		}

		if _, err := url.ParseRequestURI(constructedURL); err != nil {
			fmt.Println("Couldn't found mirror:", constructedURL)
			return err
		}
	}
	fmt.Println("Mirror URL selected : ", constructedURL)

	cmd := exec.Command("/bin/bash", "./tools/prepare.sh", constructedURL, path[0])
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path [2]string) error {
	if err := SpawnContainer([]string{"/bin/bash", "/tools/preChroot.sh"}, nil, &path); err != nil {
		return err
	}

	if isVariant {
		for _, config := range variant.Configs {
			if err := SpawnContainer([]string{"arch-chroot", config, path[1]}, nil, &path); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Configs {
		if err := SpawnContainer([]string{"arch-chroot", config, path[1]}, nil, &path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"/bin/bash", "/tools/postChroot.sh"}, nil, &path); err != nil {
		return err
	}

	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(path [2]string) error {
	if err := SpawnContainer([]string{"/bin/bash", "/tools/preChroot.sh"}, nil, &path); err != nil {
		return err
	}

	// TODO-4
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", "`/bin/bash /tools/findPackageManager.sh`", strings.Join(variant.Packages, ","), path[1]}, nil, &path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", "`/bin/bash /tools/findPackageManager.sh`", strings.Join(distribution.Packages, ","), path[1]}, nil, &path); err != nil {
		return err
	}

	if err := SpawnContainer([]string{"/bin/bash", "/tools/postChroot.sh"}, nil, &path); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, outDir string) (err error) {
	if err := IsDistro(distro); err != nil {
		flag.Usage()
		return err
	}

	basePath := outDir + "/" + distro

	fmt.Println("Building:", distro, "\nInside directory:", basePath)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	if !isAndroid {
		path := [2]string{basePath, "/root/" + distro}
		if archi := IsValidArchitecture(); archi == nil {
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

		if err := SpawnContainer([]string{"/bin/bash", "/tools/createImage.sh", hekate, baseName}, nil, &path); err != nil {
			return err
		}
	} else {
		path := [2]string{basePath, "/root/android"}
		dockerImageName = "pablozaiden/switchroot-android-build:1.0.0"
		if err := SpawnContainer(nil, []string{"ROM_NAME=" + distro}, &path); err != nil {
			return err
		}
	}
	fmt.Println("Done!")
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

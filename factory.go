package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
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
	distribution              Distribution
	variant                   Variant
	baseName, buildarch       string
	imageFile, packageManager string
	isVariant, isAndroid      = false, false
	hekate, staging           bool
	prepare, configs          bool
	packages, image           bool

	managerList = []string{"zypper", "dnf", "yum", "pacman", "apt"}

	dockerImageName = "docker.io/library/ubuntu:18.04"
	baseJSON, _     = ioutil.ReadFile("./base.json")
	basesDistro     = []Distribution{}
	_               = json.Unmarshal([]byte(baseJSON), &basesDistro)

	hekateVersion = "5.2.0"
	nyxVersion    = "0.9.0"
	hekateBin     = "hekate_ctcaer_" + hekateVersion + ".bin"
	hekateURL     = "https://github.com/CTCaer/hekate/releases/download/v" + hekateVersion + "/hekate_ctcaer_" + hekateVersion + "_Nyx_" + nyxVersion + ".zip"
	hekateZip     = hekateURL[strings.LastIndex(hekateURL, "/")+1:]
)

// DetectPackageManager :
func DetectPackageManager() (err error) {
	for _, man := range managerList {
		if _, err := os.Stat("/usr/bin/" + man); os.IsExist(err) {
			packageManager = man
			return nil
		}
	}
	if packageManager == "zypper" {
		packageManager = packageManager + " up" + packageManager + " in"
	} else if packageManager == "dnf" || packageManager == "yum" || packageManager == "apt" {
		packageManager = packageManager + " update" + packageManager + " install"
	} else if packageManager == "pacman" {
		packageManager = packageManager + " -Syu"
	}
	return err
}

/* Rootfs Image creation
* Chroot into the filesystem
 */

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
			return &buildarch
		}
	}
	return nil
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(dst string) (image string, err error) {
	var constructedURL string
	var versions, images []string

	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		if strings.Contains(avalaibleMirror, "{VERSION}") {
			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			versionBody := WalkURL(constructedURL)

			search, _ := regexp.Compile(">:?([[:digit:]]{1,3}.[[:digit:]]+|[[:digit:]]+)(?:/)")
			match := search.FindAllStringSubmatch(*versionBody, -1)

			for i := 0; i < len(match); i++ {
				for _, submatches := range match {
					versions = append(versions, submatches[1])
				}
			}

			version, err := CliSelector("Select a version: ", versions)
			if err != nil {
				return "", err
			}

			constructedURL = strings.Replace(avalaibleMirror, "{VERSION}", version, 1)
			imageBody := WalkURL(constructedURL)

			search, _ = regexp.Compile(">:?([[:alpha:]]+.*.raw.xz)")
			imageMatch := search.FindAllStringSubmatch(*imageBody, -1)

			for i := 0; i < len(imageMatch); i++ {
				for _, submatches := range imageMatch {
					images = append(images, submatches[1])
				}
			}

			if len(images) > 1 {
				imageFile, err = CliSelector("Select an image file: ", images)
				if err != nil {
					return "", err
				}
			} else {
				imageFile = images[0]
			}

			imageFile = strings.TrimSpace(imageFile)
			constructedURL = constructedURL + imageFile
			image = imageFile

		} else {
			constructedURL = avalaibleMirror
		}

		if _, err := url.ParseRequestURI(constructedURL); err != nil {
			fmt.Println("Couldn't found mirror:", constructedURL)
			return "", err
		}

		if _, err := os.Stat(dst + "/" + image); os.IsNotExist(err) {
			if err := DownloadFile(constructedURL, dst); err != nil {
				return "", err
			}
		}
	}
	return image, nil
}

// PrepareFiles :
func PrepareFiles(basePath string) (err error) {
	if err = os.MkdirAll(basePath+"/tmp/", os.ModeDir); err != nil {
		return err
	}

	if err = os.MkdirAll(basePath+"/disk/", os.ModeDir); err != nil {
		return err
	}

	if err = os.MkdirAll(basePath+"/downloadedFiles/", os.ModeDir); err != nil {
		return err
	}

	if hekate {
		if _, err := os.Stat(basePath + "/downloadedFiles/" + hekateZip); os.IsNotExist(err) {
			fmt.Println("Downloading:", hekateZip)
			if err := DownloadFile(hekateURL, basePath+"/downloadedFiles/"+hekateZip); err != nil {
				return err
			}
		}
	}

	image, err := DownloadURLfromTags(basePath + "/downloadedFiles")
	if err != nil {
		return err
	}

	fmt.Println("Extracting:", image, "in:", basePath+"/disk")
	if strings.Contains(basePath+"/downloadedFiles/"+image, ".raw") {
		if _, err := os.Stat(basePath + "/downloadedFiles/" + image[0:strings.LastIndex(image, ".")]); os.IsNotExist(err) {
			if err := ExtractFiles(basePath+"/downloadedFiles/"+image, basePath+"/downloadedFiles/"); err != nil {
				return err
			}
		}

		image = image[0:strings.LastIndex(image, ".")]
		if _, err := MountImage(basePath+"/downloadedFiles/"+image, basePath); err != nil {
			return err
		}

		if _, err := DiskCopy(basePath+"/*", basePath+"/disk/"); err != nil {
			return err
		}

		if _, err := Unmount(basePath); err != nil {
			return err
		}
	} else {
		if err := ExtractFiles(basePath+"/downloadedFiles/"+image, basePath+"/disk"); err != nil {
			return err
		}

		fmt.Println("Done preparing files!")
	}
	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path [2]string) error {
	if err := PreChroot(path); err != nil {
		return err
	}

	if isVariant {
		for _, config := range variant.Configs {
			if err := SpawnContainer([]string{"arch-chroot", config, path[1]}, nil, path); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Configs {
		if err := SpawnContainer([]string{"arch-chroot", config, path[1]}, nil, path); err != nil {
			return err
		}
	}

	if err := PostChroot(path); err != nil {
		return err
	}

	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(path [2]string) error {
	if err := PreChroot(path); err != nil {
		return err
	}

	if distribution.Name == "arch" {
		if err := SpawnContainer([]string{"arch-chroot", "pacman-key --init", "&&", "pacman-key --populate archlinuxarm", path[1]}, nil, path); err != nil {
			return err
		}
	}

	// TODO-3 : Handle staging packages
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", packageManager, strings.Join(variant.Packages, ","), path[1]}, nil, path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", packageManager, strings.Join(distribution.Packages, ","), path[1]}, nil, path); err != nil {
		return err
	}

	if err := PostChroot(path); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, dst string) (err error) {
	basePath := dst + "/" + distro

	if err := IsDistro(distro); err != nil {
		flag.Usage()
		return err
	}

	if !isAndroid {
		fmt.Println("Building:", distro, "\nInside directory:", basePath)
		path := [2]string{basePath, "/root/" + distro}

		if archi := IsValidArchitecture(); archi == nil {
			return err
		}

		if prepare {
			if err := PrepareFiles(basePath); err != nil {
				return err
			}
		}

		if configs {
			if err := InstallPackagesInChrootEnv(path); err != nil {
				return err
			}
		}

		if packages {
			if err := ApplyConfigsInChrootEnv(path); err != nil {
				return err
			}
		}

		if image {
			var imageFile string

			if isVariant {
				CreateDisk(variant.Name+".img", basePath, "ext4")
				if _, err := MountImage(basePath+"/"+variant.Name+".img", basePath+"/tmp"); err != nil {
					return err
				}
				imageFile = variant.Name + ".img"
			} else {
				CreateDisk(baseName+".img", basePath, "ext4")
				if _, err := MountImage(basePath+"/"+baseName+".img", basePath+"/tmp"); err != nil {
					return err
				}
				imageFile = baseName + ".img"
			}

			if _, err := DiskCopy(basePath+"/disk/*", basePath+"/tmp/"); err != nil {
				return err
			}

			if hekate {
				if _, err := DiskCopy(basePath+"/tmp/boot/bootloader", basePath); err != nil {
					return err
				}

				if _, err := DiskCopy(basePath+"/tmp/boot/switchroot", basePath); err != nil {
					return err
				}

				if err := os.RemoveAll(basePath + "/tmp/boot/bootloader"); err != nil {
					return err
				}

				if err := os.RemoveAll(basePath + "/tmp/boot/switchroot"); err != nil {
					return err
				}

				if err := ExtractFiles(basePath+hekateZip, basePath); err != nil {
					return err
				}

				if _, err := DiskCopy(basePath+hekateBin, basePath+"/tmp/lib/firmware/reboot_payload.bin"); err != nil {
					return err
				}

				if _, err := Unmount(basePath + "/tmp/"); err != nil {
					return err
				}

				if err := SplitFile(basePath+"/"+imageFile, basePath+"/switchroot/install", 4290772992); err != nil {
					return err
				}
				// TODO - 4 : Implement 7z compression

			} else {
				if _, err := Unmount(basePath + "/tmp/"); err != nil {
					return err
				}
			}
		}
	} else {
		path := [2]string{basePath, "/root/android"}
		dockerImageName = "pablozaiden/switchroot-android-build:1.0.0"
		if err := SpawnContainer(nil, []string{"ROM_NAME=" + distro}, path); err != nil {
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
	flag.BoolVar(&hekate, "hekate", false, "Build an hekate installable filesystem")
	flag.BoolVar(&staging, "staging", false, "Install built local packages")
	flag.BoolVar(&prepare, "prepare", false, "Build an hekate installable filesystem")
	flag.BoolVar(&configs, "configs", false, "Install built local packages")
	flag.BoolVar(&packages, "packages", false, "Build an hekate installable filesystem")
	flag.BoolVar(&image, "image", false, "Install built local packages")
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

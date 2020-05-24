package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/mholt/archiver/v3"
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
	distribution Distribution
	variant      Variant

	buildarch       string
	managerList     = []string{"zypper", "dnf", "yum", "pacman", "apt"}
	dockerImageName = "docker.io/alizkan/jet-factory:1.0.0"
	hekateVersion   = "5.2.0"
	nyxVersion      = "0.9.0"
	hekateBin       = "hekate_ctcaer_" + hekateVersion + ".bin"
	hekateURL       = "https://github.com/CTCaer/hekate/releases/download/v" + hekateVersion + "/hekate_ctcaer_" + hekateVersion + "_Nyx_" + nyxVersion + ".zip"
	hekateZip       = hekateURL[strings.LastIndex(hekateURL, "/")+1:]

	isVariant, isAndroid     = false, false
	hekate, staging, prepare bool
	configs, packages, image bool

	baseJSON, _ = ioutil.ReadFile("./base.json")
	basesDistro = []Distribution{}
	_           = json.Unmarshal([]byte(baseJSON), &basesDistro)
)

// DetectPackageManager :
func DetectPackageManager() (packageManager string, err error) {
	for _, man := range managerList {
		if _, err := os.Stat("/usr/bin/" + man); os.IsExist(err) {
			packageManager = man
		}
	}
	if packageManager == "zypper" {
		packageManager = packageManager + " up" + packageManager + " in"
	} else if packageManager == "dnf" || packageManager == "yum" || packageManager == "apt" {
		packageManager = packageManager + " update" + packageManager + " install"
	} else if packageManager == "pacman" {
		packageManager = packageManager + " -Syu"
	} else {
		return "", errors.New("Can't detect package manager")
	}
	return packageManager, nil
}

/* Rootfs Image creation
* Chroot into the filesystem
 */

// IsDistro : Checks if a distribution is avalaible in the config files
func IsDistro(name string) (err error) {
	// Check if name match a known distribution
	for i := 0; i < len(basesDistro); i++ {
		if name == basesDistro[i].Name {
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

// SelectVersion : Retrieve a URL for a distribution based on a version
func SelectVersion() (constructedURL string, err error) {
	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		if strings.Contains(avalaibleMirror, "{VERSION}") {
			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			versionBody := WalkURL(constructedURL)

			search, _ := regexp.Compile(">:?([[:digit:]]{1,3}.[[:digit:]]+|[[:digit:]]+)(?:/)")
			match := search.FindAllStringSubmatch(*versionBody, -1)
			if match == nil {
				return "", errors.New("Couldn't match regex")
			}

			versions := make([]string, 0)
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

			// TODO : Extend to fit any archive extensions
			search, _ = regexp.Compile(">:?([[:alpha:]]+.*.raw.xz)")
			imageMatch := search.FindAllStringSubmatch(*imageBody, -1)
			images := make([]string, 0)
			for i := 0; i < len(imageMatch); i++ {
				for _, submatches := range imageMatch {
					images = append(images, submatches[1])
				}
			}

			var imageFile string
			if len(images) > 1 {
				imageFile, err = CliSelector("Select an image file: ", images)
				if err != nil {
					return "", err
				}
			} else if len(images) == 1 {
				imageFile = images[0]
			} else {
				return "", err
			}

			imageFile = strings.TrimSpace(imageFile)
			constructedURL = constructedURL + imageFile

		} else {
			constructedURL = avalaibleMirror
		}
	}
	return constructedURL, nil
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(dst string) (image string, err error) {
	constructedURL, err := SelectVersion()
	if err != nil {
		return "", err
	}

	if _, err := url.ParseRequestURI(constructedURL); err != nil {
		fmt.Println("Couldn't found mirror:", constructedURL)
		return "", err
	}

	err = DownloadFile(constructedURL, dst)
	if err != nil {
		return "", err
	}

	parsedURL := strings.Split(constructedURL, "/")
	image = parsedURL[len(parsedURL)-1]
	return image, nil
}

// PrepareFiles :
func PrepareFiles(basePath string) (err error) {
	disk := basePath + "/disk/"
	dlDir := basePath + "/downloadedFiles/"

	if err = os.MkdirAll(disk, 0755); err != nil {
		return err
	}

	if err = os.MkdirAll(dlDir, 0755); err != nil {
		return err
	}

	if hekate {
		if _, err := os.Stat(dlDir + hekateZip); os.IsNotExist(err) {
			fmt.Println("Downloading:", hekateZip)
			if err := DownloadFile(hekateURL, dlDir+hekateZip); err != nil {
				return err
			}
		}
		fmt.Println("Hekate downloaded !")
	}

	image, err := DownloadURLfromTags(dlDir)
	if err != nil {
		return err
	}

	if strings.Contains(dlDir+image, ".raw") {
		fmt.Println("Extracting:", image, "in:", dlDir)
		if err := ExtractFiles(dlDir+image, dlDir); err != nil {
			return err
		}
		fmt.Println("File successfully extracted to :", dlDir)

		image = image[0:strings.LastIndex(image, ".")]
		if _, err := CopyFromDisk(dlDir+image, disk); err != nil {
			return err
		}

	} else {
		fmt.Println("Extracting:", image, "in:", disk)
		if err := ExtractFiles(dlDir+image, disk); err != nil {
			return err
		}
	}
	fmt.Println("Done preparing files!")
	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(path string) error {
	if err := PreChroot(path); err != nil {
		return err
	}

	packageManager, err := DetectPackageManager()
	if err != nil {
		return err
	}

	if distribution.Name == "arch" {
		if err := SpawnContainer([]string{"arch-chroot", "pacman-key --init", "&&", "pacman-key --populate archlinuxarm", path}, nil, path); err != nil {
			return err
		}
	}

	// TODO : Handle staging packages
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", packageManager, strings.Join(variant.Packages, ","), path}, nil, path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", packageManager, strings.Join(distribution.Packages, ","), path}, nil, path); err != nil {
		log.Println("Chroot", err)
		return err
	}

	if err := PostChroot(path); err != nil {
		return err
	}

	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path string) error {
	if err := PreChroot(path); err != nil {
		return err
	}

	if isVariant {
		for _, config := range variant.Configs {
			if err := SpawnContainer([]string{"arch-chroot", config, path}, nil, path); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Configs {
		if err := SpawnContainer([]string{"arch-chroot", config, path}, nil, path); err != nil {
			return err
		}
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
		path := "/root/" + distro + "/disk/"

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

			disk := basePath + "/disk/"
			dlDir := basePath + "/downloadedFiles/"

			if isVariant {
				CreateDisk(variant.Name+".img", disk, basePath, "ext4")
				imageFile = variant.Name + ".img"
			} else {
				CreateDisk(distribution.Name+".img", disk, basePath, "ext4")
				imageFile = distribution.Name + ".img"
			}

			if _, err := CopyToDisk(imageFile, disk); err != nil {
				return err
			}

			if hekate {
				if err := ExtractFiles(dlDir+hekateZip, basePath); err != nil {
					return err
				}

				if err := Copy(basePath+hekateBin, basePath+"/mnt/lib/firmware/reboot_payload.bin"); err != nil {
					return err
				}

				if err := CopyDirectory(basePath+"/disk/boot/bootloader", basePath); err != nil {
					return err
				}

				if err := CopyDirectory(basePath+"/disk/boot/switchroot", basePath); err != nil {
					return err
				}

				if err := os.RemoveAll(basePath + "/mnt/boot/bootloader"); err != nil {
					return err
				}

				if err := os.RemoveAll(basePath + "/mnt/boot/switchroot"); err != nil {
					return err
				}

				if err := SplitFile(basePath+"/"+imageFile, basePath+"/switchroot/install/", 4290772992); err != nil {
					return err
				}
				archiver.Archive([]string{basePath + "/switchroot/", basePath + "/bootloader/"}, basePath+"/"+distro+".rar")
			}
		}
	} else {
		path := "/root/android"
		dockerImageName = "docker.io/library/pablozaiden/switchroot-android-build:1.0.0"
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

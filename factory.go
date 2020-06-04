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
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

	isVariant, isAndroid         = false, false
	hekate, staging, skip, force bool

	baseJSON, _ = ioutil.ReadFile("./base.json")
	basesDistro = []Distribution{}
	_           = json.Unmarshal([]byte(baseJSON), &basesDistro)
)

// DetectPackageManager :
func DetectPackageManager(rootfs string) (packageManager []string, err error) {
	for _, man := range managerList {
		if Exists(rootfs + "/usr/bin/" + man) {
			if man == "zypper" || man == "dnf" || man == "yum" || man == "apt" {
				packageManager = []string{man, "install", "-y"}
			} else if man == "pacman" {
				packageManager = []string{man, "-Syu", "--noconfirm", "--cachedir=/tmp"}
			} else {
				return nil, errors.New("Couldn't detect package manager")
			}
		}
	}
	return packageManager, nil
}

// SetDistro : Checks if a distribution is avalaible in the config files
func SetDistro(name string) (err error) {
	// Check/ if name match a known distribution
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

// SelectDistro :
func SelectDistro() (*string, error) {
	var avalaibles []string
	for _, baseDistro := range basesDistro {
		for _, variantDistro := range baseDistro.Variants {
			avalaibles = append(avalaibles, variantDistro.Name)
		}
		avalaibles = append(avalaibles, baseDistro.Name)
	}

	name, err := CliSelector("", avalaibles)
	if err != nil {
		return nil, err
	}

	return &name, nil
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
func DownloadURLfromTags(srcURL, dst string) error {
	err := RetryFunction(5, 2*time.Second, func() (err error) {
		_, err = url.ParseRequestURI(srcURL)
		err = DownloadFile(srcURL, dst)

		return
	})
	if err != nil {
		return err
	}
	return nil
}

// PrepareFiles : Prepare the filesystem for chroot
func PrepareFiles(basePath, dlDir string) (err error) {
	if err = os.MkdirAll(basePath, 0755); err != nil {
		return err
	}

	if err = os.MkdirAll(dlDir, 0755); err != nil {
		return err
	}

	if !skip {
		srcURL, err := SelectVersion()
		if err != nil {
			return err
		}

		parsedURL := strings.Split(srcURL, "/")
		image := parsedURL[len(parsedURL)-1]

		if _, err := os.Stat(dlDir + image); os.IsNotExist(err) || force == true {
			fmt.Println("Downloading rootfs !")
			err = DownloadURLfromTags(srcURL, dlDir)
			if err != nil {
				return err
			}
		}

		if hekate {
			if err := DownloadFile(hekateURL, dlDir+hekateZip); err != nil {
				return err
			}

			if err := ExtractFiles(dlDir+hekateZip, dlDir); err != nil {
				return err
			}
		}

		if strings.Contains(dlDir+image, ".raw") {
			if err := ExtractFiles(dlDir+image, dlDir); err != nil {
				return err
			}

			image = image[0:strings.LastIndex(image, ".")]
			if _, err := CopyFromDisk(dlDir+image, basePath); err != nil {
				return err
			}

			if err := os.Remove(dlDir + image); err != nil {
				return err
			}

		} else {
			if err := ExtractFiles(dlDir+image, basePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// InstallPackagesInChrootEnv : Installs packages list; Returns nil if successful
func InstallPackagesInChrootEnv(path string) error {
	packageManager, err := DetectPackageManager(path)
	if err != nil {
		return err
	}

	if distribution.Name == "arch" {
		err = SpawnContainer([]string{"arch-chroot", path, "pacman-key", "--init"}, nil, path)
		if err != nil {
			return err
		}

		err = SpawnContainer([]string{"arch-chroot", path, "pacman-key", "--populate", "archlinuxarm"}, nil, path)
		if err != nil {
			return err
		}

		// if buildarch == "aarch64" {
		// 	read, err := ioutil.ReadFile(path + "/etc/pacman.d/mirrorlist")
		// 	if err != nil {
		// 		return err
		// 	}

		// 	newContents := strings.Replace(string(read), "http://mirror.archlinuxarm.org", "http://de.mirror.archlinuxarm.org", -1)

		// 	err = ioutil.WriteFile(path+"/etc/pacman.d/mirrorlist", []byte(newContents), 0)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	if isVariant {
		for _, pkg := range variant.Packages {
			err := SpawnContainer([]string{"arch-chroot", path, packageManager[0], packageManager[1], packageManager[2], packageManager[3], pkg}, nil, path)
			if err != nil {
				return err
			}
		}
	}

	for _, pkg := range distribution.Packages {
		err = SpawnContainer([]string{"arch-chroot", path, packageManager[0], packageManager[1], packageManager[2], packageManager[3], pkg}, nil, path)
		if err != nil {
			return err
		}
	}
	// TODO : Handle staging packages
	return nil
}

// ApplyConfigsInChrootEnv : Runs one or multiple command in a chroot environment; Returns nil if successful
func ApplyConfigsInChrootEnv(path string) error {
	if isVariant {
		for _, config := range variant.Configs {
			var args string
			command := strings.Split(config, " ")
			if len(command) < 2 {
				args = ""
			}
			for _, arg := range command {
				args += arg + " "
			}
			if err := SpawnContainer([]string{"arch-chroot", path, command[0], args}, nil, path); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Configs {
		var args string
		command := strings.Split(config, " ")
		if len(command) < 2 {
			args = ""
		}
		for _, arg := range command {
			args += arg + " "
		}
		if err := SpawnContainer([]string{"arch-chroot", path, command[0], args}, nil, path); err != nil {
			return err
		}
	}

	return nil
}

// Hekate : Create a Hekate installable filesystem
func Hekate(dlDir, basePath, imageFile, distro string) error {
	if err := Copy(dlDir+hekateBin, basePath+"/lib/firmware/reboot_payload.bin"); err != nil {
		return err
	}

	if _, err := CopyToDisk(imageFile, basePath); err != nil {
		return err
	}

	if err := CopyDirectory(basePath+"/boot/bootloader", basePath); err != nil {
		return err
	}

	if err := CopyDirectory(basePath+"/boot/switchroot", basePath); err != nil {
		return err
	}

	if err := os.RemoveAll(basePath + "/boot/bootloader"); err != nil {
		return err
	}

	if err := os.RemoveAll(basePath + "/boot/switchroot"); err != nil {
		return err
	}

	if err := SplitFile(basePath+"/"+imageFile, basePath+"/switchroot/install/", 4290772992); err != nil {
		return err
	}

	err := archiver.Archive([]string{basePath + "/switchroot/", basePath + "/bootloader/"}, basePath+"/"+distro+".rar")
	if err != nil {
		return err
	}
	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string) (err error) {
	var imageFile string
	if distro == "" {
		sel, err := SelectDistro()
		if err != nil {
			return err
		}
		distro = *sel
	} else if distro == "lineage" {
		// Sets default for android build
		distro = "icosa"
		isAndroid = true
	} else if distro == "opensuse" {
		// Sets default for opensuse build
		distro = "leap"
	}

	err = SetDistro(distro)
	if err != nil {
		return err
	}

	basePath := "/root/linux/" + distro

	if !isAndroid {
		dlDir := "/root/linux/downloadedFiles/"

		if archi := IsValidArchitecture(); archi == nil {
			return err
		}

		if err := PrepareFiles(basePath, dlDir); err != nil {
			log.Println(err)
			return err
		}

		if err = BinfmtSupport(); err != nil {
			return err
		}

		err = PreChroot(basePath)
		if err != nil {
			log.Println(err)
			return err
		}

		var walkFn = func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == "dev" {
				return filepath.SkipDir
			}

			if info.IsDir() && info.Name() == "proc" {
				return filepath.SkipDir
			}

			if info.IsDir() && info.Name() == "sys" {
				return filepath.SkipDir
			}

			if !info.IsDir() {
				if err := os.Chmod(file, 0755); err != nil {
					fmt.Println(err)
				}
			}

			return nil
		}

		if err = filepath.Walk(basePath, walkFn); err != nil {
			return err
		}

		if err := InstallPackagesInChrootEnv(basePath); err != nil {
			log.Println(err)
			return err
		}

		if err := ApplyConfigsInChrootEnv(basePath); err != nil {
			log.Println(err)
			return err
		}

		if isVariant {
			if _, err := CreateDisk(variant.Name, basePath, basePath, "ext4"); err != nil {
				return err
			}
			imageFile = basePath + "/" + variant.Name + ".img"
		} else {
			if _, err := CreateDisk(distribution.Name, basePath, basePath, "ext4"); err != nil {
				return err
			}
			imageFile = basePath + "/" + distribution.Name + ".img"
		}

		if hekate {
			if err := Hekate(dlDir, basePath, imageFile, distro); err != nil {
				return err
			}
		} else {
			if _, err := CopyToDisk(imageFile, basePath); err != nil {
				return err
			}
		}
		fmt.Println("\nDone!")
	} else {
		path := "/root/android"
		dockerImageName = "docker.io/pablozaiden/switchroot-android-build:1.0.0"
		if err := SpawnContainer(nil, []string{"ROM_NAME=" + distro}, path); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var distro string
	flag.StringVar(&distro, "distro", "", "Distribution to build")
	flag.StringVar(&buildarch, "archi", "aarch64", "Distribution to build")

	flag.BoolVar(&hekate, "hekate", false, "Build an hekate installable filesystem")
	flag.BoolVar(&staging, "staging", false, "Install built local packages")

	flag.BoolVar(&skip, "skip", false, "Skip file prepare")
	flag.BoolVar(&force, "force", false, "Force to redownload files")

	flag.Parse()

	Factory(distro)
}

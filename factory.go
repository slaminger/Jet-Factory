package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
)

type (
	// Distribution : Represent a distribution conatining a name, version, desktop environment and an optional list of packages
	Distribution struct {
		Name          string              `json:"name"`
		Pre           []string            `json:"pre"`
		Post          []string            `json:"post"`
		Packages      []string            `json:"packages"`
		Architectures map[string][]string `json:"buildarch"`
		Variants      []Variant           `json:"variants"`
	}

	// Variant : Represent a distribution variant
	Variant struct {
		Name     string   `json:"name"`
		Pre      []string `json:"pre"`
		Post     []string `json:"post"`
		Packages []string `json:"packages"`
	}
)

var (
	distribution Distribution
	variant      Variant

	buildarch, dockerImageName, distro string
	managerList                        = []string{"zypper", "dnf", "yum", "pacman", "apt"}
	hekateVersion, nyxVersion          = "5.3.0", "0.9.2"
	hekateBin                          = "hekate_ctcaer_" + hekateVersion + ".bin"
	hekateURL                          = "https://github.com/CTCaer/hekate/releases/download/v" + hekateVersion + "/hekate_ctcaer_" + hekateVersion + "_Nyx_" + nyxVersion + ".zip"
	hekateZip                          = hekateURL[strings.LastIndex(hekateURL, "/")+1:]

	isVariant, isAndroid         = false, false
	hekate, staging, skip, force bool

	baseJSON, _ = ioutil.ReadFile("./base.json")
	basesDistro = []Distribution{}
	_           = json.Unmarshal([]byte(baseJSON), &basesDistro)
)

// DetectPackageManager :
func DetectPackageManager(rootfs string) (packageManager string, err error) {
	for _, man := range managerList {
		if Exists(rootfs + "/usr/bin/" + man) {
			if man == "zypper" || man == "dnf" || man == "yum" || man == "apt" {
				packageManager = man + " install " + "-y "
			} else if man == "pacman" {
				packageManager = man + " -Syu " + "--noconfirm "
			} else {
				return "", errors.New("Couldn't detect package manager")
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
			distro = name
			distribution = Distribution{Name: basesDistro[i].Name, Architectures: basesDistro[i].Architectures, Pre: basesDistro[i].Pre, Post: basesDistro[i].Post, Packages: basesDistro[i].Packages}
			return nil
		}
		for j := 0; j < len(basesDistro[i].Variants); j++ {
			if name == basesDistro[i].Variants[j].Name {
				isVariant = true
				distro = name
				variant = Variant{Name: basesDistro[i].Variants[j].Name, Packages: basesDistro[i].Variants[j].Packages, Pre: basesDistro[i].Variants[j].Pre, Post: basesDistro[i].Variants[j].Post}
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

// PostChroot : Remove qemu-aarch64-static binary and unmount the binded directories
func PostChroot(mountpoint string, oldRootF *os.File) error {
	err := os.Remove(mountpoint + "/usr/bin/qemu-" + buildarch + "-static")
	if err != nil {
		return err
	}
	return nil
}

// PreChroot : Copy qemu-aarch64-static binary and mount bind the directories
func PreChroot(path string) error {
	err := Copy("/usr/bin/qemu-"+buildarch+"-static", path+"/usr/bin/qemu-"+buildarch+"-static")
	if err != nil {
		return err
	}

	if err := os.Chmod(path+"/usr/bin/qemu-"+buildarch+"-static", 0755); err != nil {
		return err
	}
	return nil
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(srcURL, dst string) error {
	err := RetryFunction(5, 2*time.Second, func() (err error) {
		_, err = url.ParseRequestURI(srcURL)
		if err != nil {
			return err
		}
		err = DownloadFile(dst, srcURL)
		if err != nil {
			return err
		}
		return
	})
	if err != nil {
		return err
	}
	return nil
}

// PrepareFiles : Prepare the filesystem for chroot
func PrepareFiles(basePath, dlDir string) (err error) {
	os.RemoveAll(basePath)

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
			err = DownloadURLfromTags(srcURL, dlDir+image)
			if err != nil {
				return err
			}
		}

		if err := ExtractFiles(dlDir+image, basePath); err != nil {
			return err
		}

		if strings.Contains(basePath+image, ".raw") {
			image = image[0:strings.LastIndex(image, ".")]
			if _, err := CopyFromDisk(basePath+"/"+image, basePath); err != nil {
				return err
			}

			if err = os.Remove(basePath + "/" + image); err != nil {
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
		err = SpawnContainer([]string{"arch-chroot", path, "pacman-key", "--init"}, nil)
		if err != nil {
			return err
		}

		err = SpawnContainer([]string{"arch-chroot", path, "pacman-key", "--populate", "archlinuxarm"}, nil)
		if err != nil {
			return err
		}

		read, err := ioutil.ReadFile(path + "/etc/pacman.conf")
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), "CheckSpace", "#CheckSpace", -1)

		err = ioutil.WriteFile(path+"/etc/pacman.conf", []byte(newContents), 0)
		if err != nil {
			return err
		}
	}

	var pkgs string
	if isVariant {
		for _, pkg := range variant.Packages {
			pkgs += pkg + " "
		}
	}

	for _, pkg := range distribution.Packages {
		pkgs += pkg + " "
	}

	err = SpawnContainer([]string{"arch-chroot", path, "bash", "-c", packageManager + " " + pkgs}, nil)
	if err != nil {
		return err
	}

	return nil
}

// PreConfigRootfs : Runs one or multiple command in a chroot environment; Returns nil if successful
func PreConfigRootfs(path string) error {
	if isVariant {
		for _, config := range variant.Pre {
			if err := SpawnContainer([]string{"arch-chroot", path, "bash", "-c", config}, nil); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Pre {
		if err := SpawnContainer([]string{"arch-chroot", path, "bash", "-c", config}, nil); err != nil {
			return err
		}
	}

	return nil
}

// PostConfigRootfs : Runs one or multiple command in a chroot environment; Returns nil if successful
func PostConfigRootfs(path string) error {
	if isVariant {
		for _, config := range variant.Post {
			if err := SpawnContainer([]string{"arch-chroot", path, "bash", "-c", config}, nil); err != nil {
				return err
			}
		}
	}

	for _, config := range distribution.Post {
		if err := SpawnContainer([]string{"arch-chroot", path, "bash", "-c", config}, nil); err != nil {
			return err
		}
	}

	return nil
}

// Hekate : Create a Hekate installable filesystem
func Hekate(dlDir, basePath, imagePath, distro string) error {
	if err := os.MkdirAll(dlDir+"switchroot/install/", 0755); err != nil {
		return err
	}

	if err := DownloadFile(dlDir+hekateZip, hekateURL); err != nil {
		return err
	}

	if err := ExtractFiles(dlDir+hekateZip, dlDir); err != nil {
		return err
	}

	if err := Copy(dlDir+hekateBin, basePath+"/lib/firmware/reboot_payload.bin"); err != nil {
		return err
	}

	args := "virt-make-fs --format=raw --type=ext4 --size=+256M " + basePath + " " + imagePath
	_, err := exec.Command("bash", "-c", args).Output()
	if err != nil {
		return err
	}

	if err := SplitFile(imagePath, dlDir+"switchroot/install/", 4290772992); err != nil {
		return err
	}

	if err := os.RemoveAll(dlDir + "SWC-" + distro + ".zip"); err != nil {
		return err
	}

	if err := archiver.DefaultZip.Archive([]string{dlDir + "switchroot/", dlDir + "bootloader/"}, dlDir+"SWC-"+distro+".zip"); err != nil {
		return err
	}
	fmt.Println("Hekate archive created in :", dlDir+"SWC-"+distro+".zip")

	if err := os.RemoveAll(dlDir + hekateZip); err != nil {
		return err
	}

	if err := os.RemoveAll(imagePath); err != nil {
		return err
	}

	if err := os.RemoveAll(dlDir + hekateBin); err != nil {
		return err
	}

	if err := os.RemoveAll(dlDir + "switchroot/"); err != nil {
		return err
	}

	if err := os.RemoveAll(dlDir + "bootloader/"); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory() (err error) {
	if distro == "" {
		sel, err := SelectDistro()
		if err != nil {
			return err
		}
		distro = *sel
	} else if distro == "opensuse" {
		// Sets default for opensuse build
		distro = "leap"
	} else if distro == "lineage" || distro == "icosa" || distro == "foster" || distro == "foster_tab" {
		// Sets default for lineage to icosa
		isAndroid = true
		if distro == "lineage" {
			distro = "icosa"
		}
	}

	if isAndroid {
		dockerImageName = "docker.io/pablozaiden/switchroot-android-build:latest"

		if err := os.MkdirAll("/root/android/lineage", 0755); err != nil {
			return err
		}

		if err := os.MkdirAll("/root/android/output", 0755); err != nil {
			return err
		}

		if err := SpawnContainer(nil, []string{"ROM_NAME=" + distro}); err != nil {
			return err
		}

		return nil
	}

	basePath := "/root/linux/" + distro
	dlDir := "/root/linux/downloadedFiles/"
	imagePath := dlDir + distro + ".img"
	dockerImageName = "docker.io/alizkan/jet-factory:1.0.0"

	if err := SetDistro(distro); err != nil {
		return err
	}

	if buildarch == "" {
		if err := SelectArchitecture(); err != nil {
			return err
		}
	}

	if archi := IsValidArchitecture(); archi == nil {
		return err
	}

	if err := PrepareFiles(basePath, dlDir); err != nil {
		return err
	}

	if err := BinfmtSupport(); err != nil {
		return err
	}

	if err := PreChroot(basePath); err != nil {
		return err
	}

	if err := PreConfigRootfs(basePath); err != nil {
		return err
	}

	if err := InstallPackagesInChrootEnv(basePath); err != nil {
		return err
	}

	if err := PostConfigRootfs(basePath); err != nil {
		return err
	}

	if hekate {
		if err := Hekate(dlDir, basePath, imagePath, distro); err != nil {
			return err
		}
	} else {
		args := "virt-make-fs --format=raw --type=ext4 --size=+256M " + basePath + " " + imagePath
		_, err := exec.Command("bash", "-c", args).Output()
		if err != nil {
			return err
		}
	}

	if err := os.RemoveAll(basePath); err != nil {
		return err
	}

	fmt.Println("\nDone")
	return nil
}

func main() {
	flag.StringVar(&distro, "distro", "", "Distribution to build")
	flag.StringVar(&buildarch, "archi", "", "Distribution to build")

	flag.BoolVar(&hekate, "hekate", false, "Build an hekate installable filesystem")
	// flag.BoolVar(&staging, "staging", false, "Install built local packages")
	flag.BoolVar(&skip, "skip", false, "Skip file prepare")
	flag.BoolVar(&force, "force", false, "Force to redownload files")

	flag.Parse()

	Factory()
}

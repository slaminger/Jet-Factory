package main

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
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Mirantis/virtlet/pkg/diskimage/guestfs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	strip "github.com/grokify/html-strip-tags-go"
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
	distribution         Distribution
	variant              Variant
	baseName, buildarch  string
	imageFile            string
	isVariant, isAndroid = false, false
	hekate, staging      = false, false

	dockerImageName = "docker.io/library/ubuntu:18.04"
	baseJSON, _     = ioutil.ReadFile("./setup/base.json")
	basesDistro     = []Distribution{}
	_               = json.Unmarshal([]byte(baseJSON), &basesDistro)

	hekateVersion = "5.2.0"
	nyxVersion    = "0.9.0"
	hekateBin     = "hekate_ctcaer_" + hekateVersion + ".bin"
	hekateURL     = "https://github.com/CTCaer/hekate/releases/download/v" + hekateVersion + "/hekate_ctcaer_" + hekateVersion + "_Nyx_" + nyxVersion + ".zip"
	hekateZip     = hekateURL[strings.LastIndex(hekateURL, "/")+1:]
)

/* Menu - CLI
* Selector function
 */

// CliSelector : Select an item in a menu froim cli
func CliSelector(label string, items []string) (answer *string, err error) {
	var inputValue string
	prompt := &survey.Select{
		Message: label,
		Options: items,
	}
	if err := survey.AskOne(prompt, &inputValue); err != nil {
		return nil, err
	}
	return &inputValue, nil
}

/* Net utilities
* WalkURL using regex
* Wget to download a file
 */

// WalkURL : Walk a URL using a regex, and return the matches
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

// Wget : Download a file in given path
func Wget(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

/* System utilities
* Docker Spawn container
* Archive extractor
* Guestfs mount, create disk
 */

// SpawnContainer : Spawns a container based on dockerImageName
func SpawnContainer(cmd, env []string, volume [2]string) error {
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

// MountImage :
func MountImage(imgFile, mountDir string) (*guestfs.GuestfsError, error) {
	disk := imgFile

	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}

	/* Attach the disk image read-only to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        true,
	}
	if err := g.Add_drive(disk, &optargs); err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	if err := g.Launch(); err != nil {
		return err, nil
	}

	/* Ask libguestfs to inspect for operating systems. */
	roots, err := g.Inspect_os()
	if err != nil {
		return err, nil
	}
	if len(roots) == 0 {
		panic("inspect-vm: no operating systems found")
	}

	var root string
	if len(roots) == 1 {
		root = roots[0]
	}

	if err := g.Mount(root, mountDir); err != nil {
		return err, nil
	}

	return nil, nil
}

// CreateDisk :
func CreateDisk(name, outDir, format string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	// TODO - 1: Create disk, use DU to get dir size
	// if err := g. (name, "raw", 512*1024*1024); err != nil {
	//return err
	//}

	/* Set the trace flag so that we can see each libguestfs call. */
	g.Set_trace(true)

	/* Attach the disk image to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	if err := g.Add_drive(name, &optargs); err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	if err := g.Launch(); err != nil {
		return err, nil
	}

	/* Get the list of devices.  Because we only added one drive
	 * above, we expect that this list should contain a single
	 * element.
	 */
	devices, err := g.List_devices()
	if err != nil {
		return err, nil
	}
	if len(devices) != 1 {
		panic("expected a single device from list-devices")
	}

	/* Partition the disk as one single MBR partition. */
	err = g.Part_disk(devices[0], "mbr")
	if err != nil {
		return err, nil
	}

	/* Get the list of partitions.  We expect a single element, which
	 * is the partition we have just created.
	 */
	partitions, err := g.List_partitions()
	if err != nil {
		return err, nil
	}
	if len(partitions) != 1 {
		panic("expected a single partition from list-partitions")
	}

	/* Create a filesystem on the partition. */
	err = g.Mkfs(format, partitions[0], nil)
	if err != nil {
		return err, nil
	}

	err = g.Close()
	if err != nil {
		return err, nil
	}

	return nil, nil
}

// CopyToDisk :
func CopyToDisk(mountedPath, outPath string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	err := g.Cp_r(mountedPath, outPath)
	if err != nil {
		return err, nil
	}

	if err = g.Shutdown(); err != nil {
		return err, nil
	}

	return nil, nil
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

// PreChroot : Copy qemu-aarch64-static binary and mount bind the directories
func PreChroot(mount [2]string) error {
	err := SpawnContainer(
		[]string{
			"cp", "/usr/bin/qemu-aarch64-static",
			mount[1] + "/usr/bin",

			"&&", "mount", "--bind",
			mount[1], mount[1],

			"&&", "mount", "--bind",
			mount[1] + "/bootloader",
			mount[1] + "/boot",
		},
		nil,
		mount,
	)
	if err != nil {
		return err
	}
	return nil
}

// PostChroot : Remove qemu-aarch64-static binary and unmount the binded directories
func PostChroot(mounted [2]string) error {
	err := SpawnContainer(
		[]string{
			"rm", mounted[1] + "/usr/bin/qemu-aarch64-static",
			"&&", "umount", mounted[1],
			"&&", "mount", mounted[1] + "/boot",
		},
		nil,
		mounted,
	)
	if err != nil {
		return err
	}
	return nil
}

// PrepareFiles :
func PrepareFiles(basePath string) (err error) {
	if err = os.MkdirAll(basePath+"/tmp", os.ModePerm); err != nil {
		return err
	}

	if err = os.MkdirAll(basePath+"/downloadedFiles", os.ModePerm); err != nil {
		return err
	}

	if hekate {
		if _, err := os.Stat(basePath + "/downloadedFiles/" + hekateZip); os.IsNotExist(err) {
			fmt.Println("Downloading:", hekateZip)
			if err := Wget(hekateURL, basePath+"/downloadedFiles/"+hekateZip); err != nil {
				return err
			}
		}
	}

	image := *DownloadURLfromTags(basePath + "/downloadedFiles")

	if _, err := os.Stat(basePath + "/downloadedFiles/" + image[0:strings.LastIndex(image, ".")]); os.IsNotExist(err) {
		fmt.Println("Extracting:", image, "in: ./tmp")
		// TODO-2 : Implement correctly extract function
		if err := archiver.Extract(basePath+"/downloadedFiles/"+image, "", basePath+"/tmp"); err != nil {
			log.Println(err)
			return err
		}
	}

	if strings.Contains(basePath+"/downloadedFiles/"+image, ".raw") {
		image = image[0:strings.LastIndex(image, ".")]
		if _, err := MountImage(basePath+"/downloadedFiles/"+image, basePath+"/tmp"); err != nil {
			return err
		}

		if _, err := CopyToDisk(basePath+"/tmp/*", basePath); err != nil {
			return err
		}
	}
	return nil
}

// DownloadURLfromTags : Retrieve a URL for a distribution based on a version
func DownloadURLfromTags(filepath string) (image *string) {
	var constructedURL string
	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		if strings.Contains(avalaibleMirror, "{VERSION}") || strings.Contains(avalaibleMirror, "{BUILDARCH}") {
			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			regexURL := WalkURL(constructedURL, "(?m)^([[:digit:]]{1,3}.[[:digit:]]+|[[:digit:]]+)(?:/)")

			version, err := CliSelector("Select a version: ", regexURL)
			if err != nil {
				return nil
			}

			constructedURL = strings.Replace(avalaibleMirror, "{VERSION}/", *version, 1)
			regexURL = WalkURL(constructedURL, ".*.raw.xz")

			imageFile, err := CliSelector("Select an image file: ", regexURL)
			if err != nil {
				return nil
			}
			*imageFile = strings.TrimSpace(*imageFile)
			constructedURL = constructedURL + *imageFile
			image = imageFile
		} else {
			constructedURL = avalaibleMirror
		}

		if _, err := url.ParseRequestURI(constructedURL); err != nil {
			fmt.Println("Couldn't found mirror:", constructedURL)
			return nil
		}
		if _, err := os.Stat(filepath + "/" + imageFile); os.IsNotExist(err) {
			fmt.Println("Mirror URL selected:", constructedURL)
			if err := Wget(constructedURL, filepath+"/"+imageFile); err != nil {
				return nil
			}
		}
	}
	return image
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

	// TODO-3 : Handle staging packages
	if isVariant {
		if err := SpawnContainer([]string{"arch-chroot", "`/bin/bash /tools/findPackageManager.sh`", strings.Join(variant.Packages, ","), path[1]}, nil, path); err != nil {
			return err
		}
	}

	if err := SpawnContainer([]string{"arch-chroot", "`/bin/bash /tools/findPackageManager.sh`", strings.Join(distribution.Packages, ","), path[1]}, nil, path); err != nil {
		return err
	}

	if err := PostChroot(path); err != nil {
		return err
	}

	return nil
}

// Factory : Build your distribution with the setted options; Returns a pointer on the location of the produced build
func Factory(distro string, outDir string) (err error) {
	basePath := outDir + "/" + distro

	if err := IsDistro(distro); err != nil {
		flag.Usage()
		return err
	}

	fmt.Println("Building:", distro, "\nInside directory:", basePath)
	if !isAndroid {
		path := [2]string{basePath, "/root/" + distro}

		if archi := IsValidArchitecture(); archi == nil {
			fmt.Println(buildarch, "is not a valid architecture !")
			return err
		}
		fmt.Println("Found valid architecture: ", buildarch)

		if err := PrepareFiles(basePath); err != nil {
			return err
		}

		if err := InstallPackagesInChrootEnv(path); err != nil {
			return err
		}

		if err := ApplyConfigsInChrootEnv(path); err != nil {
			return err
		}

		// CreateDisk("", "", "ext4")

		if hekate {
			// TODO - 4 : Implement split function and 7z
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
	flag.Bool("hekate", hekate, "Build an hekate installable filesystem")
	flag.Bool("staging", staging, "Install built local packages")
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

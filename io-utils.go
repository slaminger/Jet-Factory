package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Mirantis/virtlet/pkg/diskimage/guestfs"
	"github.com/mholt/archiver/v3"
	"github.com/xi2/xz"
)

// MountImage :
func MountImage(disk, mountDir string) (*guestfs.GuestfsError, error) {
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
		log.Println("inspect-vm: no operating systems found")
		return err, nil
	}

	var root string
	if len(roots) == 1 {
		root = roots[0]
	} else {
		root, _ = CliSelector("Select root partition to use:", roots)
		if root == "" {
			return err, nil
		}
	}

	mountDir, ok := filepath.Abs(mountDir)
	if ok != nil {
		return nil, ok
	}

	if err := g.Mount(root, mountDir); err != nil {
		return err, nil
	}

	return nil, nil
}

// CreateDisk :
func CreateDisk(disk, outDir, format string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	size, err := g.Du(outDir)
	if err != nil {
		return err, nil
	}

	fmt.Println("Estimated size:", size)

	if ret := exec.Command("dd", "of="+outDir+"/"+disk+".img", "bs=1", "count=0", "seek="+string(size*1024*1024)+"M"); ret != nil {
		return err, nil
	}

	/* Set the trace flag so that we can see each libguestfs call. */
	g.Set_trace(true)

	/* Attach the disk image to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	if err := g.Add_drive(disk, &optargs); err != nil {
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
		log.Println("expected a single device from list-devices")
		return err, nil
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
		log.Println("expected a single partition from list-partitions")
		return err, nil
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

// DiskCopy :
func DiskCopy(disk, dst string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	/* Attach the disk image to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	if err := g.Add_drive(disk, &optargs); err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	if err := g.Launch(); err != nil {
		return err, nil
	}

	err := g.Cp_r(disk, dst)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

// Unmount :
func Unmount(disk string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	/* Attach the disk image to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	if err := g.Add_drive(disk, &optargs); err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	if err := g.Launch(); err != nil {
		return err, nil
	}

	err := g.Umount(disk, nil)
	if err != nil {
		return err, nil
	}

	if err = g.Shutdown(); err != nil {
		return err, nil
	}
	return nil, nil
}

// ExtractFiles :
func ExtractFiles(archivePath, dst string) (err error) {
	if strings.Contains(archivePath, ".raw.xz") {
		data, err := ioutil.ReadFile(archivePath)
		if err != nil {
			return err
		}
		r, err := xz.NewReader(bytes.NewReader(data), 0)
		if err != nil {
			return err
		}

		outputFile := strings.Split(filepath.Base(archivePath), ".xz")[0]
		destination, err := os.Create(dst + "/" + outputFile)
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, r)
		if err != nil {
			return err
		}

		err = destination.Sync()
		if err != nil {
			return err
		}

	} else if strings.Contains(archivePath, ".tar") {
		if err := archiver.Extract(archivePath, "", dst); err != nil {
			return err
		}
	} else {
		fmt.Println("Couldn't recognize archive type for:", archivePath)
		return err
	}
	return nil
}

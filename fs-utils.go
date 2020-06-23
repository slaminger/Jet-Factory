package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Mirantis/virtlet/pkg/diskimage/guestfs"
	"github.com/codeclysm/extract"
	"github.com/xi2/xz"
)

// Copy : Copy a file to a dst; returns nil if success; returns err otherwise;
func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// RetryFunction : Retry a function that failed, n Times; returns nil on success; returns err otherwise;
func RetryFunction(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		fmt.Println("Retrying after error:", err)
	}
	return fmt.Errorf("After %d attempts, last error: %s", attempts, err)
}

// SplitFile : Split a file into pieces of sizeInBytes size; returns nil on success; returns err otherwise;
func SplitFile(filepath, outpath string, sizeInBytes int64) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(sizeInBytes)))

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(float64(sizeInBytes), float64(fileSize-int64(i*uint64(sizeInBytes)))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := "l4t.0" + strconv.FormatUint(i, 10)

		_, err := os.Create(outpath + fileName)
		if err != nil {
			return err
		}

		// write/save buffer to disk
		ioutil.WriteFile(outpath+fileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to : ", outpath+fileName)
	}
	return nil
}

// ExtractFiles : Extract compressed files, supports TAR, GZ, ZIP, TBZ2; returns nil on success; returns err otherwise;
func ExtractFiles(archivePath, dst string) (err error) {
	fmt.Println("Extracting:", archivePath, "in:", dst)

	if strings.Contains(archivePath, ".xz") {
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

	} else if strings.Contains(archivePath, ".tar") || strings.Contains(archivePath, ".zip") || strings.Contains(archivePath, ".rar") || strings.Contains(archivePath, ".tbz2") || strings.Contains(archivePath, ".bz2") {
		data, err := ioutil.ReadFile(archivePath)
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(data)

		filepath.Ext(archivePath)

		switch filepath.Ext(archivePath) {
		case ".bz2":
			err = extract.Bz2(context.Background(), buffer, dst, nil)
		case ".tbz2":
			err = extract.Bz2(context.Background(), buffer, dst, nil)
		case ".gz":
			err = extract.Gz(context.Background(), buffer, dst, nil)
		case ".zip":
			err = extract.Zip(context.Background(), buffer, dst, nil)
		case ".tar":
			err = extract.Tar(context.Background(), buffer, dst, nil)
		default:
			fmt.Println("Unknown error")
			return err
		}

		if err != nil {
			fmt.Println(archivePath, ": Should not fail: "+err.Error())
		}

	} else {
		fmt.Println("\nCouldn't recognize archive type for:", archivePath)
		return err
	}
	fmt.Println("Successfully extracted:", archivePath)
	return nil
}

// // CreateDisk :
// func CreateDisk(src, dst, diskImage, format string) (*guestfs.GuestfsError, error) {
// 	diskImage = dst + diskImage + ".img"
// 	size := DirSizeMB(src)
// 	rootSize := size + (256 * 1024 * 1024)

// 	fmt.Println("Estimated size:", rootSize/1024/1024, "MB")

// 	if err := os.Remove(diskImage); err != nil {
// 		return nil, err
// 	}

// 	// TODO : Create writable disk Image
// 	diskObj, ferr := diskfs.Create(diskImage, rootSize, diskfs.Raw)
// 	if ferr != nil {
// 		return nil, ferr
// 	}

// 	// create a partition table
// 	diskObj.LogicalBlocksize = 2048
// 	fspec := disk.FilesystemSpec{Partition: 0, FSType: filesystem.TypeISO9660, VolumeLabel: "label"}
// 	fs, err := diskObj.CreateFilesystem(fspec)
// 	if err != nil {
// 		return nil, err
// 	}

// 	iso, ok := fs.(*iso9660.FileSystem)
// 	if !ok {
// 		return nil, errors.New("not an iso9660 filesystem")
// 	}

// 	if err := iso.Finalize(iso9660.FinalizeOptions{}); err != nil {
// 		return nil, err
// 	}

// 	g, errno := guestfs.Create()
// 	if errno != nil {
// 		return nil, errno
// 	}
// 	defer g.Close()

// 	/* Attach the disk image to libguestfs. */
// 	optargs := guestfs.OptargsAdd_drive{
// 		Format_is_set:   true,
// 		Format:          "raw",
// 		Readonly_is_set: true,
// 		Readonly:        false,
// 	}
// 	if gerr := g.Add_drive(diskImage, &optargs); gerr != nil {
// 		return gerr, nil
// 	}

// 	/* Run the libguestfs back-end. */
// 	if gerr := g.Launch(); gerr != nil {
// 		return gerr, nil
// 	}

// 	/* Get the list of devices.  Because we only added one drive
// 	 * above, we expect that this list should contain a single
// 	 * element.
// 	 */
// 	devices, gerr := g.List_devices()
// 	if gerr != nil {
// 		return gerr, nil
// 	}
// 	if len(devices) != 1 {
// 		return nil, errors.New("\nexpected a single device from list-devices")
// 	}

// 	/* Partition the disk as one single MBR partition. */
// 	if gerr := g.Part_disk(devices[0], "mbr"); gerr != nil {
// 		return gerr, nil
// 	}

// 	/* Get the list of partitions.  We expect a single element, which
// 	 * is the partition we have just created.
// 	 */
// 	partitions, gerr := g.List_partitions()
// 	if gerr != nil {
// 		return gerr, nil
// 	}
// 	if len(partitions) != 1 {
// 		return nil, errors.New("\nexpected a single partition from list-partitions")
// 	}

// 	/* Create a filesystem on the partition. */
// 	if format == "ext4" {
// 		if gerr := g.Mke2fs(partitions[0], &guestfs.OptargsMke2fs{Fstype: format}); gerr != nil {
// 			return gerr, nil
// 		}
// 	} else if format == "fat32" {
// 		if gerr := g.Mkfs(format, partitions[0], nil); gerr != nil {
// 			return gerr, nil
// 		}
// 	} else {
// 		return nil, errors.New("filesystem type is currently not handled or does not exist")
// 	}

// 	if gerr := g.Shutdown(); gerr != nil {
// 		return gerr, nil
// 	}

// 	return nil, nil
// }

// CopyFromDisk : Use libguestfs backend to extract files from a partition; returns nil on success; returns err otherwise;
func CopyFromDisk(disk, dst string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	/* Attach the disk image read-only to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	if gerr := g.Add_drive(disk, &optargs); gerr != nil {
		return gerr, nil
	}

	/* Run the libguestfs back-end. */
	if gerr := g.Launch(); gerr != nil {
		return gerr, nil
	}

	/* Ask libguestfs to inspect for operating systems. */
	roots, gerr := g.Inspect_os()
	if gerr != nil {
		return gerr, nil
	}

	var root string
	var ferr error

	if len(roots) > 1 {
		root, ferr = CliSelect("Select root partition to use:", roots)
		if ferr != nil {
			return nil, ferr
		}
	} else if len(roots) == 1 {
		root = roots[0]
	} else {
		return nil, errors.New("inspect-vm: no operating systems found for")
	}

	if gerr := g.Mount(root, "/"); gerr != nil {
		return gerr, nil
	}

	if gerr := g.Tar_out("/", dst+"/"+distro+".tar", nil); gerr != nil {
		return gerr, nil
	}

	if gerr := g.Umount_all(); gerr != nil {
		return gerr, nil
	}

	if gerr := g.Shutdown(); gerr != nil {
		return gerr, nil
	}

	if err := ExtractFiles(dst+"/"+distro+".tar", dst); err != nil {
		return nil, err
	}

	if err := os.Remove(dst + "/" + distro + ".tar"); err != nil {
		return nil, err
	}

	return nil, nil
}

// CopyToDisk :
// func CopyToDisk(disk, src string) (*guestfs.GuestfsError, error) {
// 	g, errno := guestfs.Create()
// 	if errno != nil {
// 		return nil, errno
// 	}
// 	defer g.Close()

// 	/* Attach the disk image read-only to libguestfs. */
// 	optargs := guestfs.OptargsAdd_drive{
// 		Format_is_set:   true,
// 		Format:          "raw",
// 		Readonly_is_set: true,
// 		Readonly:        false,
// 	}
// 	if gerr := g.Add_drive(disk, &optargs); gerr != nil {
// 		return gerr, nil
// 	}

// 	/* Run the libguestfs back-end. */
// 	if gerr := g.Launch(); gerr != nil {
// 		return gerr, nil
// 	}

// 	/* Get the list of devices.  Because we only added one drive
// 	 * above, we expect that this list should contain a single
// 	 * element.
// 	 */
// 	devices, gerr := g.List_devices()
// 	if gerr != nil {
// 		return gerr, nil
// 	}
// 	if len(devices) != 1 {
// 		return nil, errors.New("\nexpected a single device from list-devices")
// 	}

// 	/* Get the list of partitions.  We expect a single element, which
// 	 * is the partition we have just created.
// 	 */
// 	partitions, gerr := g.List_partitions()
// 	if gerr != nil {
// 		return gerr, nil
// 	}
// 	if len(partitions) != 1 {
// 		return nil, errors.New("\nexpected a single partition from list-partitions")
// 	}

// 	if gerr := g.Mount(partitions[0], "/"); gerr != nil {
// 		return gerr, nil
// 	}

// 	// TODO : Fix CopyToDisk properly

// 	if gerr := g.Umount_all(); gerr != nil {
// 		return gerr, nil
// 	}

// 	if gerr := g.Shutdown(); gerr != nil {
// 		return gerr, nil
// 	}

// 	return nil, nil
// }

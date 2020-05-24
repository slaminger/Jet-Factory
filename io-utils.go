package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Mirantis/virtlet/pkg/diskimage/guestfs"
	"github.com/mholt/archiver/v3"
	"github.com/xi2/xz"
)

// CreateDisk :
func CreateDisk(disk, src, dst, format string) (*guestfs.GuestfsError, error) {
	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	/* Set the trace flag so that we can see each libguestfs call. */
	g.Set_trace(true)

	size := DirSizeMB(src)
	fmt.Println("Estimated size:", size, "MB")

	if ret := exec.Command("dd", "of="+dst+"/"+disk+".img", "bs=1", "count=0", "seek="+string(int64(size*1024*1024)/4+4+512)+"M"); ret != nil {
		err := errors.New("disk image creation failed")
		return nil, err
	}

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
		fmt.Println("expected a single device from list-devices")
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
		fmt.Println("expected a single partition from list-partitions")
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

// CopyFromDisk :
func CopyFromDisk(disk, dst string) (*guestfs.GuestfsError, error) {
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
		fmt.Println("inspect-vm: no operating systems found")
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

	if err := g.Mount(root, "/"); err != nil {
		return err, nil
	}

	if err := g.Mount_local("/mnt", nil); err != nil {
		return err, nil
	}

	go g.Mount_local_run()

	files, ok := filepath.Glob("/mnt/*")
	if ok != nil {
		return nil, ok
	}

	for _, dir := range files {
		if err := CopyDirectory(dir, dst); err != nil {
			return nil, err
		}
	}

	if err := g.Umount(root, nil); err != nil {
		return err, nil
	}

	if err = g.Shutdown(); err != nil {
		return err, nil
	}
	return nil, nil
}

// CopyToDisk :
func CopyToDisk(disk, src string) (*guestfs.GuestfsError, error) {
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
		fmt.Println("inspect-vm: no operating systems found")
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

	if err := g.Mount(root, "/"); err != nil {
		return err, nil
	}

	if err := g.Mount_local("/mnt", nil); err != nil {
		return err, nil
	}

	go g.Mount_local_run()

	if err := CopyDirectory(src, "/mnt/"); err != nil {
		return nil, err
	}

	if err := g.Umount(root, nil); err != nil {
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

	} else if strings.Contains(archivePath, ".tar") || strings.Contains(archivePath, ".zip") || strings.Contains(archivePath, ".rar") {
		if err := archiver.Unarchive(archivePath, dst); err != nil {
			return err
		}
	} else {
		fmt.Println("Couldn't recognize archive type for:", archivePath)
		return err
	}
	return nil
}

// SplitFile :
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

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(float64(sizeInBytes), float64(fileSize-int64(i*uint64(sizeInBytes)))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := "l4t.0" + string(i)
		_, err := os.Create(outpath + "/" + fileName)

		if err != nil {
			return err
		}

		// write/save buffer to disk
		ioutil.WriteFile(outpath+"/"+fileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to : ", outpath+"/"+fileName)
	}
	return nil
}

// DirSizeMB :
func DirSizeMB(path string) float64 {
	var dirSize int64 = 0

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}

		return nil
	}

	filepath.Walk(path, readSize)

	sizeMB := float64(dirSize) / 1024.0 / 1024.0

	return sizeMB
}

// CopyDirectory :
func CopyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy :
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

// Exists :
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateIfNotExists :
func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// CopySymLink :
func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
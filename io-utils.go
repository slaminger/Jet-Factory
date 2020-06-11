package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Mirantis/virtlet/pkg/diskimage/guestfs"
)

// CreateDisk :
func CreateDisk(src, dst, disk, format string) (*guestfs.GuestfsError, error) {
	size := DirSizeMB(src)
	rootSize := int64(size) + 256
	fmt.Println("Estimated size:", rootSize, "MB")

	// if err := ExecWrapper("dd", "of="+dst+"/"+disk+".img", "bs=1", "count=0", "seek="+rootSize+"M"); err != nil {
	// 	return nil, err
	// }

	g, errno := guestfs.Create()
	if errno != nil {
		return nil, errno
	}
	defer g.Close()

	f, ferr := os.Create(dst + "/" + disk + ".img")
	if ferr != nil {
		return nil, ferr
	}
	defer f.Close()

	ferr = f.Truncate(rootSize * 1024 * 1024)
	if ferr != nil {
		return nil, ferr
	}

	/* Attach the disk image to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        true,
	}
	err := g.Add_drive(disk, &optargs)
	if err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	err = g.Launch()
	if err != nil {
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
		fmt.Println("\nexpected a single device from list-devices")
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
		fmt.Println("\nexpected a single partition from list-partitions")
		return err, nil
	}

	/* Create a filesystem on the partition. */
	err = g.Mkfs(format, partitions[0], nil)
	if err != nil {
		return err, nil
	}

	if err = g.Shutdown(); err != nil {
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
	defer g.Close()

	/* Attach the disk image read-only to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        true,
	}
	err := g.Add_drive(disk, &optargs)
	if err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	err = g.Launch()
	if err != nil {
		return err, nil
	}

	/* Ask libguestfs to inspect for operating systems. */
	roots, err := g.Inspect_os()
	if err != nil {
		return err, nil
	}

	var root string
	var ferr error
	if len(roots) > 1 {
		root, ferr = CliSelector("Select root partition to use:", roots)
		if ferr != nil {
			return nil, ferr
		}
	} else if len(roots) == 1 {
		root = roots[0]
	} else {
		return nil, errors.New("inspect-vm: no operating systems found for")
	}

	err = g.Mount(root, "/")
	if err != nil {
		return err, nil
	}

	err = g.Mount_local("/mnt", nil)
	if err != nil {
		return err, nil
	}

	go g.Mount_local_run()

	files, ok := filepath.Glob("/mnt/*")
	if ok != nil {
		return nil, ok
	}

	for _, dir := range files {
		if err := CopyDirectory(dir, dst); err != nil {
			log.Println(err)
		}
	}

	err = g.Shutdown()
	if err != nil {
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
	defer g.Close()

	/* Attach the disk image read-only to libguestfs. */
	optargs := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}
	err := g.Add_drive(disk, &optargs)
	if err != nil {
		return err, nil
	}

	/* Run the libguestfs back-end. */
	err = g.Launch()
	if err != nil {
		return err, nil
	}

	err = g.Mount(disk, "/")
	if err != nil {
		return err, nil
	}

	err = g.Mount_local("/mnt", nil)
	if err != nil {
		return err, nil
	}

	go g.Mount_local_run()

	errno = CopyDirectory(src, "/mnt/")
	if err != nil {
		return nil, errno
	}

	err = g.Shutdown()
	if err != nil {
		return err, nil
	}

	return nil, nil
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

	fmt.Printf("\nSplitting to %d pieces.\n", totalPartsNum)

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
		fmt.Println("\nSplit to : ", outpath+"/"+fileName)
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
			return nil
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("Failed to get raw syscall.Stat_t data for '%s'", sourcePath)
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
			CopySymLink(sourcePath, destPath)
		default:
			Copy(sourcePath, destPath)
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			os.Chmod(destPath, entry.Mode())
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
		return fmt.Errorf("Failed to create directory: '%s', error: '%s'", dir, err.Error())
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

// RetryFunction :
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

// ExecWrapper :
func ExecWrapper(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

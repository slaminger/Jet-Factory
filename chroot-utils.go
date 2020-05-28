package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/pkg/mount"
	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

// PreChroot : Copy qemu-aarch64-static binary and mount bind the directories
func PreChroot(mountpoint string) error {
	err := Copy("/usr/bin/qemu-"+buildarch+"-static", mountpoint+"/usr/bin/qemu-"+buildarch+"-static")
	if err != nil {
		return err
	}
	if mounted, _ := mount.Mounted(mountpoint); !mounted {
		if err := mount.Mount(mountpoint, mountpoint, "bind", "rbind,rw"); err != nil {
			return err
		}
	}
	return nil
}

// PostChroot : Remove qemu-aarch64-static binary and unmount the binded directories
func PostChroot(mountpoint string) error {
	err := ExecWrapper("rm", mountpoint+"/usr/bin/qemu-"+buildarch+"-static")
	if err != nil {
		return err
	}

	if mounted, _ := mount.Mounted(mountpoint); mounted {
		if err := mount.Unmount(mountpoint); err != nil {
			return err
		}
	}
	return nil
}

// Chroot :
func Chroot(path string) error {
	if err := cg(); err != nil {
		return err
	}

	caps, err := capability.NewPid(0)
	if err != nil {
		return err
	}

	// if the process doesn't have CAP_SYS_ADMIN, but does have CAP_SYS_CHROOT, we need to use the actual chroot
	if !caps.Get(capability.EFFECTIVE, capability.CAP_SYS_ADMIN) && caps.Get(capability.EFFECTIVE, capability.CAP_SYS_CHROOT) {
		return realChroot(path)
	}

	if err := unix.Unshare(unix.CLONE_NEWNS); err != nil {
		return fmt.Errorf("Error creating mount namespace before pivot: %v", err)
	}

	// make everything in new ns private
	if err := mount.MakeRPrivate("/"); err != nil {
		return err
	}

	if mounted, _ := mount.Mounted(path); !mounted {
		if err := mount.Mount(path, path, "bind", "rbind,rw"); err != nil {
			return realChroot(path)
		}
	}

	// setup oldRoot for pivot_root
	pivotDir, err := ioutil.TempDir(path, ".pivot_root")
	if err != nil {
		return fmt.Errorf("Error setting up pivot dir: %v", err)
	}

	var mounted bool
	defer func() {
		if mounted {
			// make sure pivotDir is not mounted before we try to remove it
			if errCleanup := unix.Unmount(pivotDir, unix.MNT_DETACH); errCleanup != nil {
				if err == nil {
					err = errCleanup
				}
				return
			}
		}

		errCleanup := os.Remove(pivotDir)
		// pivotDir doesn't exist if pivot_root failed and chroot+chdir was successful
		// because we already cleaned it up on failed pivot_root
		if errCleanup != nil && !os.IsNotExist(errCleanup) {
			errCleanup = fmt.Errorf("Error cleaning up after pivot: %v", errCleanup)
			if err == nil {
				err = errCleanup
			}
		}
	}()

	if err := unix.PivotRoot(path, pivotDir); err != nil {
		// If pivot fails, fall back to the normal chroot after cleaning up temp dir
		if err := os.Remove(pivotDir); err != nil {
			return fmt.Errorf("Error cleaning up after failed pivot: %v", err)
		}
		return realChroot(path)
	}
	mounted = true

	// This is the new path for where the old root (prior to the pivot) has been moved to
	// This dir contains the rootfs of the caller, which we need to remove so it is not visible during extraction
	pivotDir = filepath.Join("/", filepath.Base(pivotDir))

	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("Error changing to new root: %v", err)
	}

	// Make the pivotDir (where the old root lives) private so it can be unmounted without propagating to the host
	if err := unix.Mount("", pivotDir, "", unix.MS_PRIVATE|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("Error making old root private after pivot: %v", err)
	}

	// Now unmount the old root so it's no longer visible from the new root
	if err := unix.Unmount(pivotDir, unix.MNT_DETACH); err != nil {
		return fmt.Errorf("Error while unmounting old root after pivot: %v", err)
	}
	mounted = false

	return nil
}

func realChroot(path string) error {
	if err := cg(); err != nil {
		return err
	}

	if err := unix.Chroot(path); err != nil {
		return fmt.Errorf("Error after fallback to chroot: %v", err)
	}
	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("Error changing to new root after chroot: %v", err)
	}
	return nil
}

func cg() error {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "liz"), 0755)
	err := ioutil.WriteFile(filepath.Join(pids, "liz/pids.max"), []byte("20"), 0700)
	if err != nil {
		return err
	}
	// Removes the new cgroup in place after the container exits
	err = ioutil.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	if err != nil {
		return err
	}
	return nil
}

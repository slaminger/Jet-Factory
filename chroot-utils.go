package main

import (
	"os"
)

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

	return nil
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeclysm/extract"
	"github.com/xi2/xz"
)

// ExtractFiles :
func ExtractFiles(archivePath, dst string) (err error) {
	fmt.Println("\nExtracting:", archivePath, "in:", dst)
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
	} else if strings.Contains(archivePath, ".tar.gz") || strings.Contains(archivePath, ".zip") || strings.Contains(archivePath, ".rar") {
		data, err := ioutil.ReadFile(archivePath)
		if err != nil {
			return err
		}
		buffer := bytes.NewBuffer(data)
		switch filepath.Ext(archivePath) {
		case ".bz2":
			err = extract.Bz2(context.Background(), buffer, dst, nil)
		case ".gz":
			err = extract.Gz(context.Background(), buffer, dst, nil)
		case ".zip":
			err = extract.Zip(context.Background(), buffer, dst, nil)
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

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

/* Net utilities
* WalkURL using regex
* Wget to download a file
 */

// WalkURL : Walk a URL, and return the body
func WalkURL(source string) *string {
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

	body := string(bodyBytes)
	return &body
}

// PrintDownloadPercent :
func PrintDownloadPercent(done chan int64, path string, total int64) {
	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}

			fi, err := file.Stat()
			if err != nil {
				fmt.Println(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

// DownloadFile :
func DownloadFile(url string, dest string) (err error) {

	file := path.Base(url)

	fmt.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	start := time.Now()

	out, err := os.Create(path.String())

	if err != nil {
		fmt.Println(path.String())
		return err
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		return err
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		return err
	}

	done := make(chan int64)

	go PrintDownloadPercent(done, path.String(), int64(size))

	var (
		response *http.Response
		retries  int = 5
	)
	for retries > 0 {
		response, err = http.Get(url)
		if err != nil {
			log.Println(err)
			retries--
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
		n, err := io.Copy(out, response.Body)

		if err != nil {
			return err
		}

		done <- n

		elapsed := time.Since(start)
		fmt.Printf("Download completed in %s", elapsed)
	}

	return nil
}

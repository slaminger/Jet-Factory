package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

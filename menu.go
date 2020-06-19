package main

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

/* Menu - CLI
* Selector function
 */

// CliSelect : Select an item in a menu from cli
func CliSelect(label string, items []string) (answer string, err error) {
	var inputValue string
	prompt := &survey.Select{
		Message: label,
		Options: items,
	}
	if err := survey.AskOne(prompt, &inputValue); err != nil {
		return "", err
	}
	return inputValue, nil
}

// CliInput : Set an input in a menu from cli
func CliInput(label string) (answer string, err error) {
	var inputValue string
	prompt := &survey.Input{
		Message: label,
	}
	if err := survey.AskOne(prompt, &inputValue); err != nil {
		return "", err
	}
	return inputValue, nil
}

// SelectDistro :
func SelectDistro() (*string, error) {
	var avalaibles []string
	for _, baseDistro := range basesDistro {
		for _, variantDistro := range baseDistro.Variants {
			avalaibles = append(avalaibles, variantDistro.Name)
		}
		avalaibles = append(avalaibles, baseDistro.Name)
	}

	name, err := CliSelect("", avalaibles)
	if err != nil {
		return nil, err
	}

	return &name, nil
}

// SelectArchitecture :
func SelectArchitecture() error {
	var avalaible []string
	for archi := range distribution.Architectures {
		avalaible = append(avalaible, archi)
	}
	arch, err := CliSelect("Choose an avalaible Architecture :", avalaible)
	if err != nil {
		return err
	}
	buildarch = arch
	return nil
}

// SelectVersion : Retrieve a URL for a distribution based on a version
func SelectVersion() (constructedURL string, err error) {
	for _, avalaibleMirror := range distribution.Architectures[buildarch] {
		// If the string contains the tag {VERSION} then try to replace the tag by walking the URL

		if strings.Contains(avalaibleMirror, "{VERSION}") {

			constructedURL = strings.Split(avalaibleMirror, "/{VERSION}")[0]
			versionBody := WalkURL(constructedURL)

			// TODO : Rework this
			search, _ := regexp.Compile(">:?([[:digit:]]{1,3}.[[:digit:]]+|[[:digit:]]+)(?:/)")
			match := search.FindAllStringSubmatch(*versionBody, -1)

			if match == nil {
				return "", errors.New("Couldn't find any match for regex")
			}
			// TODO END

			versions := make([]string, 0)
			for i := 0; i < len(match); i++ {
				for _, submatches := range match {
					versions = append(versions, submatches[1])
				}
			}

			version, err := CliSelect("Select a version: ", versions)
			if err != nil {
				return "", err
			}
			constructedURL = strings.Replace(avalaibleMirror, "{VERSION}", version, 1)
			imageBody := WalkURL(constructedURL)

			log.Println("ImageBody:", *imageBody)
			search, _ = regexp.Compile(">:?([[:alpha:]]+.*.raw.xz)")
			imageMatch := search.FindAllStringSubmatch(*imageBody, -1)
			images := make([]string, 0)

			log.Println("ImageMatch:", imageMatch)
			for i := 0; i < len(imageMatch); i++ {
				for _, submatches := range imageMatch {
					images = append(images, submatches[1])
				}
			}

			var imageFile string
			if len(images) > 1 {
				imageFile, err = CliSelect("Select an image file: ", images)
				if err != nil {
					return "", err
				}
			} else if len(images) == 1 {
				imageFile = images[0]
			} else {
				return "", err
			}

			log.Println("ImageFile:", imageFile)

			return strings.TrimSpace(constructedURL + imageFile), nil

			// TODO : Rework this following ugly stuff
		} else if strings.Contains(avalaibleMirror, ".raw.") || strings.Contains(avalaibleMirror, ".tar.") || strings.Contains(avalaibleMirror, ".tbz2") || strings.Contains(avalaibleMirror, ".zip") || strings.Contains(avalaibleMirror, ".rar") || strings.Contains(avalaibleMirror, ".gz") {
			return avalaibleMirror, nil
		} else {
			constructedURL, err := CliInput("No URL found in config, input one that point directly to your rootfs :")
			if err != nil {
				return "", err
			}
			return constructedURL, nil
		}
	}
	return "", errors.New("Unknown issue occured")
}

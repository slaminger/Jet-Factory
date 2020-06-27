package main

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// CliSelect : Select an item in a menu from cli; returns nil and the answer on success; returns err otherwise;
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

// CliInput : Set an input in a menu from cli; returns nil and the input on success; returns err otherwise;
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

// SelectDistro : Select an avalaible distribution to build; returns nil and the selected distro name on success; returns err otherwise;
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

// SelectArchitecture : Select an avalaible build architecture; returns nil on success; returns err otherwise;
func SelectArchitecture() error {
	var avalaible []string

	if isVariant {
		for archi := range variant.Architectures {
			avalaible = append(avalaible, archi)
		}
	} else {
		for archi := range base.Architectures {
			avalaible = append(avalaible, archi)
		}
	}

	arch, err := CliSelect("Choose an avalaible Architecture :", avalaible)
	if err != nil {
		return err
	}
	buildarch = arch
	return nil
}

// SelectVersion : Retrieve a URL for a distribution based on a version; returns nil and the constructed URL on success; returns err otherwise;
func SelectVersion() (constructedURL string, err error) {
	var avalaibleMirrors []string
	var selectedMirror string

	if isVariant {
		for _, avalaibleMirror := range variant.Architectures[buildarch] {
			avalaibleMirrors = append(avalaibleMirrors, avalaibleMirror)
		}
	} else {
		for _, avalaibleMirror := range base.Architectures[buildarch] {
			avalaibleMirrors = append(avalaibleMirrors, avalaibleMirror)
		}
	}

	if len(avalaibleMirrors) > 1 {
		selectedMirror, err = CliSelect("Choose a URL :", avalaibleMirrors)
		if err != nil {
			return "", err
		}
	} else if len(avalaibleMirrors) == 1 {
		selectedMirror = avalaibleMirrors[0]
	}

	// TODO : Rework this following ugly stuff
	if strings.Contains(selectedMirror, ".raw.") || strings.Contains(selectedMirror, ".tar.") || strings.Contains(selectedMirror, ".tbz2") || strings.Contains(selectedMirror, ".zip") || strings.Contains(selectedMirror, ".rar") || strings.Contains(selectedMirror, ".gz") {
		return selectedMirror, nil
	}

	constructedURL, err = CliInput("No URL found in config, input one that point directly to your rootfs :")
	if err != nil {
		return "", err
	}

	return constructedURL, nil
}

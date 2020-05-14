package main

import "github.com/AlecAivazis/survey/v2"

/* Menu - CLI
* Selector function
 */

// CliSelector : Select an item in a menu froim cli
func CliSelector(label string, items []string) (answer string, err error) {
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

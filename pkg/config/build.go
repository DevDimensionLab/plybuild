package config

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

var validateNumber = func(input string) error {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}

func BuildConfigInteractive(config *ProjectConfiguration) (err error) {
	config.GroupId, err = InteractivePrompt("GroupId", config.GroupId, nil)
	config.ArtifactId, err = InteractivePrompt("ArtifactId", config.ArtifactId, nil)
	config.Package, err = InteractivePrompt("PackageName", config.Package, nil)
	config.Language, err = InteractivePrompt("Language", config.Language, nil)
	config.Name, err = InteractivePrompt("Name", config.Name, nil)
	config.Description, err = InteractivePrompt("Description", config.Description, nil)
	return
}

func InteractivePrompt(label, defaultValue string, validator func(input string) error) (value string, err error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: validator,
	}
	return prompt.Run()
}

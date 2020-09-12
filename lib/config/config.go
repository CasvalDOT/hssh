// Package config provides ...
package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// Config structure include some configurations
// params about hssh config. look below
type Config struct {
	GitlabBaseURL      string
	GitlabPrivateToken string
	GitlabProjectID    string
	GitlabFiles        []string
}

// Parse Return a tuple contain key value found
// in the configuration file
func Parse(content string) (string, string) {
	parsed := strings.Split(content, "=")
	return parsed[0], parsed[1]
}

// Read the specific configuration file at given path
func Read(path string) (Config, error) {
	config := Config{}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value := Parse(scanner.Text())

		switch key {
		case "GITLAB_BASE_URL":
			config.GitlabBaseURL = value
			break
		case "GITLAB_PRIVATE_TOKEN":
			config.GitlabPrivateToken = value
			break
		case "GITLAB_PROJECT_ID":
			config.GitlabProjectID = value
			break
		case "GITLAB_FILES":
			config.GitlabFiles = strings.Split(value, ",")
			break
		}
	}

	return config, nil

}

// Get the highter file configuration available
func Get() (Config, error) {
	homeDir, _ := os.UserHomeDir()

	config := Config{}
	var allowedPathConfigurations [2]string

	allowedPathConfigurations[0] = homeDir + "/.config/hssh/config"
	allowedPathConfigurations[1] = "/etc/hssh/config"

	// Start read conf files, from highter to lower
	for _, path := range allowedPathConfigurations {
		config, err := Read(path)
		if err != nil {
			continue
		}

		return config, nil
	}

	return config, errors.New("Missing file configuration")

}

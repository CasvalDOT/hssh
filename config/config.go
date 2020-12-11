package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v2"
)

var allowedPaths = []string{
	"/etc/hssh",
	"{{HOME}}/.config/hssh",
}

// Config ..
type Config interface {
	fuzzyBinaryExist() bool
	read(string) error

	Load() error
	Create(string) error

	GetProvider() providerConfig
	GetDefaultProvider() string
	GetFuzzyEngine() string
}

type providerConfig struct {
	Host         string   `yaml:"host"`
	PrivateToken string   `yaml:"private_token"`
	ProjectID    string   `yaml:"project_id"`
	Files        []string `yaml:"files"`
}

type config struct {
	Provider        providerConfig `yaml:"provider"`
	FuzzyEngine     string         `yaml:"fuzzyengine"`
	DefaultProvider string         `yaml:"default_provider"`
}

func replaceHomePlaceholder(path string) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	regex := regexp.MustCompile("{{HOME}}")
	return regex.ReplaceAllString(path, userHomeDir), nil
}

func (c *config) read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	if c.FuzzyEngine == "" {
		c.FuzzyEngine = "fzf"
	}

	if c.fuzzyBinaryExist() == false {
		c.FuzzyEngine = ""
	}

	return nil
}

func (c *config) fuzzyBinaryExist() bool {
	cmdOutput := &bytes.Buffer{}
	cmd := exec.Command("which", c.FuzzyEngine)
	cmd.Stdout = cmdOutput
	cmd.Stderr = nil
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

// Create ...
func (c *config) Create(content string) error {
	for _, path := range allowedPaths {
		fileName := "config.yml"
		pathFolder, err := replaceHomePlaceholder(path)
		if err != nil {
			continue
		}

		pathWithFile := pathFolder + "/" + fileName

		_, err = os.Stat(pathWithFile)
		if err == nil {
			fmt.Println("FILE EXIST")
			pathWithFile = pathWithFile + ".example"
		}

		fmt.Println(pathWithFile)
		file, err := os.Create(pathWithFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer file.Close()
		file.WriteString(content)

	}
	return nil
}

// Load ...
func (c *config) Load() error {
	fileReads := 0
	for _, path := range allowedPaths {
		pathWithFile := path + "/config.yml"
		pathWithFile, err := replaceHomePlaceholder(pathWithFile)
		if err != nil {
			continue
		}

		_, err = os.Stat(pathWithFile)
		if os.IsNotExist(err) {
			continue
		}

		c.read(pathWithFile)
		if err != nil {
			continue
		}
		fileReads = fileReads + 1
	}

	if fileReads == 0 {
		return errors.New("NO_VALID_CONF_FILES")
	}

	return nil
}

func (c *config) GetProvider() providerConfig {
	return c.Provider
}

func (c *config) GetDefaultProvider() string {
	return c.DefaultProvider
}

func (c *config) GetFuzzyEngine() string {
	return c.FuzzyEngine
}

// New ...
func New() Config {
	// Read files here of configuration
	return &config{}
}

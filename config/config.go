package config

/*
Config entity is use to interact with the file
configuration and the environemnt of the cli.
*/
import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v2"
)

var resolverBinary = "which"
var configFileName = "config.yml"
var defaultFuzzysearch = "fzf"
var homePlaceholder = "{{HOME}}"
var defaultFolderName = "hssh"

var allowedPaths = []string{
	"/etc/" + defaultFolderName,
	homePlaceholder + "/.config/" + defaultFolderName,
}

// IConfig ..
type IConfig interface {
	fuzzyBinaryExist() bool
	read(string) error

	Load() error
	Create(string) error
	GetProvider() providerConfig
	GetDefaultProvider() string
	GetFuzzysearch() string
}

type providerConfig struct {
	Host         string `yaml:"host"`
	Path         string `yaml:"path"`
	PrivateToken string `yaml:"private_token"`
	ProjectID    string `yaml:"project_id"`
}

type config struct {
	Provider        providerConfig `yaml:"provider"`
	Fuzzysearch     string         `yaml:"fuzzysearch"`
	DefaultProvider string         `yaml:"default_provider"`
}

// replaceHomePlaceholder
/*............................................................................*/
func replaceHomePlaceholder(path string) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	regex := regexp.MustCompile(homePlaceholder)
	return regex.ReplaceAllString(path, userHomeDir), nil
}

// read
/*............................................................................*/
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

	if c.Fuzzysearch == "" {
		c.Fuzzysearch = defaultFuzzysearch
	}

	/*
		Fuzzy search engine is not mandatory
		for hssh. So instead generate and error
		if binary path cannot be found, we "unset"
		fuzzysearch
	*/
	if c.fuzzyBinaryExist() == false {
		c.Fuzzysearch = ""
	}

	return nil
}

// fuzzyBinaryExist
/*............................................................................*/
func (c *config) fuzzyBinaryExist() bool {
	cmdOutput := &bytes.Buffer{}
	cmd := exec.Command(resolverBinary, c.Fuzzysearch)
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
/*............................................................................*/
func (c *config) Create(content string) error {
	for _, path := range allowedPaths {
		fileName := configFileName
		pathFolder, err := replaceHomePlaceholder(path)
		if err != nil {
			continue
		}

		pathWithFile := pathFolder + "/" + fileName

		_, err = os.Stat(pathWithFile)
		if err == nil {
			fmt.Println("File", pathWithFile, "just exist")
			pathWithFile = pathWithFile + ".example"
		}

		file, err := os.Create(pathWithFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("File", pathWithFile, "created")

		defer file.Close()
		file.WriteString(content)
	}
	return nil
}

// Load ...
/*............................................................................*/
func (c *config) Load() error {
	fileReads := 0
	for _, path := range allowedPaths {
		pathWithFile := path + "/" + configFileName
		pathWithFile, err := replaceHomePlaceholder(pathWithFile)
		if err != nil {
			continue
		}

		_, err = os.Stat(pathWithFile)
		if err != nil || os.IsNotExist(err) {
			continue
		}

		c.read(pathWithFile)

		fileReads = fileReads + 1
	}

	if fileReads == 0 {
		return errors.New("No valid configuration file")
	}

	return nil
}

// GetProvider
/*............................................................................*/
func (c *config) GetProvider() providerConfig {
	return c.Provider
}

// GetDefaultProvider
/*............................................................................*/
func (c *config) GetDefaultProvider() string {
	return c.DefaultProvider
}

// GetFuzzysearch
/*............................................................................*/
func (c *config) GetFuzzysearch() string {
	return c.Fuzzysearch
}

// New ...
/*............................................................................*/
func New() IConfig {
	return &config{}
}

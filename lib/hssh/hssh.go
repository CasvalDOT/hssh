// Package hssh provide a simple way for serach and connect ssh servers
package hssh

import (
	"bytes"
	"errors"
	"fmt"
	"hssh/lib/gitlab"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

// ProviderConfig structure
/*
	Describe the param for connect and fetch configuration files
	contains the list of the SSH servers
*/
type ProviderConfig struct {
	Host         string   `yaml:"host"`
	PrivateToken string   `yaml:"private_token"`
	ProjectID    string   `yaml:"project_id"`
	Files        []string `yaml:"files"`
}

// Config structure
// Describe the list of providers stored in the configurationo
type Config struct {
	FuzzysearchBinary string         `yaml:"fuzzysearch"`
	Gitlab            ProviderConfig `yaml:"gitlab"`
}

// HSSH structure
// Rapresent the hssh instance. It contain the configuration file attributes
type HSSH struct {
	Configuration Config
}

// ReplaceObject structure
// Rapresent a search and replace structure
type ReplaceObject struct {
	value   string
	replace string
}

/*
---------------------------------------------------------------------------
Private functions
----------------------------------------------------------------------------
*/

/*
	_readConfig function

	Read che configuration file at path provided
*/
func _readConfig(path string) (Config, error) {

	var cfg Config

	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

/*
	_getSSHDirectory function

	Return the absolute path of the ssh configuration folder download
*/
func _getSSHDirectory(folder string) string {
	homeDir, _ := os.UserHomeDir()
	folderPath := homeDir + "/.ssh/" + folder

	return folderPath
}

/*
	_createSSHListFile

	Create the file contains the list of ssh connections configuration
	inside the .ssh/<folder> directory
*/
func _createSSHListFile(folder string, fileName string, content []byte) {

	// Get folder path
	folderPath := _getSSHDirectory(folder)

	// Create configuration directory
	os.MkdirAll(folderPath, os.ModePerm)

	// Create file
	file, err := os.Create(folderPath + "/" + fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	if _, err := file.Write(content); err != nil {
		panic(err)
	}
	if err := file.Sync(); err != nil {
		panic(err)
	}
}

/*
	_createTempFile function

	Create a temp file to store
	the formatted ssh files list configuration
	after parsing it
*/
func createTempFile(content string) *os.File {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "hssh-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return tmpFile
}

/*
	_resolveFuzzysearchBinary function

	Search executable path of fuzzysearch engine
*/
func _resolveFuzzysearchBinary(engine string) bool {
	cmdOutput := &bytes.Buffer{}

	cmd := exec.Command("which", engine)
	cmd.Stdout = cmdOutput
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

/*
	_search function

	Search for connections in a wall of text using fuzzysearch utility
	with pattern provided using the standard input
*/
func (hssh *HSSH) _search(connections string, fuzzysearch bool) string {
	// Store connections parsed in a temp file
	tmpFile := createTempFile(connections)
	defer os.Remove(tmpFile.Name())

	/*
		Create a command to execute
		for retrieve connections
	*/
	command := "cat " + tmpFile.Name()
	if fuzzysearch {

		// Check for binary of fuzzysearch engine in use
		isBinResolved := _resolveFuzzysearchBinary(hssh.Configuration.FuzzysearchBinary)

		if isBinResolved == true {
			command = command + " | " + hssh.Configuration.FuzzysearchBinary
		}

	}

	// Initialize a variable to store output of the command
	cmdOutput := &bytes.Buffer{}

	// Excute command
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = cmdOutput
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	/*
		Close the temporaly file
		previously created
	*/
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Return output
	return string(string(cmdOutput.Bytes()))
}

/*
	_parseConnections function

	Normalize the ssh connections file in a provided format
	useful for a simply and rapic search and overview
*/
func (hssh *HSSH) _parseConnections(format string) string {

	files := hssh.Configuration.Gitlab.Files

	var connections string = ""

	/*
		Read all config files and chain it
		NOTE. If an error occured during read
		the file will be skipped
	*/
	for i := 0; i < len(files); i++ {
		re := regexp.MustCompile(`(%2F|\/)`)
		path := re.ReplaceAllString(files[i], `,`)
		fileFolder := strings.Split(path, ",")

		folder := fileFolder[0]
		file := fileFolder[1]

		folderPath := _getSSHDirectory(folder)

		content, err := ioutil.ReadFile(folderPath + "/" + file)
		if err != nil {
			continue
		}
		connections = connections + string(content)
	}

	regex := [7]ReplaceObject{
		{`(?mi)^(#.*)$`, ""},
		{`(?mi)\n`, ""},
		{`(?mi)Host `, "\nHost "},
		{`(Hostname|Host|User|Port)`, ""},
		{` +`, " "},
		{`(?m)^ `, ""},
		{`(?m)^(.*) (.*) (.*) (.*)$`, format},
	}

	for _, reg := range regex {
		re := regexp.MustCompile(reg.value)
		connections = re.ReplaceAllString(connections, reg.replace)
	}

	return connections

}

/*
---------------------------------------------------------------------------
Public functions
----------------------------------------------------------------------------
*/

// Sync function
/*
	Fetch connections files from provided services
	used and save it in the ~/.ssh folder
*/
func (hssh *HSSH) Sync() {

	// Take project ID
	var projectID string = hssh.Configuration.Gitlab.ProjectID

	// Define a gitlab instance using authentication envs
	var gitlab = gitlab.Gitlab{
		BaseURL:      hssh.Configuration.Gitlab.Host,
		PrivateToken: hssh.Configuration.Gitlab.PrivateToken,
	}

	// Define a list of files to take from gitlab
	var files = hssh.Configuration.Gitlab.Files

	var wg sync.WaitGroup

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		// Go routine start to fetch files using
		// HTTP Requests
		go func(url string) {

			defer wg.Done()

			fileDecoded, err := gitlab.GetFile(projectID, url)

			// Get folder path
			re := regexp.MustCompile(`(\/|%2F).*`)
			folder := re.ReplaceAllString(url, ``)

			if err != nil {
				fmt.Println(err)
			} else {
				_createSSHListFile(folder, fileDecoded.Name, fileDecoded.Content)
			}

		}(files[i])

	}

	wg.Wait()
}

// List function
/*
	Aggregate the list of ssh connections stored in the
	downloaded files and print a list. You can
	even search inside the list if you provide the fuzzysearch flag
*/
func (hssh *HSSH) List(fuzzysearch bool) string {
	// Obtain connections
	connections := hssh._parseConnections(`$1 -> $3@$2:$4`)

	// Perform a search
	output := hssh._search(connections, fuzzysearch)

	// Return output
	return output
}

// Exec function
// Perfotm an SSH connection after selected from the flatten list
func (hssh *HSSH) Exec() {
	// Obtain connections
	connections := hssh._parseConnections(`$1 -> ssh $3@$2 -p $4`)

	// Perform a search (fuzzyseacrh is set to true)
	command := hssh._search(connections, true)

	// Remove unused string part to obtain a valid ssh command
	re := regexp.MustCompile(`^.*?->\s`)
	command = re.ReplaceAllString(command, "")

	// Execute SSH command returned
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// LoadConfig function
/*
	Read and load the configuration found in the filesystem.
	The configurations are hierarchical. In order:
	- /etc/hssh/config.yml
	- ~/.config/hssh/config.yml
*/
func (hssh *HSSH) LoadConfig() (Config, error) {

	homeDir, _ := os.UserHomeDir()
	var cfg Config
	var isConfigLoad bool = false

	var allowedPathConfigurations [2]string

	allowedPathConfigurations[0] = homeDir + "/.config/hssh/config.yml"
	allowedPathConfigurations[1] = "/etc/hssh/config.yml"

	for _, path := range allowedPathConfigurations {

		conf, err := _readConfig(path)

		if err != nil {
			continue
		}

		cfg = conf

		isConfigLoad = true
	}

	if isConfigLoad == false {
		return cfg, errors.New("Error loading config")
	}

	hssh.Configuration = cfg

	return cfg, nil

}

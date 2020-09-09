// Package hssh provides some functions for a rapid search on ssh connections
package hssh

import (
	"bytes"
	"fmt"
	"hssh/lib/config"
	"hssh/lib/gitlab"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// HSSH rapresent the structure of the HSSH content
type HSSH struct {
	Configuration config.Config
}

func getSSHDirectory(folder string) string {
	// Configuration ssh path
	homeDir, _ := os.UserHomeDir()
	folderPath := homeDir + "/.ssh/" + folder

	return folderPath
}

func createConfigurationFile(folder string, fileName string, content []byte) {

	// Get folder path
	folderPath := getSSHDirectory(folder)

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
Create a temp file to store
the formatted ssh files configuration
after parsing it
*/
func createTempFile(content string) *os.File {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "hssh-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Example writing to the file
	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return tmpFile
}

/*
	Search in a list of connections provided
*/
func search(connections string, fuzzysearch bool) string {
	// Store connections parsed in a temp file
	tmpFile := createTempFile(connections)
	defer os.Remove(tmpFile.Name())

	// Create a command to execute
	// for retrieve connections
	command := "cat " + tmpFile.Name()
	if fuzzysearch {
		command = command + " | fzf"
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

	// Close the temporaly file
	// previously created
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Return output
	return string(string(cmdOutput.Bytes()))
}

/*
Parse the SSH connections in a more human readable
format, putting it in one line for perform rapid search
*/
func (hssh *HSSH) parseConnections(format string) string {

	files := hssh.Configuration.GitlabFiles

	var connections string = ""

	// Read all cofig files and chain it
	// NOTE. If an error occured during read
	// the file will be skipped
	for i := 0; i < len(files); i++ {
		re := regexp.MustCompile(`^.*?%2F`)
		file := re.ReplaceAllString(files[i], ``)

		// Get folder path
		re = regexp.MustCompile(`(\/|%2F).*`)
		folder := re.ReplaceAllString(files[i], ``)

		folderPath := getSSHDirectory(folder)

		content, err := ioutil.ReadFile(folderPath + "/" + file)
		if err != nil {
			continue
		}
		connections = connections + string(content)
	}

	// Remove comments
	var re = regexp.MustCompile(`(?mi)^(#.*)$`)
	connections = re.ReplaceAllString(connections, ``)

	// Remove \n
	connections = strings.ReplaceAll(connections, "\n", "")

	// Replace Host with \nHost.
	// this allowed to have a configuration in one line
	connections = strings.ReplaceAll(connections, "Host ", "\nHost ")

	// Remove attributes names
	re = regexp.MustCompile(`(Hostname|Host|User|Port)`)
	connections = re.ReplaceAllString(connections, ``)

	// Remove empty spaces
	re = regexp.MustCompile(` +`)
	connections = re.ReplaceAllString(connections, ` `)

	re = regexp.MustCompile(`(?m)^ `)
	connections = re.ReplaceAllString(connections, ``)

	// Apply format provided
	re = regexp.MustCompile(`(?m)^(.*) (.*) (.*) (.*)$`)
	connections = re.ReplaceAllString(connections, format)

	return connections

}

/*
Sync new SSH configurations files
from Gitlab repository.
Then save it the the .ssh user home folder
*/
func (hssh *HSSH) Sync() {

	// Take project ID
	var projectID string = hssh.Configuration.GitlabProjectID

	// Define a gitlab instance using authentication envs
	var gitlab = gitlab.Gitlab{
		BaseURL:      hssh.Configuration.GitlabBaseURL,
		PrivateToken: hssh.Configuration.GitlabPrivateToken,
	}

	// Define a list of files to take from gitlab
	// NOTE
	// In the configurations files are comma separated
	var files = hssh.Configuration.GitlabFiles

	// TODO
	// Use go routines
	for i := 0; i < len(files); i++ {
		fileDecoded, err := gitlab.GetFile(projectID, files[i])

		// Get folder path
		re := regexp.MustCompile(`(\/|%2F).*`)
		folder := re.ReplaceAllString(files[i], ``)

		if err != nil {
			fmt.Println(err)
			continue
		}

		createConfigurationFile(folder, fileDecoded.Name, fileDecoded.Content)
	}
}

/*
List the SSH connections configured
in the format provided. List can be searchable
using fuzzysearch, if -f is provided
*/
func (hssh *HSSH) List(fuzzysearch bool) string {
	// Obtain connections
	connections := hssh.parseConnections(`$1 -> $3@$2:$4`)

	// Perform a search
	output := search(connections, fuzzysearch)

	// Return output
	return output
}

/*
Exec an SSH connection after seraching inside the list
*/
func (hssh *HSSH) Exec() {
	// Obtain connections
	connections := hssh.parseConnections(`$1 -> ssh $3@$2 -p $4`)

	// Perform a search
	// NOTE that fuzzysearch arguments is set to true
	command := search(connections, true)

	// Remove unused string part
	// to obtain a valid ssh command
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

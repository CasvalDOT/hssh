// Package cli ...
/*
The cli package contains a bunch of methods
mapped 1:1 with the flags provided by
*/
package cli

import (
	"bytes"
	"fmt"
	"hssh/providers"
	"hssh/sshconfig"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// ICli ...
type ICli interface {
	List() (string, error)
	Sync(string, string)
	Exec() error
	Print(string)

	search(string) (string, error)
	getConnections(string) (string, error)
}

type cli struct {
	fuzzysearch   string
	filesToSearch []string
	provider      providers.IProvider
	colors        bool
}

func (c *cli) getConnections(format string) (string, error) {
	connections, err := sshconfig.Format(format)
	if err != nil {
		return "", err
	}

	return c.search(connections)
}

func (c *cli) search(connections string) (string, error) {
	tmpFile, err := sshconfig.Temporize(connections)
	if err != nil {
		return "", err
	}
	command := "cat " + tmpFile.Name()

	if c.fuzzysearch != "" {
		command = command + " | " + c.fuzzysearch
	}

	cmdOutput := &bytes.Buffer{}

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = cmdOutput
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	if err := tmpFile.Close(); err != nil {
		return "", err
	}

	return string(string(cmdOutput.Bytes())), nil
}

// List ...
func (c *cli) List() (string, error) {
	return c.getConnections(`{{.Name}} -> {{.User}}@{{.Hostname}}:{{.Port}}`)
}

// Exec ...
func (c *cli) Exec() error {
	command, err := c.getConnections(`{{.Name}} -> ssh {{.User}}@{{.Hostname}} -p {{.Port}}`)

	// Remove unused string part to obtain a valid ssh command
	re := regexp.MustCompile(`^.*?->\s`)
	command = re.ReplaceAllString(command, "")

	// Execute SSH command returned
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Sync ...
func (c *cli) Sync(projectID string, path string) {
	var wg sync.WaitGroup

	filesOfTheProject, err := c.provider.GetFiles(projectID, path)
	if err != nil {
		return
	}

	for _, fileFromProvider := range filesOfTheProject {
		wg.Add(1)

		go func(fileID string, filePath string) {
			defer wg.Done()

			content, err := c.provider.GetFile(projectID, fileID)
			if err != nil {
				return
			}

			fmt.Println(filePath)
			splits := strings.Split(filePath, "/")
			folder := splits[0]
			fileName := splits[1]

			sshconfig.Create(folder, fileName, content)

		}(fileFromProvider.ID, fileFromProvider.Path)
	}

	wg.Wait()
}

// Print
func (c *cli) Print(content string) {

	if c.colors {
		re := regexp.MustCompile(`(?m)(.*)\s->\s(ssh\s|)(.*?)@(.*?)(:|\s-p\s)(.*)$`)
		content = re.ReplaceAllString(content, "\033[36m$1\033[0m -> \033[32m$3\033[0m@\033[33m$4\033[0m$5\033[31m$6\033[0m")
	}

	fmt.Printf(content)
}

// New ...
func New(fuzzysearch string, p providers.IProvider, colors bool) ICli {
	return &cli{
		fuzzysearch: fuzzysearch,
		provider:    p,
		colors:      colors,
	}
}

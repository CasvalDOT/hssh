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
	"sync"
)

// ICli ...
type ICli interface {
	List() (string, error)
	Sync(string)
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
func (c *cli) Sync(projectID string) {
	var wg sync.WaitGroup

	for _, fileFromProvider := range c.filesToSearch {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			fileDecoded, err := c.provider.GetFile(projectID, file)
			if err != nil {
				return
			}

			re := regexp.MustCompile(`(\/|%2F).*`)
			folder := re.ReplaceAllString(file, ``)

			sshconfig.Create(folder, fileDecoded.Name, fileDecoded.Content)

		}(fileFromProvider)
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
func New(fuzzysearch string, p providers.IProvider, filesToSearch []string, colors bool) ICli {
	return &cli{
		fuzzysearch:   fuzzysearch,
		provider:      p,
		filesToSearch: filesToSearch,
		colors:        colors,
	}
}

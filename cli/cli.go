// Package cli ...
/*
The cli package contains a bunch of methods
mapped 1:1 with the flags provided by
*/
package cli

import (
	"bytes"
	"fmt"
	"hssh/connections"
	"hssh/providers"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var tempFileName = "hssh.tmp"
var printRegexp = `(?m)(.*)\s->\s(ssh\s|)(.*?)@(.*?)(:|\s-p\s)(.*)$`
var defaultConnectionFormat = "[{{.ID}}] {{.Name}} -> {{.User}}@{{.Hostname}}:{{.Port}}"

// ICli ...
type ICli interface {
	List() (string, error)
	Sync(string, string)
	Exec() error
	Print(string)

	search(string) (string, error)
	toTempFile(string) (*os.File, error)
}

type cli struct {
	fuzzysearch string
	sshUA       connections.ISSHUA
	provider    providers.IProvider
	colors      bool
}

func (c *cli) toTempFile(format string) (*os.File, error) {
	connections := c.sshUA.List()
	var context = ""
	for _, connection := range connections {
		formattedConnection, err := connection.ToString(format)
		if err != nil {
			continue
		}
		context += formattedConnection + "\n"
	}

	return createTempFile(context)
}

func (c *cli) getIDFromSelection(selection string) (int, error) {
	rgx := regexp.MustCompile("(?mi).*\\[(.*?)\\].*\n")
	idString := rgx.ReplaceAllString(selection, "$1")
	return strconv.Atoi(idString)
}

func (c *cli) search(format string) (string, error) {
	tmpFile, err := c.toTempFile(format)
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
	connectionFormat := defaultConnectionFormat
	return c.search(connectionFormat)
}

// Exec ...
func (c *cli) Exec() error {
	results, err := c.List()
	if err != nil {
		return err
	}

	id, err := c.getIDFromSelection(results)
	if err != nil {
		return err
	}

	connection := c.sshUA.SearchConnectionByID(id)

	connection.Connect()

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

			c.sshUA.Create(folder, fileName, content)

		}(fileFromProvider.ID, fileFromProvider.Path)
	}

	wg.Wait()
}

// Print
func (c *cli) Print(content string) {
	if c.colors {
		re := regexp.MustCompile(printRegexp)
		content = re.ReplaceAllString(content, "\033[36m$1\033[0m -> \033[32m$3\033[0m@\033[33m$4\033[0m$5\033[31m$6\033[0m")
	}

	fmt.Printf(content)
}

// New ...
func New(fuzzysearch string, p providers.IProvider, sshUA connections.ISSHUA, colors bool) ICli {
	return &cli{
		fuzzysearch: fuzzysearch,
		provider:    p,
		sshUA:       sshUA,
		colors:      colors,
	}
}

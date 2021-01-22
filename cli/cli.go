// Package cli ...
/*
This package is a wrapper around the entire
hssh logic. Instead of dirt the main file with lot of logics
i've created this one.
*/
package cli

import (
	"bytes"
	"fmt"
	"hssh/providers"
	"hssh/sshuseragent"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var tempFileName = "hssh.tmp"
var printRegexp = `(?m)(.*)\s->\s(ssh\s|)(.*?)@(.*?)(:|\s-p\s)(.*)$`
var defaultConnectionFormat = "[{{.ID}}] {{.Name}} -> {{.User}}@{{.Hostname}}:{{.Port}}"

// ICli ...
type ICli interface {
	List() (string, error)
	Sync(string)
	Connect() error
	Print(string)

	list(string) (string, error)
}

type cli struct {
	fuzzysearch string
	sshUA       sshuseragent.IsshUserAgent
	provider    providers.IProvider
	colors      bool
}

/*
	getListOfConnections
	...........................................................
	obtain a list of connections from different sources.
	First check for a cached file, otherwise read each config file
	in SSH folder
*/
func (c *cli) getListOfConnections(format string) string {
	return getConnectionsFromFiles(c, format)
}

/*
	list
	.........................................................
	Show the list of connections available
*/
func (c *cli) list(format string) (string, error) {
	context := c.getListOfConnections(format)

	command := cat(context)
	command = pipeFuzzysearch(command, c.fuzzysearch)

	cmdOutput := &bytes.Buffer{}

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = cmdOutput
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil && err.Error() != "exit status 130" {
		return context, err
	}

	context = string(cmdOutput.Bytes())

	return context, nil
}

/*
	List
	...............................................................
	Return the full list of connections in the format provided.
*/
func (c *cli) List() (string, error) {
	return c.list(defaultConnectionFormat)
}

/*
	Connect
	................................................................
	Allow to select a connections from the list using fuzzysearch.
	Once selected an ssh command start
*/
func (c *cli) Connect() error {
	results, err := c.List()
	if err != nil {
		return err
	}

	id, err := getIDFromSelection(results)
	if err != nil {
		return err
	}

	return c.sshUA.Connect(id)
}

/*
	Sync
	................................................................
	Download and save the files with ssh configurations
	from the provided declared in configuration file
*/
func (c *cli) Sync(providerConnectionString string) {
	var wg sync.WaitGroup

	projectID, path, err := getProjectIDAndPath(providerConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesOfTheProject, err := c.provider.GetFiles(projectID, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fileFromProvider := range filesOfTheProject {
		wg.Add(1)

		go func(fileID string, filePath string) {
			defer wg.Done()

			content, err := c.provider.GetFile(projectID, fileID)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(filePath, "[OK]")
			splits := strings.Split(filePath, "/")
			folder := splits[0]
			fileName := splits[1]

			c.sshUA.CreateConfig(folder, fileName, content)

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
func New(fuzzysearch string, p providers.IProvider, sshUA sshuseragent.IsshUserAgent, colors bool) ICli {
	return &cli{
		fuzzysearch: fuzzysearch,
		provider:    p,
		sshUA:       sshUA,
		colors:      colors,
	}
}

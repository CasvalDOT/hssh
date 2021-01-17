package cli

import (
	"errors"
	"hssh/cache"
	"regexp"
	"strconv"
	"strings"
)

/*............................................................................*/
func getFromCache(c *cli, format string) string {
	context, _ := cache.Read(tempFileName)
	return context
}

/*............................................................................*/
func getFromFiles(c *cli, format string) string {
	context := ""
	connections := c.sshUA.List()
	for _, connection := range connections {
		formattedConnection, err := connection.ToString(format)
		if err != nil {
			continue
		}
		context += formattedConnection + "\n"
	}

	return context
}

/*............................................................................*/
func cat(context string) string {
	return "echo -e '" + context + "'"
}

/*............................................................................*/
func pipeFuzzysearch(command string, fuzzyBinary string) string {
	if fuzzyBinary == "" {
		return command
	}
	command += " | " + fuzzyBinary
	return command
}

/*............................................................................*/
func getIDFromSelection(selection string) (int, error) {
	rgx := regexp.MustCompile("(?mi).*\\[(.*?)\\].*\n")
	idString := rgx.ReplaceAllString(selection, "$1")
	return strconv.Atoi(strings.Trim(idString, " "))
}

/*............................................................................*/
func getProjectIDAndPath(providerConnectionString string) (string, string, error) {
	rgx := regexp.MustCompile("^.*:/(.*)@(.*)$")
	matches := rgx.FindAllStringSubmatch(providerConnectionString, 1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", "", errors.New("Cannot find project ID or Path in the provided string")
	}

	return matches[0][1], matches[0][2], nil
}

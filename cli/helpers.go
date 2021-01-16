package cli

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

/*............................................................................*/
func createTempFile(content string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFileName)
	if err != nil {
		return nil, err
	}

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

/*............................................................................*/
func cat(tmpFile *os.File) string {
	return "cat " + tmpFile.Name()
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
	return strconv.Atoi(idString)
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

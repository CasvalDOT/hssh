package cli

import (
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

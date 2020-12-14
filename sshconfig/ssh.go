// Package sshconfig ...
/*
This package provide a small set of functions to interact with ssh files
configuration.
*/
package sshconfig

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var pairRegex = " (.*?)(\\s|#|$)"
var nameRegex = "(.*?) Hostname"
var sshConnectionAttributes = [][]string{
	{"Name", "-"},
	{"Hostname", "-"},
	{"IdentityFile", ""},
	{"User", "root"},
	{"Port", "22"},
}
var configRegexs = [5]replaceObject{
	{`(?mi)#.*\s`, ""},
	{`(?mi)Host `, "#Name "},
	{`(?mi)\s\s+`, `\n`},
	{`(?mi)\\n`, ""},
	{`(Name|Hostname|IdentityFile|User|Port)`, ` $1`},
}

type replaceObject struct {
	value   string
	replace string
}

func homeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func makePath(folderName string) string {
	absolutePath := homeDir() + "/.ssh/" + folderName
	return absolutePath
}

func extractData(key string, context string) string {
	regexToCompile := regexp.MustCompile(key + pairRegex)
	matchs := regexToCompile.FindAllStringSubmatch(context, -1)
	if len(matchs) < 1 {
		return ""
	}

	return matchs[0][1]
}

func valueOrFallback(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func listConnections(rawConnections string) []string {
	for _, reg := range configRegexs {
		re := regexp.MustCompile(reg.value)
		rawConnections = re.ReplaceAllString(rawConnections, reg.replace)
	}

	return strings.Split(rawConnections, "#")
}

func fromRawToFormattedConnection(rawConnection string, format string) string {
	tmpl, err := template.New("T").Parse(format)
	if err != nil {
		return ""
	}

	var templateBuffer bytes.Buffer
	var templateData = map[string]interface{}{}

	for _, key := range sshConnectionAttributes {
		value := valueOrFallback(
			extractData(key[0], rawConnection),
			key[1],
		)
		templateData[key[0]] = value
	}

	tmpl.Execute(&templateBuffer, templateData)

	return string(templateBuffer.Bytes())
}

// ReadAll ...
func ReadAll(files []string) string {
	var content string = ""
	for _, file := range files {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		content = content + string(f)
	}
	return content
}

// Create ...
func Create(folderName string, fileName string, content []byte) (string, error) {
	folderPath := makePath(folderName)
	filePath := folderPath + "/" + fileName

	os.MkdirAll(folderPath, os.ModePerm)

	newFile, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}

	defer newFile.Close()

	if _, err := newFile.Write(content); err != nil {
		return filePath, err
	}
	if err := newFile.Sync(); err != nil {
		return filePath, err
	}

	return filePath, nil
}

// Temporize ...
func Temporize(content string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "hssh-")
	if err != nil {
		return nil, err
	}

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

// Search ...
/*
Search in .ssh homedir folder and take the files
that contains the "config" string in the path;
for example:
 config.test.d/servers
 myconnections/config
*/
func Search() ([]string, error) {
	var files []string

	sshAbsolutePath := makePath("")
	filepath.Walk(sshAbsolutePath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if strings.Contains(path, "config") == false {
			return nil
		}

		if info.IsDir() == false {
			files = append(files, path)
		}

		return nil
	})

	return files, nil
}

// Format ...
func Format(format string) (string, error) {
	parsedConnections := ""

	configFiles, err := Search()
	if err != nil {
		return parsedConnections, err
	}

	connections := ReadAll(configFiles)

	listOfRawConnections := listConnections(connections)

	for _, rawConnection := range listOfRawConnections {
		formattedConnection := fromRawToFormattedConnection(rawConnection, format)
		parsedConnections = parsedConnections + formattedConnection + "\n"
	}

	return parsedConnections, nil
}

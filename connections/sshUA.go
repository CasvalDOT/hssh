package connections

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nameKey = "Name"
var portKey = "Port"
var identityFileKey = "IdentityFile"
var hostnameKey = "Hostname"
var userKey = "User"
var defaultConnectionName = "-"
var defaultConnectionHostname = "-"
var defaultConnectionIdentityFile = ""
var defaultConnectionPort = "22"
var defaultConnectionUser = "root"
var filePatternMatch = "config"
var templateName = "format_connection"
var pairRegex = " (.*?)(\\s|#|$)"
var nameRegex = "(.*?) " + hostnameKey
var connectionSeparator = "#"

type replaceObject struct {
	value   string
	replace string
}

var configRegexs = [5]replaceObject{
	{`(?mi)#.*\s`, ""},
	{`(?mi)Host `, "#" + nameKey + " "},
	{`(?mi)\s\s+`, `\n`},
	{`(?mi)\\n`, ""},
	{`(` + nameKey + `|` + hostnameKey + `|` + identityFileKey + `|` + userKey + `|` + portKey + `)`, ` $1`},
}

// ISSHUA ...
type ISSHUA interface {
	Load()
	List() []IConnection
	SearchConnectionByID(int) IConnection
	CreateConfig(string, string, []byte) (string, error)
}

type sshUA struct {
	connections []IConnection
	configFiles []string
}

// normalizeContextConfig
/*............................................................................*/
func (ssh *sshUA) normalizeContextConfig(context string) []string {
	for _, reg := range configRegexs {
		re := regexp.MustCompile(reg.value)
		context = re.ReplaceAllString(context, reg.replace)
	}

	return strings.Split(context, connectionSeparator)
}

// readConnectionsFromConfig
/*............................................................................*/
func (ssh *sshUA) readConnectionsFromConfig(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	connectionsString := ssh.normalizeContextConfig(string(file))

	for _, connectionString := range connectionsString {
		name := extractData(nameKey, connectionString, defaultConnectionName)
		hostname := extractData(hostnameKey, connectionString, defaultConnectionHostname)
		user := extractData(userKey, connectionString, defaultConnectionUser)
		port := extractData(portKey, connectionString, defaultConnectionPort)
		identityFile := extractData(identityFileKey, connectionString, defaultConnectionIdentityFile)

		if name == defaultConnectionName {
			continue
		}

		ssh.connections = append(ssh.connections, NewConnection(len(ssh.connections), name, hostname, user, identityFile, port))
	}

	return nil
}

// loadConnections
/*............................................................................*/
func (ssh *sshUA) loadConnections() {
	for _, configFile := range ssh.configFiles {
		err := ssh.readConnectionsFromConfig(configFile)
		if err != nil {
			continue
		}
	}

}

// loadConfigs
/*............................................................................*/
func (ssh *sshUA) loadConfigs() {
	var configFiles []string
	sshAbsolutePath := makePath("")
	filepath.Walk(sshAbsolutePath, func(path string, info os.FileInfo, err error) error {
		if err != nil ||
			strings.Contains(path, filePatternMatch) == false ||
			info.IsDir() == true {
			return nil
		}

		configFiles = append(configFiles, path)

		return nil
	})

	ssh.configFiles = configFiles
}

func (ssh *sshUA) Load() {
	ssh.loadConfigs()
	ssh.loadConnections()
}

func (ssh *sshUA) List() []IConnection {
	ssh.Load()
	return ssh.connections
}

// SearchConnectionByID
/*............................................................................*/
func (ssh *sshUA) SearchConnectionByID(id int) IConnection {
	for _, connection := range ssh.connections {
		if connection.GetID() == id {
			return connection
		}
	}

	return nil
}

// CreateConfig
/*............................................................................*/
func (ssh *sshUA) CreateConfig(folderName string, fileName string, content []byte) (string, error) {
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

// NewSSHUA ...
/*............................................................................*/
func NewSSHUA() ISSHUA {
	return &sshUA{}
}

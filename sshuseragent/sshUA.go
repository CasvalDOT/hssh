package sshuseragent

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SSH Connection attributes
var nameKey = "Name"
var portKey = "Port"
var identityFileKey = "IdentityFile"
var hostnameKey = "Hostname"
var userKey = "User"

// SSH Default values for connection attributes
var defaultConnectionName = "-"
var defaultConnectionHostname = "-"
var defaultConnectionIdentityFile = ""
var defaultConnectionPort = "22"
var defaultConnectionUser = "root"

// We use this string to identify a valid configuration
// file in .ssh folder.
var filePatternMatch = "config"

// Other data ...
var templateName = "format_connection"
var pairRegex = " (.*?)(\\s|#|$)"
var nameRegex = "(.*?) " + hostnameKey
var connectionSeparator = "#"

// TODO: shell binary can be a parameter in config
var shellBinary = "bash"

/*
 Replace object rapresent two
 strings: one is the part to replace, the
 second is the content to apply. For example:
 having a string "hello world"
 if we need to replace hello with "hi",
 the form of the replace object mus be:

 replaceObject{"hello", "hi"}
*/
type replaceObject struct {
	value   string
	replace string
}

/*
For well reading the entire list
of connections params, we apply this list
of regex of our content.
*/
var configRegexs = [5]replaceObject{
	{`(?mi)#.*\s`, ""},                  // Remove all comments
	{`(?mi)Host `, "#" + nameKey + " "}, // Replace "Host " with #Name char
	{`(?mi)\s\s+`, `\n`},                // Remove multiple spaces
	{`(?mi)\\n`, ""},                    // Remove new lines
	{`(` + nameKey + `|` + hostnameKey + `|` + identityFileKey + `|` + userKey + `|` + portKey + `)`, ` $1`}, // Replace <key> with " <key>. Is useful for next logic and easy argument taking"
}

// IsshUserAgent ...
type IsshUserAgent interface {
	Load()
	Connect(int) error
	List() []IConnection
	SearchConnectionByID(int) IConnection
	CreateConfig(string, string, []byte) (string, error)
}

type sshUserAgent struct {
	connections []IConnection
	configFiles []string
}

// normalizeContextConfig
/*
............................................................................
We must normalize the list of connections
apply the regex defined at the start of this file
*/
func (ssh *sshUserAgent) normalizeContextConfig(context string) []string {
	for _, reg := range configRegexs {
		re := regexp.MustCompile(reg.value)
		context = re.ReplaceAllString(context, reg.replace)
	}

	return strings.Split(context, connectionSeparator)
}

// readConnectionsFromConfig
/*............................................................................*/
func (ssh *sshUserAgent) readConnectionsFromConfig(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	connectionsString := ssh.normalizeContextConfig(string(file))

	/*
		Once normalized extract data , for each connection found,
		create a new connection structure to append to a list of structured connections.
	*/
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
func (ssh *sshUserAgent) loadConnections() {
	for _, configFile := range ssh.configFiles {
		err := ssh.readConnectionsFromConfig(configFile)
		if err != nil {
			continue
		}
	}

}

// loadConfigs
/*............................................................................*/
func (ssh *sshUserAgent) loadConfigs() {
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

// Load
/*............................................................................*/
func (ssh *sshUserAgent) Load() {
	ssh.loadConfigs()
	ssh.loadConnections()
}

// List
/*............................................................................*/
func (ssh *sshUserAgent) List() []IConnection {
	ssh.Load()
	return ssh.connections
}

// SearchConnectionByID
/*............................................................................*/
func (ssh *sshUserAgent) SearchConnectionByID(id int) IConnection {
	for _, connection := range ssh.connections {
		if connection.GetID() == id {
			return connection
		}
	}

	return nil
}

// Connect
/*............................................................................*/
func (ssh *sshUserAgent) Connect(connectionID int) error {
	connection := ssh.SearchConnectionByID(connectionID)
	cmd := exec.Command(shellBinary, "-c", connection.GetSSHConnection())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// CreateConfig
/*............................................................................*/
func (ssh *sshUserAgent) CreateConfig(folderName string, fileName string, content []byte) (string, error) {
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

// NewSSHUserAgent ...
/*............................................................................*/
func NewSSHUserAgent() IsshUserAgent {
	ua := sshUserAgent{}
	return &ua
}

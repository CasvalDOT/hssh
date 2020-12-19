package connections

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"
)

// IConnection ...
type IConnection interface {
	ToString(string) (string, error)
	Connect() error
	GetID() int
}

type connection struct {
	ID           int
	Name         string
	User         string
	Port         string
	IdentityFile string
	Hostname     string
}

// GetID ...
/*............................................................................*/
func (c *connection) GetID() int {
	return c.ID
}

// GetName ...
/*............................................................................*/
func (c *connection) GetName() string {
	return c.Name
}

// GetUser ...
/*............................................................................*/
func (c *connection) GetUser() string {
	return c.User
}

// GetPort ...
/*............................................................................*/
func (c *connection) GetPort() string {
	return c.Port
}

// GetIdentityFile ...
/*............................................................................*/
func (c *connection) GetIdentityFile() string {
	return c.IdentityFile
}

// GetHostname ...
/*............................................................................*/
func (c *connection) GetHostname() string {
	return c.Hostname
}

// ToString ...
/*............................................................................*/
func (c *connection) ToString(format string) (string, error) {
	tmpl, err := template.New(templateName).Parse(format)
	var connectionString string
	if err != nil {
		return connectionString, err
	}

	var templateBuffer bytes.Buffer
	tmpl.Execute(&templateBuffer, c)

	connectionString = string(templateBuffer.Bytes())

	return connectionString, nil
}

// Connect ...
/*............................................................................*/
func (c *connection) Connect() error {
	var argsOfTheCommand = "ssh "
	argsOfTheCommand += c.User + "@" + c.Hostname + " -p " + c.Port

	cmd := exec.Command("bash", "-c", argsOfTheCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// NewConnection ...
/*............................................................................*/
func NewConnection(
	id int,
	name string,
	hostname string,
	user string,
	identityFile string,
	port string,
) IConnection {
	return &connection{
		ID:           id,
		Name:         name,
		User:         user,
		Hostname:     hostname,
		IdentityFile: identityFile,
		Port:         port,
	}
}

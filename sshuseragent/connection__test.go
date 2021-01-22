package sshuseragent

import (
	"testing"
)

// TestToString ...
func TestToString(t *testing.T) {
	c := NewConnection(1, "Test", "10.0.30.3", "root", "identity", "22")
	connectionToString, err := c.ToString("{{.ID}} - {{.Name}} - {{.Port}} - {{.Hostname}} - {{.User}} - {{.IdentityFile}}")

	if err != nil {
		t.Errorf("Should not return any error")
	}

	if connectionToString != "1 - Test - 22 - 10.0.30.3 - root - identity" {
		t.Errorf("Should return the expected string, not %s", connectionToString)
	}
}

// TestGetSSHConnection
func TestGetSSHConnection(t *testing.T) {
	c := NewConnection(1, "Test", "10.0.30.3", "root", "", "22")
	connectionString := c.GetSSHConnection()

	if connectionString != "ssh root@10.0.30.3 -p 22" {
		t.Errorf("Should return a valid connection, not %s", connectionString)
	}
}

// TestGetSSHConnectionWithIdentity
func TestGetSSHConnectionWithIdentity(t *testing.T) {
	c := NewConnection(1, "Test", "10.0.30.3", "root", "~/.ssh/id_rsa", "22")
	connectionString := c.GetSSHConnection()

	if connectionString != "ssh root@10.0.30.3 -p 22 -i ~/.ssh/id_rsa" {
		t.Errorf("Should return a valid connection, not %s", connectionString)
	}
}

package test

import (
	"hssh/cli"
	"testing"
)

// TestGitlabStart ...
/*
	Providing a valid connection string, token and other params
	must be set correctly
*/
func TestGetProjectIDAndPath(t *testing.T) {
	projectID, path, err := cli.TgetProjectIDAndPath("gitlab://test:/123456@main/folder")
	if err != nil {
		t.Log("Should not return any errors")
		t.Fail()
	}

	if projectID != "123456" {
		t.Log("Project ID should be 123456 instead of", projectID)
		t.Fail()
	}

	if path != "main/folder" {
		t.Log("Path should be main/folder instead of", path)
		t.Fail()
	}
}

func TestRealConnectionString(t *testing.T) {
	projectID, path, err := cli.TgetProjectIDAndPath("github://1234567890:/CasvalDOT/hssh@providers")
	if err != nil {
		t.Log("Should not return any errors")
		t.Fail()
	}

	if projectID != "CasvalDOT/hssh" {
		t.Log("Project ID should be CasvalDOT/hssh instead of", projectID)
		t.Fail()
	}

	if path != "providers" {
		t.Log("Path should be providers instead of", path)
		t.Fail()
	}
}

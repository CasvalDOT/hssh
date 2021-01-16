package test

import (
	"hssh/providers"
	"testing"
)

// TestConnectionString ...
/*
	Providing a valid connection string, token and other params
	must be set correctly
*/
func TestConnectionString(t *testing.T) {
	g := providers.NewGitlab("gitlab://token")
	token := g.GetPrivateToken()
	if token != "token" {
		t.Log("should return 'token' but instead returns", token)
		t.Fail()
	}
}

// TestConnectionStringWithMultipleColons ...
/*
	Providing a valid connection string, token and other params
	must be set correctly
*/
func TestConnectionStringWithMultipleColons(t *testing.T) {
	g := providers.NewGitlab("gitlab://t:ok:en:/123/folder")
	token := g.GetPrivateToken()
	if token != "t:ok:en" {
		t.Log("should return 't:ok:en' but instead returns", token)
		t.Fail()
	}
}

// TestConnectionStringWithAdditionalData ...
/*
	Providing a valid connection string, token and other params
	must be set correctly
*/
func TestConnectionStringWithAdditionalData(t *testing.T) {
	g := providers.NewGitlab("gitlab://token:/123/folder")
	token := g.GetPrivateToken()
	if token != "token" {
		t.Log("should return 'token' but instead returns", token)
		t.Fail()
	}
}

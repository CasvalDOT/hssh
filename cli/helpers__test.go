package cli

import (
	"testing"
)

// TestGetIDFromSelection ...
func TestGetIDFromSelection(t *testing.T) {
	index, e := getIDFromSelection("[1] hosting -> root@10.10.10.1:22\n")
	if e != nil {
		t.Errorf("Should not return any error")
	}

	if index != 1 {
		t.Errorf("Should return 1 instead of %d", index)
	}

	index, e = getIDFromSelection("[ 1 ] hosting -> root@10.10.10.1:22\n")
	if e != nil {
		t.Errorf("Should not return any error")
	}

	if index != 1 {
		t.Errorf("Should return 1 even with spaces instead of %d", index)
	}
}

// TestGetProjectIDAndPath ...
func TestGetProjectIDAndPath(t *testing.T) {
	projectID, path, err := getProjectIDAndPath("gitlab://test:123456:/CasvalDOT/hssh@providers")
	if err != nil {
		t.Errorf("Should not return any error")
	}

	if projectID != "CasvalDOT/hssh" {
		t.Errorf("Should return the projectID provided, not %s", projectID)
	}

	if path != "providers" {
		t.Errorf("Should return the path provided, not %s", path)
	}
}

// TestPipeFuzzySearch ...
func TestPipeFuzzySearch(t *testing.T) {
	command := pipeFuzzysearch("echo test", "fzf")
	if command != "echo test | fzf" {
		t.Errorf("Should return the command with fzf concatenated, not %s", command)
	}
}

// TestCat ..
func TestCat(t *testing.T) {
	command := cat("alfa,beta,gamma")
	if command != "echo -e 'alfa,beta,gamma'" {
		t.Errorf("Should return the echo of the content provided, not %s", command)
	}
}

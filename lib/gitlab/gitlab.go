// Package gitlab provides a basic function for interact with Gitlab
package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Gitlab structure
// Include some configurations about Gitlab API
type Gitlab struct {
	PrivateToken string
	BaseURL      string
}

// File strucrture
// Describe the structure of a Gitlab file repository
type File struct {
	Content string `json:"content"`
	Name    string `json:"file_name"`
}

// FileDecoded structure
// Descrive the Gitlab file converted in bytes
type FileDecoded struct {
	Content []byte
	Name    string
}

/*
---------------------------------------------------------------------------
Public functions
----------------------------------------------------------------------------
*/

// Get function
// Perform a request using the provided Gitlab API
func (gitlab *Gitlab) Get(endpoint string) ([]byte, error) {

	req, err := http.NewRequest("GET", gitlab.BaseURL+endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlab.PrivateToken)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetFile function
// Return the content about a file in the Gitlab repository
func (gitlab *Gitlab) GetFile(projectID string, filePath string) (FileDecoded, error) {
	res, err := gitlab.Get("/projects/" + projectID + "/repository/files/" + filePath + "?ref=master")

	var fileDecoded = FileDecoded{}

	if err != nil {
		return fileDecoded, err
	}

	file := File{}

	jsonErr := json.Unmarshal(res, &file)

	if jsonErr != nil {
		return fileDecoded, jsonErr
	}

	// Decode the file content from base 64 to array of bytes
	sd, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return fileDecoded, err
	}

	fileDecoded.Content = sd
	fileDecoded.Name = file.Name

	return fileDecoded, nil

}

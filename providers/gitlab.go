package providers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type gitlab struct {
	provider
}

func (g *gitlab) get(endpoint string) ([]byte, error) {
	request, err := http.NewRequest("GET", g.url+endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("PRIVATE-TOKEN", g.privateToken)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (g *gitlab) GetFile(projectID string, filePath string) (*fileDecoded, error) {
	endpoint := "/projects/" + projectID + "/repository/files/" + filePath + "?ref=master"
	bodyInBytes, err := g.get(endpoint)
	if err != nil {
		return nil, err
	}

	f := file{}
	err = json.Unmarshal(bodyInBytes, &f)
	if err != nil {
		return nil, err
	}

	contentBase64, err := base64.StdEncoding.DecodeString(f.Content)
	if err != nil {
		return nil, err
	}

	return &fileDecoded{
		Content: contentBase64,
		Name:    f.Name,
	}, nil
}

// NewGitlab ...
func NewGitlab(url string, privateToken string) IProvider {
	return &gitlab{
		provider: provider{
			url:          url,
			privateToken: privateToken,
		},
	}
}

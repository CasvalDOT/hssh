package providers

// Provider ...
type Provider interface {
	iGet
	iGetFile
}

// iGet
type iGet interface {
	get(string) ([]byte, error)
}

type iGetFile interface {
	GetFile(string, string) (*fileDecoded, error)
}

type file struct {
	Content string `json:"content"`
	Name    string `json:"file_name"`
}

type fileDecoded struct {
	Content []byte
	Name    string
}

type provider struct {
	url          string
	privateToken string
}

// New ...
func New(driver string, url string, privateToken string) Provider {
	if driver == "gitlab" {
		return &gitlab{
			provider: provider{
				url:          url,
				privateToken: privateToken,
			},
		}
	}
	return nil
}

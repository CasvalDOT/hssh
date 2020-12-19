package providers

/*
Provider is a abstract class that decscribe
the concrete classes used to fetch the connections
files from a remote repository

NOTE: Now is currently supported gitlab
*/

// IProvider ...
type IProvider interface {
	iGet
	iGetFile
	iGetFiles
}

type iGet interface {
	get(string, []queryParam) ([]byte, error)
}

type iGetFiles interface {
	GetFiles(string, string) ([]file, error)
}

type iGetFile interface {
	GetFile(string, string) ([]byte, error)
}

type file struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Name    string `json:"file_name"`
	Path    string `json:"path"`
}

type queryParam struct {
	key   string
	value string
}

/*
Provider use two attributes
url and privateToken.

url is the repo link where files can be found.

privateToken instead permit to
authenticate to the service
*/
type provider struct {
	url          string
	privateToken string
}

// New ...
func New(driver string, url string, privateToken string) IProvider {
	if driver == "gitlab" {
		return NewGitlab(url, privateToken)
	}
	panic("INVALID PROVIDER")
}

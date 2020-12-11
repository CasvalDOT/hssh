package worker

import (
	"hssh/config"
	"hssh/helpers"
	"hssh/providers"
	"regexp"
	"sync"
)

// Worker ...
type Worker interface {
	Exec()
	List()
	Search()
	Sync()
	Assets()
}

type worker struct {
	provider providers.Provider
	config   config.Config
}

// Exec ...
func (w *worker) Exec() {

}

// List ...
func (w *worker) List() {

}

// Search ...
func (w *worker) Search() {

}

// Sync ...
func (w *worker) Sync() {

	provider := w.config.GetProvider()
	var wg sync.WaitGroup

	for _, fileFromProvider := range provider.Files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			fileDecoded, err := w.provider.GetFile(provider.ProjectID, file)
			if err != nil {
				return
			}

			// Get folder path
			re := regexp.MustCompile(`(\/|%2F).*`)
			folder := re.ReplaceAllString(file, ``)

			helpers.CreateSSHConfig(folder, fileDecoded.Name, fileDecoded.Content)

		}(fileFromProvider)
	}

}

// Assets ...
func (w *worker) Assets() {

}

// New ...
func New() Worker {
	return &worker{}
}

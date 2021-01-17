package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// get file in cache
func get(fileName string) string {
	filePath := ""
	filepath.Walk(os.TempDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil ||
			strings.Contains(path, fileName) == false {
			return nil
		}

		filePath = path
		return nil
	})

	return filePath
}

// Clear cache
func Clear(pattern string) {
	filepath.Walk(os.TempDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil ||
			strings.Contains(path, pattern) == false {
			return nil
		}

		os.Remove(path)
		return nil
	})
}

// Read cache
func Read(fileName string) (string, error) {
	fileToRead := get(fileName)
	content, err := ioutil.ReadFile(fileToRead)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Write cache
func Write(fileName string, content string) (*os.File, error) {
	Clear(fileName)

	tmpFile, err := ioutil.TempFile(os.TempDir(), fileName)
	if err != nil {
		return nil, err
	}

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

package cli

import (
	"io/ioutil"
	"os"
)

func createTempFile(content string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), tempFileName)
	if err != nil {
		return nil, err
	}

	text := []byte(content)
	if _, err = tmpFile.Write(text); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

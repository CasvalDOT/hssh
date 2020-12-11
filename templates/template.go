package templates

import (
	"bytes"
	"text/template"
)

// Print ...
func Print(content string, iface map[string]interface{}) error {
	tmpl, err := template.New("T").Parse(content)
	if err != nil {
		return err
	}
	var templateBuffer bytes.Buffer
	tmpl.Execute(&templateBuffer, iface)

	return nil
}

// SaveInFile ...
func SaveInFile(filePath string, content string) {
}

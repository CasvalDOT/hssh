package connections

import (
	"os"
	"regexp"
)

func homeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func makePath(folderName string) string {
	absolutePath := homeDir() + "/.ssh/" + folderName
	return absolutePath
}

func extractData(key string, context string, defaultValue string) string {
	regexToCompile := regexp.MustCompile(key + pairRegex)
	matchs := regexToCompile.FindAllStringSubmatch(context, -1)
	if len(matchs) < 1 {
		return defaultValue
	}
	return matchs[0][1]
}

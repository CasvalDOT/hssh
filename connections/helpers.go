package connections

import (
	"os"
	"regexp"
)

// homeDir
/*............................................................................*/
func homeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

// makePath
/*............................................................................*/
func makePath(folderName string) string {
	absolutePath := homeDir() + "/.ssh/" + folderName
	return absolutePath
}

// extractData
/*............................................................................*/
func extractData(key string, context string, defaultValue string) string {
	regexToCompile := regexp.MustCompile(key + pairRegex)
	matchs := regexToCompile.FindAllStringSubmatch(context, -1)
	if len(matchs) < 1 {
		return defaultValue
	}
	return matchs[0][1]
}

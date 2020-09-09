// Package hssh provides a way for sync heply servers
package main

import (
	"flag"
	"fmt"
	"hssh/lib/config"
	"hssh/lib/hssh"
	"log"
	"os"
	"regexp"
)

func out(content string, withColors bool) {

	if withColors {
		re := regexp.MustCompile(`(?m)(.*)\s->\s(ssh\s|)(.*?)@(.*?)(:|\s-p\s)(.*)$`)
		content = re.ReplaceAllString(content, "\033[36m$1\033[0m -> \033[32m$3\033[0m@\033[33m$4\033[0m$5\033[31m$6\033[0m")
	}

	fmt.Printf(content)
}

func help() {
	fmt.Println("")
	fmt.Println("HSSH - An heply utility to connect into the server's company")
	fmt.Println("")
	fmt.Println("OPTIONS")
	fmt.Println("-l return a list of connections available")
	fmt.Println("-s sync the files in the repository in your local machine")
	fmt.Println("-f enable fuzzy search (it work only in conjuction with -l flag)")
	fmt.Println("-le search inside the list in fuzzy search mode and perform an ssh connection to selected host")
	fmt.Println("")
}

// TODO
// Try a better approach for flags
func main() {
	// Define flags
	isFuzzy := flag.Bool("f", false, "Enable fuzzysearch using FZF. Default is set to false.")
	isList := flag.Bool("l", false, "Return the list of ssh connections.")
	isListFuzzy := flag.Bool("lf", false, "Return the list of ssh connections and apply fuzzysearch")
	isColor := flag.Bool("c", false, "Enable a colored output")
	isListExecutable := flag.Bool("le", false, "Search inside the list of connections using fuzzysearch and start a SSH connection")
	isSync := flag.Bool("s", false, "Sync new updates in the repository fetching new files from Gitlab")
	isHelp := flag.Bool("h", false, "Print the help")

	// Parse flags
	flag.Parse()

	// Read the configuration file
	configuration, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	// Init hssh instance
	var hsshInstance = hssh.HSSH{
		Configuration: configuration,
	}

	// Command assignation
	command := ""

	if *isListFuzzy {
		*isFuzzy = true
	}

	if *isList || *isListFuzzy {
		command = "list"
	}

	if *isSync {
		command = "sync"
	}

	if *isListExecutable {
		command = "exec"
	}

	if *isHelp {
		command = "help"
	}

	switch command {
	case "sync":
		hsshInstance.Sync()
		break
	case "list":
		result := hsshInstance.List(*isFuzzy)
		out(result, *isColor)
		break
	case "exec":
		hsshInstance.Exec()
		break
	case "help":
		help()
		break
	default:
		fmt.Println("Invalid action")
		os.Exit(1)
	}

}

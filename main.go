package main

import (
	"flag"
	"fmt"
	"hssh/cli"
	"hssh/config"
	"hssh/connections"
	"hssh/providers"
	"hssh/templates"
	"os"
)

func printHelp() {
	fmt.Println(templates.Help)
}

func main() {
	withFuzzysearch := flag.Bool("f", false, templates.MsgFuzzySearchFlag)
	isList := flag.Bool("l", false, templates.MsgListFlag)
	isColor := flag.Bool("c", false, templates.MsgColorFlag)
	isExec := flag.Bool("e", false, templates.MsgExecFlag)
	isSync := flag.Bool("s", false, templates.MsgSyncFlag)
	isNewConfig := flag.Bool("C", false, templates.MsgNewConfigFlag)
	isHelp := flag.Bool("h", false, templates.MsgHelpFlag)
	flag.Parse()

	conf := config.New()

	err := conf.Load()
	if err != nil {
		fmt.Println(templates.ErrLoadConfiguration, err)
	}

	providerConnectionString := conf.GetProvider()

	p := providers.New(
		providerConnectionString,
	)

	sshUA := connections.NewSSHUA()

	fuzzysearch := conf.GetFuzzysearch()
	if *withFuzzysearch == false && *isExec == false {
		fuzzysearch = ""
	}

	if fuzzysearch == "" && (*isExec == true || *withFuzzysearch == true) {
		fmt.Println(templates.ErrInvalidFuzzysearchBInary)
		os.Exit(1)
	}

	c := cli.New(
		fuzzysearch,
		p,
		sshUA,
		*isColor,
	)

	if *isHelp == true {
		printHelp()
		os.Exit(0)
	}

	if *isList == true {
		out, _ := c.List()
		c.Print(out)
		os.Exit(0)
	}

	if *isExec == true {
		c.Connect()
		os.Exit(0)
	}

	if *isSync == true {
		c.Sync(providerConnectionString)
		os.Exit(0)
	}

	if *isNewConfig == true {
		conf.Create(templates.Config)
		os.Exit(0)
	}

	printHelp()

}

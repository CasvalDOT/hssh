package main

import (
	"flag"
	"fmt"
	"os"

	"hssh/cli"
	"hssh/config"
	"hssh/providers"
	"hssh/sshuseragent"
	"hssh/templates"
)

func main() {
	withFuzzysearch := flag.Bool("f", false, templates.MsgFuzzySearchFlag)
	isList := flag.Bool("l", false, templates.MsgListFlag)
	isColor := flag.Bool("c", false, templates.MsgColorFlag)
	isSync := flag.Bool("s", false, templates.MsgSyncFlag)
	isNewConfig := flag.Bool("C", false, templates.MsgNewConfigFlag)
	isHelp := flag.Bool("h", false, templates.MsgHelpFlag)
	flag.Parse()

	conf := config.New()

	/*
		Instead of generate an error
		and leave the user "alone" with debug or TODOs,
		if the configuration file cannot be found we can create an empty one
	*/
	err := conf.Load()
	if err != nil {
		fmt.Println(templates.ErrLoadConfiguration, err)
		conf.Create(templates.Config)
	}

	providerConnectionString := conf.GetProvider()
	p := providers.New(
		providerConnectionString,
	)

	sshUA := sshuseragent.NewSSHUserAgent()

	fuzzysearch := conf.GetFuzzysearch()
	if *isList == true && *withFuzzysearch == false {
		fuzzysearch = ""
	}

	if fuzzysearch == "" && *withFuzzysearch == true {
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
		fmt.Println(templates.Help)
		os.Exit(0)
	}

	if *isList == true {
		out, _ := c.List()
		c.Print(out)
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

	if fuzzysearch != "" {
		c.Connect()
		os.Exit(0)
	}

}

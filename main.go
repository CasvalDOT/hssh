package main

import (
	"flag"
	"fmt"
	"hssh/config"
	"hssh/engine"
	"hssh/providers"
	"hssh/templates"
	"os"
)

func printHelp() {
	fmt.Println(templates.Help)
}

func main() {
	withFuzzyEngine := flag.Bool("f", false, templates.MsgFuzzySearchFlag)
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
		fmt.Println(err)
	}

	providerConfig := conf.GetProvider()
	defaultProvider := conf.GetDefaultProvider()

	p := providers.New(
		defaultProvider,
		providerConfig.Host,
		providerConfig.PrivateToken,
	)

	fuzzyEngine := conf.GetFuzzyEngine()
	if *withFuzzyEngine == false {
		fuzzyEngine = ""
	}

	e := engine.New(
		fuzzyEngine,
		p,
		providerConfig.Files,
		*isColor,
	)

	if *isExec == true {
		fuzzyEngine = conf.GetFuzzyEngine()
	}

	if *isHelp == true {
		printHelp()
		os.Exit(0)
	}

	if *isList == true {
		out, _ := e.List()
		e.Print(out)
		os.Exit(0)
	}

	if *isExec == true {
		e.Exec()
		os.Exit(0)
	}

	if *isSync == true {
		e.Sync(providerConfig.ProjectID)
		os.Exit(0)
	}

	if *isNewConfig == true {
		conf.Create(templates.Config)
		os.Exit(0)
	}

	printHelp()

}

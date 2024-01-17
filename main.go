package main

import (
	"flag"
	"fmt"
	"os"

	"ekival-canvas/config"
	"ekival-canvas/router"
)

var cmdlineFlags struct {
	configFile string
}

func main() {
	// Load config
	flag.StringVar(
		&cmdlineFlags.configFile,
		"config",
		"",
		"path to config file to load",
	)
	flag.Parse()

	err := config.Load(cmdlineFlags.configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err)
		os.Exit(1)
	}

	config.WalletSetup()

	router.RouterInit()

}

package main

import (
	"fmt"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

func main() {
	app := cli.NewCli(eprinttools.Version)

	// Document non-option Command Line Parameters

	// Add Help Docs

	// Setup Environment

	// Setup Standard Options

	// Setup App Options

	// Process options and environment

	// Setup IO

	// Process options

	// Run App
	fmt.Printf("%s Not implemented\n", app.AppName())
}

// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2022, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package main

//
// ep3genfeeds implements feed generator rendering JSON documents and
// non-templated Markdown documents to a directory structure in an 
// htdoc folder specified in the configuration file.
//

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `%% {app_name}(1) user manual
%% R. S. Doiel
%% 2022-11-28

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTION] JSON_SETTINGS_FILENAME

# DESCRIPTION

{app_name} is a command line program that renders the EPrint harvested
metadata and aggregation tables to JSON documents, non-templated
Markdown documents and the necessary directory structures needed for
representing EPrints repositories as a static site.

The configuration needs to be previously created using the 
ep3harvester tool.

# OPTIONS

`

	examples = `
# EXAMPLES

Harvesting repositories for week month of May, 2022.

~~~
    {app_name} harvester-settings.json
~~~

`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
	verbose bool

	// App Option
	people bool
	groups bool
)

func main() {
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")
	flagSet.BoolVar(&verbose, "verbose", false, "use verbose logging")
	flagSet.BoolVar(&people, "people", false, "render people feeds")
	flagSet.BoolVar(&groups, "groups", false, "render groups feeds")


	// We're ready to process args
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	out := os.Stdout

	// Handle options
	if showHelp {
		eprinttools.DisplayUsage(out, appName, flagSet, description, examples)
		os.Exit(0)
	}
	if showLicense {
		eprinttools.DisplayLicense(out, appName)
		os.Exit(0)
	}
	if showVersion {
		eprinttools.DisplayVersion(out, appName)
		os.Exit(0)
	}
	settings := ""
	if len(args) > 0 {
		settings = args[0]
	}

	t0 := time.Now()
	switch {
		case people:
			if err := eprinttools.RunGenPeople(settings, verbose); err != nil {
				log.Print(err)
				os.Exit(1)
			}
		case groups:
			if err := eprinttools.RunGenGroups(settings, verbose); err != nil {
				log.Print(err)
				os.Exit(1)
			}
		default:
			if err := eprinttools.RunGenfeeds(settings, verbose); err != nil {
				log.Print(err)
				os.Exit(1)
			}
	}
	log.Printf("Total run time %v", time.Now().Sub(t0).Truncate(time.Second))
}

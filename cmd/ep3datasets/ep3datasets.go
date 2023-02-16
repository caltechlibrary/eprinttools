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
// ep3datasets renders dataset collections from previously harvested
// EPrints repositories based on the settings.json configuration file.
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	helpText = `---
title: "{app_name} (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-02-07
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTION] JSON_SETTINGS_FILE

# DESCRIPTION

{app_name} is a command line program renders dataset collections
from previously harvested EPrint repositories based on the
configuration in the JSON_SETTINGS_FILE.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-verbose
: use verbose logging

-repo
: write out the dataset for a specific repo in JSON_SETTINGS_FILE

# EXAMPLES

Rendering harvested repositories for settings.json.

~~~
    {app_name} settings.json
~~~

Render the harvested repository caltechauthors based on settings.json.

~~~
	{app_name} -repo caltechauthors settings.json
~~~

{app_name} {version}

`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Option
	verbose bool
	repoName string
	simplified bool
)

func fmtTxt(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&verbose, "verbose", false, "use verbose logging")
	flag.BoolVar(&simplified, "simplified", false, "use a simplified records structure in collection")
	flag.StringVar(&repoName, "repo", "", "Harvest a specific repository id defined in configuration")
	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	// Setup I/O
	var err error

	//in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	// Handle options
	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, eprinttools.Version)		)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", eprinttools.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s\n", appName, eprinttools.Version)
		os.Exit(0)
	}
	settings := ""
	if len(args) > 0 {
		settings = args[0]
	}

	t0 := time.Now()
	if repoName != "" {
		err = eprinttools.RunDataset(settings, repoName, verbose)
	} else {
		err = eprinttools.RunDatasets(settings, verbose)
	}
	if err != nil {
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	log.Printf("Total run time %v", time.Now().Sub(t0).Truncate(time.Second))
}

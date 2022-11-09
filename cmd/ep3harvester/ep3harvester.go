// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
// ep3harvester implements an EPrints Metadata to JSON store harvester.
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `
USAGE

	{app_name} [OPTION] JSON_SETTINGS_FILENAME \
	           [START_TIMESTAMP] [END_TIMESTAMP]

{app_name} is a command line program for metadata harvester
for EPrints repositories.

DESCRIPTION

{app_name} takes a JSON settings file and harvests
all the EPrint repositories defined in the settings file
into a JSON store implemented in MySQL 8. One repository per
MySQL 8 table.

Each MySQL 8 table has two columns id, src (holding the JSON
document as a JSON column).

CONFIGURING YOUR JSON STORE

{app_name} uses a MySQL 8 database for a JSON document store.
It will generate one table for EPrint repository. You can
generate a SQL program for creating the MySQL database and
tables from your settings JSON file using the "-sql-schema"
option. Using the option will require a JSON settings filename
parameter. E.g.

    {app_name} -sql-schema settings.json

OPTIONS
`

	examples = `
EXAMPLES

Harvesting repositories for week month of May, 2022.

    {app_name} settings.json "2022-05-01 00:00:00" "2022-05-31 59:59:59"

`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Option
	showSqlSchema bool
)

func main() {
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")
	flagSet.BoolVar(&showSqlSchema, "sql-schema", false, "display SQL schema for installing MySQL jsonstore DB")

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
	settings, start, end := "", "", ""
	if len(args) > 0 {
		settings = args[0]
	}
	if len(args) > 1 {
		start = args[1]
	}
	if len(args) > 2 {
		end = args[2]
	}

	// Handle request to show schema.
	if showSqlSchema {
		if settings == "" {
			log.Printf("%s -sql-schema option requires a settings file parameter", appName)
			os.Exit(1)
		}
		src, err := eprinttools.HarvesterInitDB(settings)
		if err != nil {
			log.Printf("%s -sql-schema error: %s", appName, err)
			os.Exit(1)
		}
		fmt.Print(src)
		os.Exit(0)
	}
	err := eprinttools.RunHarvester(settings, start, end)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

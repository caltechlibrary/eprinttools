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
	"time"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	helpText = `---
title: "{app_name} (1) user manual"
pubDate: 2023-02-16
author: "R. S. Doiel"
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] JSON_SETTINGS_FILENAME \
           [START_TIMESTAMP] [END_TIMESTAMP]

# DESCRIPTION

{app_name} is a command line program for metadata harvesting
of EPrints repositories.

{app_name} takes a JSON settings file and harvests
all the EPrint repositories defined in the settings file
into a JSON store implemented in MySQL 8. One repository per
MySQL 8 table.

Each MySQL 8 table has several columns id, src (holding the JSON
document as a JSON column) and an updated (holding the timestamp
of when the metadata was harvested).

# CONFIGURATION

{app_name} can generate an example settings JSON document. You
can then edit it with any plain text editor (e.g. nano). Then
you'll need to setup a MySQL 8 database and tables to store
havested data in.

{app_name} uses a MySQL 8 database for a JSON document store.
It will generate one table for EPrint repository. You can
generate a SQL program for creating the MySQL database and
tables from your settings JSON file using the "-sql-schema"
option. Using the option will require a JSON settings filename
parameter. E.g.

~~~
    {app_name} -init harvester-settings.json
    nano harvester-settings.json
    {app_name} -sql-schema harvester-settings.json >collections.sql
~~~

# OPTIONS

-help
: display help

-version
: display version

-license
: display license

-groups
: Harvest groups from CSV files included configuration

-init
: generate a settings JSON file

-people
: Harvest people from CSV files included configuration

-people-groups
: Harvest people and groups from CSV files included configuration

-repo string
: Harvest a specific repository id defined in configuration

-simple
: Crosswalk the harvested eprint record to the simplified record model
before saving the JSON to the SQL database.

-sql-schema
: display SQL schema for installing MySQL jsonstore DB

-verbose
: use verbose logging

# EXAMPLES

Harvesting repositories for the month of May, 2022.

~~~
    {app_name} harvester-settings.json \
        "2022-05-01 00:00:00" "2022-05-31 59:59:59"
~~~

Harvesting a caltechauthors repo using harvester-settings.json
for week month of the month of May, 2022.

~~~
	{app_name} -repo caltechauthors harvester-settings.json \ 
        "2022-05-01 00:00:00" "2022-05-31 59:59:59"
~~~

{app_name} {version}

`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Option
	showSqlSchema   bool
	initialize      bool
	people          bool
	groups          bool
	repoName        string
	peopleAndGroups bool
	useSimplifiedRecord bool
	verbose         bool
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
	flag.BoolVar(&showSqlSchema, "sql-schema", false, "display SQL schema for installing MySQL jsonstore DB")
	flag.BoolVar(&verbose, "verbose", false, "use verbose logging")
	flag.BoolVar(&initialize, "init", false, "generate a settings JSON file")
	flag.BoolVar(&people, "people", false, "Harvest people from CSV files included configuration")
	flag.BoolVar(&groups, "groups", false, "Harvest groups from CSV files included configuration")
	flag.BoolVar(&peopleAndGroups, "people-groups", false, "Harvest people and groups from CSV files included configuration")
	flag.BoolVar(&useSimplifiedRecord, "simple", false, "Crosswalk harvested eprint records storing simplified model")
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
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, eprinttools.Version))
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

	if initialize {
		if settings != "" {
			var err error
			out, err = os.Create(settings)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer out.Close()
		}
		fmt.Fprintf(out, "%s\n", eprinttools.DefaultConfig())
		os.Exit(0)
	}
	if settings == "" {
		log.Printf("%s a settings file parameter", appName)
		os.Exit(1)
	}

	t0 := time.Now()
	// Handle request to show schema.
	switch {
	case showSqlSchema:
		src, err := eprinttools.HarvesterDBSchema(settings)
		if err != nil {
			fmt.Fprintf(eout, "%s -sql-schema error: %s", appName, err)
			os.Exit(1)
		}
		fmt.Fprintf(out, "%s\n", src)
		os.Exit(1)
	case people:
		err = eprinttools.RunHarvestPeople(settings, verbose)
	case groups:
		err = eprinttools.RunHarvestGroups(settings, verbose)
	case peopleAndGroups:
		err = eprinttools.RunHarvestPeople(settings, verbose)
		if err != nil {
			fmt.Fprintln(eout, err)
			os.Exit(1)
		}
		err = eprinttools.RunHarvestGroups(settings, verbose)
	case repoName != "":
		err = eprinttools.RunHarvestRepoID(settings, repoName, start, end, verbose)
	default:
		err = eprinttools.RunHarvester(settings, start, end, useSimplifiedRecord, verbose)
	}
	if err != nil {
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	log.Printf("total run time %v", time.Now().Sub(t0).Truncate(time.Second))
	os.Exit(0)
}

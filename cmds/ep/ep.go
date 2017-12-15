//
// Package ep is a collection of structures and functions for working with the EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	ep "github.com/caltechlibrary/eprinttools"
)

var (
	// cli help text

	description = `

SYNOPSIS

%s wraps the REST API for EPrints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

CONFIGURATION

ep can be configured with following environment variables

EP_EPRINTS_URL the URL to your EPrints installation

EP_DATASET the dataset and collection name for exporting, site building, and content retrieval

`

	examples = `

EXAMPLE

    %s -export all

Would export the entire EPrints repository public content defined by the
environment virables EP_API_URL, EP_DATASET.

    %s -export 2000

Would export 2000 EPrints from the repository with the heighest ID values.

   %s -export-modified 2017-07-01

Would export the EPrint records modified since July 1, 2017.

   %s -export-modified 2017-07-01,2017-07-31 \
      -export-save-keys=july-keys.txt 

Would export the EPrint records with modified times in July 2017 and
save the keys for the records exported with one key per line. 

SUPPRESSING NOTES FIELD

Sometimes the notes field is used for internal processing and should
not be harvested. If this is the case for you use the "-suppress-notes"
option.
`

	// Standard Options
	showHelp             bool
	showVersion          bool
	showLicense          bool
	showExamples         bool
	outputFName          string
	quiet                bool
	generateMarkdownDocs bool

	// App Options
	verbose     bool
	useAPI      bool
	prettyPrint bool

	apiURL      string
	datasetName string

	updatedSince          string
	exportEPrints         string
	exportEPrintsModified string
	exportSaveKeys        string

	authMethod string
	userName   string
	userSecret string

	// NOTE: supressNote added to handle the case where Note field is internal use only and not to be harvested
	suppressNote bool
)

func main() {
	app := cli.NewCli(ep.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.AddParams("[EP_EPRINTS_URL|ONE_OR_MORE_EPRINT_ID]")

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(ep.LicenseText, appName, ep.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName, appName)))

	// App Environment
	app.EnvStringVar(&apiURL, "EP_EPRINT_URL", "", "Sets the EPRints API URL")
	app.EnvStringVar(&datasetName, "EP_DATASET", "", "Sets the dataset collection for storing EPrint harvested records")
	app.EnvBoolVar(&suppressNote, "EP_SUPPRESS_NOTE", false, "Suppress the note field on harvesting")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "generation markdown documentation")
	app.BoolVar(&quiet, "quiet", false, "suppress error output")
	//app.BoolVar(&newLine, "nl,newline", true, "include trailing newline in output")

	// App Specific options
	app.BoolVar(&suppressNote, "suppress-note", false, "suppress note")
	app.BoolVar(&verbose, "verbose", true, "verbose logging")
	app.StringVar(&authMethod, "auth", "", "set the authentication method (e.g. none, basic, oauth, shib)")
	app.StringVar(&userName, "un,username", "", "set the username")
	app.StringVar(&userSecret, "pw,password", "", "set the password")

	app.StringVar(&apiURL, "api", "", "url for EPrints API")
	app.StringVar(&datasetName, "dataset", "", "dataset/collection name")

	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print JSON output")
	app.BoolVar(&useAPI, "read-api", false, "read the contents from the API without saving in the database")
	app.StringVar(&exportEPrints, "export", "", "export N EPrints from highest ID to lowest")
	app.StringVar(&exportEPrintsModified, "export-modified", "", "export records by date or date range (e.g. 2017-07-01)")
	app.StringVar(&exportSaveKeys, "export-save-keys", "", "save the keys exported in a file with provided filename")
	app.StringVar(&updatedSince, "updated-since", "", "list EPrint IDs updated since a given date (e.g 2017-07-01)")

	// Parse environment and options
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.Out, err = cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, app.In)

	// Set log to output
	log.SetOutput(app.Out)

	// Process Options
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Println(app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showVersion {
		fmt.Println(app.Version())
		os.Exit(0)
	}
	if showLicense {
		fmt.Println(app.License())
		os.Exit(0)
	}

	// Required configuration
	if apiURL == "" {
		fmt.Fprintf(app.Eout, "EPrint URL not provided\n")
		os.Exit(1)
	}
	if datasetName == "" {
		fmt.Fprintf(app.Eout, "Missing dataset (EP_DATASET) name\n")
		os.Exit(1)
	}

	// This will read in the settings from the app
	// and configure access to the EPrints API
	api, err := ep.New(apiURL, datasetName, suppressNote, authMethod, userName, userSecret)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}

	if exportEPrints != "" {
		t0 := time.Now()
		exportNo := -1
		if exportEPrints != "all" {
			exportNo, err = strconv.Atoi(exportEPrints)
			if err != nil {
				log.Fatalf("Export count should be %q or an integer, %s", exportEPrints, err)
			}
		}
		log.Printf("%s %s (pid %d)", appName, ep.Version, os.Getpid())
		log.Printf("Export started, %s", t0)
		if err := api.ExportEPrints(exportNo, exportSaveKeys, verbose); err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Export completed, running time %s", time.Now().Sub(t0))
		os.Exit(0)
	}
	if exportEPrintsModified != "" {
		s := exportEPrintsModified
		e := time.Now().Format("2006-01-02")
		if strings.Contains(s, ",") {
			p := strings.SplitN(s, ",", 2)
			s, e = p[0], p[1]
		}
		start, err := time.Parse("2006-01-02", s)
		if err != nil {
			log.Fatalf("%s", err)
		}
		end, err := time.Parse("2006-01-02", e)
		if err != nil {
			log.Fatalf("%s", err)
		}
		t0 := time.Now()
		log.Printf("%s %s (pid %d)", appName, ep.Version, os.Getpid())
		log.Printf("Export from %s to %s, started %s", start.Format("2006-01-02"), end.Format("2006-01-02"), t0.Format("2006-01-02 15:04:05 MST"))
		if err := api.ExportModifiedEPrints(start, end, exportSaveKeys, verbose); err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Export completed, running time %s", time.Now().Sub(t0))
		os.Exit(0)
	}

	//
	// Generate JSON output
	//
	var (
		src  []byte
		data interface{}
	)
	switch {
	case updatedSince != "":
		// date should be formatted YYYY-MM-DD, 2006-01-02
		end := time.Now()
		start, err := time.Parse("2006-01-02", updatedSince)
		if err != nil {
			fmt.Fprintf(os.Stderr, "updated since %q, %s", updatedSince, err)
			os.Exit(1)
		}
		data, err = api.ListModifiedEPrintURI(start, end, verbose)
	case useAPI == true:
		if len(args) == 1 {
			data, _, err = api.GetEPrint(args[0])
		} else {
			data, err = api.ListEPrintsURI()
		}
	default:
		if len(args) == 1 {
			data, err = api.Get(args[0])
		} else if len(args) > 1 {
			records := []*ep.Record{}
			for _, id := range args {
				if rec, err := api.Get(id); err == nil {
					records = append(records, rec)
				} else {
					fmt.Fprintf(os.Stderr, "Can't read EPrint id %s, %s\n", id, err)
				}
			}
			data = records
		} else {
			data, err = api.ListID(0, -1)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	if prettyPrint == true {
		src, _ = json.MarshalIndent(data, "", "    ")
	} else {
		src, _ = json.Marshal(data)
	}
	fmt.Fprintf(app.Out, "%s", src)
}

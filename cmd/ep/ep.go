//
// ep is an eprinttools based cli for harvesting EPrints content into a dataset collection.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2018, Caltech
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
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/eprinttools/harvest"
)

var (
	// cli help text
	description = `
%s uses the REST API for EPrints 3.x to harvest EPrints content into
a dataset collection. If you don't need dataset integration use eputil 
instead. If you want to view  the harvested content then use the
dataset command.

CONFIGURATION

ep can be configured with following environment variables

EPRINT_URL the URL to your EPrints installation

DATASET the dataset collection name to use for storing your harvested EPrint content.
`

	examples = `
Save a list the URI end points for eprint records found at EPRINT_URL.

	%s -o uris.txt

Export the entire EPrints repository public content defined by the
environment variables EPRINT_URL, DATASET.

    %s -export all

Export 2000 EPrints from the repository with the heighest ID values.

    %s -export 2000

Export the EPrint records modified since July 1, 2017.

    %s -export-modified 2017-07-01

Explore a specific list of keys (e.g. "101,102,1304" or
if list is '-' then read from standard input, one key per line)

	%s -export-keys "101,102,1304"

	%s -export-keys "-"

Export the EPrint records with modified times in July 2017 and
save the keys for the records exported with one key per line. 

    %s -export-modified 2017-07-01,2017-07-31 \
       -export-save-keys=july-keys.txt 
`

	// Standard Options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	outputFName      string
	quiet            bool
	generateMarkdown bool
	generateManPage  bool
	newLine          bool

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
	exportEPrintsKeyList  string
	exportEPrintDocs      bool

	authMethod string
	userName   string
	userSecret string

	// NOTE: suppressSuggestions (Internal Note) added to handle the case
	// where Suggestions field is internal use only and not to be harvested
	suppressSuggestions bool

	thisProcessID int
)

func main() {
	var (
		apiURLEnv      string
		datasetNameEnv string
	)
	app := cli.NewCli(eprinttools.Version)
	appName := app.AppName()
	thisProcessID = os.Getpid()

	// Document non-option parameters
	app.SetParams("[EPRINT_URL]", "[ONE_OR_MORE_EPRINT_ID]")

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(eprinttools.LicenseText, appName, eprinttools.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName, appName, appName, appName, appName)))

	// App Environment
	app.EnvStringVar(&apiURLEnv, "EPRINT_URL", "", "Sets the EPRints API URL")
	app.EnvStringVar(&datasetNameEnv, "DATASET", "", "Sets the dataset collection for storing EPrint harvested records")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generation markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generation man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error output")
	app.BoolVar(&newLine, "nl,newline", true, "set to false to exclude trailing newline")

	// App Specific options
	app.BoolVar(&verbose, "verbose", true, "verbose logging")
	app.StringVar(&authMethod, "auth", "", "set the authentication method (e.g. none, basic, oauth, shib)")
	app.StringVar(&userName, "un,username", "", "set the username")
	app.StringVar(&userSecret, "pw,password", "", "set the password")

	app.StringVar(&apiURL, "api", "", "url for EPrints API")
	app.StringVar(&datasetName, "dataset", "", "dataset collection name")

	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print JSON output")
	app.BoolVar(&useAPI, "read-api", false, "read the contents from the API without saving in the database")
	app.StringVar(&exportEPrints, "export", "", "export N EPrints from highest ID to lowest")
	app.StringVar(&exportEPrintsModified, "export-modified", "", "export records by date or date range (e.g. 2017-07-01)")
	app.StringVar(&exportSaveKeys, "export-save-keys", "", "save the keys exported in a file with provided filename")
	app.StringVar(&exportEPrintsKeyList, "export-keys", "", "export using a delimited list of EPrint keys")
	app.BoolVar(&exportEPrintDocs, "export-with-docs", false, "include EPrint documents with export")
	app.StringVar(&updatedSince, "updated-since", "", "list EPrint IDs updated since a given date (e.g 2017-07-01)")
	app.BoolVar(&suppressSuggestions, "suppress-suggestions", true, "suppress the suggestions field from output")

	// Parse environment and options
	if err := app.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong parsing env and options!, %s\n", err)
		os.Exit(1)
	}
	args := app.Args()
	if apiURL == "" {
		apiURL = app.Getenv("EPRINT_URL")
		if len(args) > 0 {
			for _, val := range args {
				if strings.Contains(val, "://") {
					apiURL = val
					break
				}
			}
		}
	}
	if datasetName == "" {
		datasetName = app.Getenv("DATASET")
	}

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.Out, err = cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, app.In)

	// Process Options
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
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

	if exportEPrintDocs {
		harvest.ExportEPrintDocs = true
	}

	// Required configuration, let Env overide options if not options are defaults
	if apiURL == "" {
		if apiURLEnv == "" {
			fmt.Fprintf(app.Eout, "EPrint URL not provided, -api ..., 'export EPRINT_URL=...'\n")
			os.Exit(1)
		}
		apiURL = apiURLEnv
	}
	if datasetName == "" {
		if datasetNameEnv == "" {
			fmt.Fprintf(app.Eout, "Missing dataset (e.g. --dataset ..., 'export DATASET=...') name\n")
			os.Exit(1)
		}
		datasetName = datasetNameEnv
	}

	// This will read in the settings from the app
	// and configure access to the EPrints API
	u, err := url.Parse(apiURL)
	if err == nil {
		if strings.HasSuffix(u.Path, "/rest/eprint") == true {
			u.Path = strings.TrimSuffix(u.Path, "/rest/eprint")
			apiURL = u.String()
		}
	} else {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	api, err := eprinttools.New(apiURL, datasetName, suppressSuggestions, authMethod, userName, userSecret)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}

	// NOTE: Sanitize u (the URL struct) so we can log it rather than the raw apiURL.
	u.User = nil
	log.Printf("(pid: %d) %s %s", thisProcessID, appName, eprinttools.Version)
	log.Printf("(pid: %d) Harvesting from %s", thisProcessID, u.String())
	t0 := time.Now()
	switch {
	case exportEPrintsKeyList != "":
		//FIXME: need to see if this is a short comma limited list
		// provided on the command line or if we're reading from a file
		// or standard in.
		log.Printf("(pid: %d) Export started, %s", thisProcessID, t0)
		keys_to_export := []string{}
		if exportEPrintsKeyList == "-" {
			scanner := bufio.NewScanner(app.In)
			for scanner.Scan() {
				key := scanner.Text()
				key = strings.TrimSpace(key)
				if key != "" {
					keys_to_export = append(keys_to_export, key)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Fatalf("(pid: %d) %s")
			}
		} else {
			for _, key := range strings.Split(exportEPrintsKeyList, ",") {
				key = strings.TrimSpace(key)
				if key != "" {
					keys_to_export = append(keys_to_export, key)
				}
			}
		}
		if err := harvest.ExportEPrintsKeyList(api, keys_to_export, exportSaveKeys, verbose); err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		}
	case exportEPrintsModified != "":
		s := exportEPrintsModified
		e := time.Now().Format("2006-01-02")
		if strings.Contains(s, ",") {
			p := strings.SplitN(s, ",", 2)
			s, e = p[0], p[1]
		}
		start, err := time.Parse("2006-01-02", s)
		if err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		}
		end, err := time.Parse("2006-01-02", e)
		if err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		}
		log.Printf("(pid: %d) Export from %s to %s, started %s", thisProcessID, start.Format("2006-01-02"), end.Format("2006-01-02"), t0.Format("2006-01-02 15:04:05 MST"))
		if err := harvest.ExportModifiedEPrints(api, start, end, exportSaveKeys, verbose); err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		}
	case exportEPrints != "":
		exportNo := -1
		if exportEPrints != "all" {
			exportNo, err = strconv.Atoi(exportEPrints)
			if err != nil {
				log.Fatalf("(pid: %d) Export count should be %q or an integer, %s", thisProcessID, exportEPrints, err)
			}
		}
		log.Printf("(pid: %d) Export started, %s", thisProcessID, t0)
		if err := harvest.ExportEPrints(api, exportNo, exportSaveKeys, verbose); err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		}
	default:
		if uris, err := api.ListEPrintsURI(); err != nil {
			log.Fatalf("(pid: %d) %s", thisProcessID, err)
		} else {
			fmt.Fprintf(app.Out, "%s\n", strings.Join(uris, "\n"))
		}
	}
	log.Printf("(pid: %d) Export completed, running time %s", thisProcessID, time.Now().Sub(t0).Round(time.Second))
	os.Exit(0)
}

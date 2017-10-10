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
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	ep "github.com/caltechlibrary/eprinttools"
)

var (
	// cli help text
	usage = `USAGE %s [OPTIONS] [EP_EPRINTS_URL|ONE_OR_MORE_EPRINT_ID]`

	description = `
SYNOPSIS

%s wraps the REST API for EPrints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

CONFIGURATION

ep can be configured with following environment variables

EP_EPRINTS_URL the URL to your EPrints installation

EP_DATASET the dataset and collection name for exporting, site building, and content retrieval`

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
`

	// Standard Options
	showHelp     bool
	showVersion  bool
	showLicense  bool
	showExamples bool
	outputFName  string
	verbose      bool

	// App Options
	useAPI      bool
	prettyPrint bool

	apiURL      string
	datasetName string

	updatedSince          string
	exportEPrints         string
	exportEPrintsModified string
	exportSaveKeys        string
	feedSize              int

	authMethod string
	userName   string
	userSecret string
)

func init() {
	// Setup options
	feedSize = ep.DefaultFeedSize

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")
	flag.StringVar(&outputFName, "o", "", "output filename (logging)")
	flag.StringVar(&outputFName, "output", "", "output filename (logging)")
	flag.BoolVar(&verbose, "verbose", true, "verbose logging")

	// App Specific options
	flag.StringVar(&authMethod, "auth", "", "set the authentication method (e.g. none, basic, oauth, shib)")
	flag.StringVar(&userName, "username", "", "set the username")
	flag.StringVar(&userName, "un", "", "set the username")
	flag.StringVar(&userSecret, "pw", "", "set the password")

	flag.StringVar(&apiURL, "api", "", "url for EPrints API")
	flag.StringVar(&datasetName, "dataset", "", "dataset/collection name")

	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
	flag.BoolVar(&prettyPrint, "pretty", false, "pretty print JSON output")
	flag.BoolVar(&useAPI, "read-api", false, "read the contents from the API without saving in the database")
	flag.StringVar(&exportEPrints, "export", "", "export N EPrints from highest ID to lowest")
	flag.StringVar(&exportEPrintsModified, "export-modified", "", "export records by date or date range (e.g. 2017-07-01)")
	flag.StringVar(&exportSaveKeys, "export-save-keys", "", "save the keys exported in a file with provided filename")
	flag.StringVar(&updatedSince, "updated-since", "", "list EPrint IDs updated since a given date (e.g 2017-07-01)")
}

func check(cfg *cli.Config, key, value string) string {
	if value == "" {
		log.Fatalf("Missing %s_%s", cfg.EnvPrefix, strings.ToUpper(key))
		return ""
	}
	return value
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	// Populate cfg from the environment
	cfg := cli.New(appName, "EP", ep.Version)
	cfg.LicenseText = fmt.Sprintf(ep.LicenseText, appName, ep.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.OptionText = "OPTIONS"
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName)

	// Handle the default options
	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}

	if showExamples == true {
		if len(args) > 0 {
			fmt.Println(cfg.Example(args...))
		} else {
			fmt.Println(cfg.ExampleText)
		}
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	out, err := cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprint(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, out)

	// Log to out
	log.SetOutput(out)

	// Required configuration
	apiURL = check(cfg, "eprint_url", cfg.MergeEnv("eprint_url", apiURL))
	datasetName = check(cfg, "dataset", cfg.MergeEnv("dataset", datasetName))

	// Optional configuration
	authMethod = cfg.MergeEnv("auth_method", authMethod)
	userName = cfg.MergeEnv("username", userName)
	userSecret = cfg.MergeEnv("password", userSecret)

	// This will read in any settings from the environment
	api, err := ep.New(cfg)
	if err != nil {
		log.Fatal(err)
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
			log.Printf("%s", err)
			os.Exit(1)
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
			log.Printf("%s", err)
			os.Exit(1)
		}
		end, err := time.Parse("2006-01-02", e)
		if err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
		t0 := time.Now()
		log.Printf("%s %s (pid %d)", appName, ep.Version, os.Getpid())
		log.Printf("Export from %s to %s, started %s", start.Format("2006-01-02"), end.Format("2006-01-02"), t0.Format("2006-01-02 15:04:05 MST"))
		if err := api.ExportModifiedEPrints(start, end, exportSaveKeys, verbose); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
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
	fmt.Printf("%s", src)
}

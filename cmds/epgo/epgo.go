//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
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

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/epgo"
)

var (
	// cli help text
	usage = `USAGE %s [OPTIONS] [EPRINT_URI]`

	description = `
 %s wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
 a JSON view of the XML presentation as well as generates feeds and web pages.

 CONFIG

 %s can be configured with following environment variables

 + EPGO_API_URL (required) the URL to your E-Prints installation
 + EPGO_DATASET   (required) the dataset and collection name for exporting, site building, and content retrieval
 + EPGO_BLEVE (optional) the name for the Bleve index/db
 + EPGO_SITE_URL (optional) the URL to your public website (might be the same as E-Prints)
 + EPGO_HTDOCS   (optional) the htdocs root for site building
 + EPGO_TEMPLATE_PATH (optional) the template directory to use for site building

 If EPRINT_URI is provided then an individual EPrint is return as
 a JSON structure (e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are
 returned.

`

	examples = `
EXAMPLE

    %s -export all

Would export the entire EPrints repository public content defined by the
environment virables EPGO_API_URL, EPGO_DATASET.

    %s -export 2000

Would export 2000 EPrints from the repository with the heighest ID values.

    %s -select

Would (re)build the select lists based on contents of $EPGO_DATASET.

    %s -select -export all

Would export all eprints and rebuild the select lists.`

	// Standard Options
	showHelp    bool
	showVersion bool
	showLicense bool
	outputFName string

	// App Options
	useAPI      bool
	prettyPrint bool

	apiURL      string
	datasetName string

	exportEPrints   string
	feedSize        int
	publishedNewest int
	articlesNewest  int

	genSelectLists bool
)

func init() {
	// Setup options
	publishedNewest = 0
	feedSize = epgo.DefaultFeedSize

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.StringVar(&outputFName, "o", "", "output filename (logging)")
	flag.StringVar(&outputFName, "output", "", "output filename (logging)")

	// App Specific options
	flag.StringVar(&apiURL, "api", "", "url for EPrints API")
	flag.StringVar(&datasetName, "dataset", "", "dataset/collection name")

	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
	flag.BoolVar(&useAPI, "read-api", false, "read the contents from the API without saving in the database")
	flag.IntVar(&feedSize, "feed-size", feedSize, "number of items rendering in feeds")
	flag.StringVar(&exportEPrints, "export", "", "export N EPrints from highest ID to lowest")
	flag.IntVar(&publishedNewest, "published-newest", 0, "list the N newest published items")
	flag.IntVar(&articlesNewest, "articles-newest", 0, "list the N newest published articles")
	flag.BoolVar(&genSelectLists, "s", false, "generate select lists in dataset")
	flag.BoolVar(&genSelectLists, "select", false, "generate select lists in dataset")
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
	cfg := cli.New(appName, appName, fmt.Sprintf(epgo.LicenseText, appName, epgo.Version), epgo.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName, appName, appName)

	// Handle the default options
	if showHelp == true {
		fmt.Println(cfg.Usage())
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

	apiURL = check(cfg, "api_url", cfg.MergeEnv("api_url", apiURL))
	datasetName = check(cfg, "dataset", cfg.MergeEnv("dataset", datasetName))

	// This will read in any settings from the environment
	api, err := epgo.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if exportEPrints != "" {
		exportNo := -1
		if exportEPrints != "all" {
			exportNo, err = strconv.Atoi(exportEPrints)
			if err != nil {
				log.Fatalf("Export count should be %q or an integer, %s", exportEPrints, err)
			}
		}
		log.Printf("%s %s", appName, epgo.Version)
		log.Println("Export started")
		if err := api.ExportEPrints(exportNo); err != nil {
			log.Printf("%s", err)
		} else {
			log.Println("Export completed")
		}
		if genSelectLists != true {
			os.Exit(0)
		}
		log.Printf("Ready to run `%s -select` to rebuild select lists\n", appName)
	}
	if genSelectLists == true {
		log.Printf("%s %s", appName, epgo.Version)
		log.Println("Generating Select lists")
		api.BuildSelectLists()
		log.Println("Select lists complete")
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
	case publishedNewest > 0:
		data, err = api.GetPublications(0, publishedNewest)
	case articlesNewest > 0:
		data, err = api.GetArticles(0, articlesNewest)
	case useAPI == true:
		if len(args) == 1 {
			data, err = api.GetEPrint(args[0])
		} else {
			data, err = api.ListEPrintsURI()
		}
	default:
		if len(args) == 1 {
			data, err = api.Get(args[0])
		} else {
			data, err = api.ListURI(0, 1000000)
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

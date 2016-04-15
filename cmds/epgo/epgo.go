//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
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
	"io/ioutil"
	"log"
	"os"
	"strings"

	// 3rd Party Packages
	"github.com/robertkrimen/otto"

	// Caltech Library Packages
	"github.com/caltechlibrary/epgo"
	"github.com/caltechlibrary/ostdlib"
)

var (
	// CLI options
	showHelp       bool
	showVersion    bool
	showLicense    bool
	restAPI        bool
	prettyPrint    bool
	buildSite      bool
	runInteractive bool
	runJavaScript  bool

	exportEPrints   int
	feedSize        int
	publishedOldest int
	publishedNewest int
	articlesOldest  int
	articlesNewest  int
)

func init() {
	publishedNewest = 0
	publishedOldest = 0
	feedSize = epgo.DefaultFeedSize

	flag.BoolVar(&showHelp, "h", false, "display help info")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
	flag.BoolVar(&showLicense, "l", false, "show license information")
	flag.BoolVar(&runJavaScript, "js", false, "run JavaScript files")
	flag.BoolVar(&runInteractive, "i", false, "run interactive JavaScript REPL")
	flag.BoolVar(&restAPI, "api", false, "read the contents from the API without saving in the database")
	flag.BoolVar(&buildSite, "build", false, "build pages and feeds from database")
	flag.IntVar(&feedSize, "feed-size", feedSize, "number of items rendering in feeds")
	flag.IntVar(&exportEPrints, "export", 0, "export N EPrints from highest ID to lowest")
	flag.IntVar(&publishedOldest, "published-oldest", 0, "list the N oldest published items")
	flag.IntVar(&publishedNewest, "published-newest", 0, "list the N newest published items")
	flag.IntVar(&articlesOldest, "articles-oldest", 0, "list the N oldest published articles")
	flag.IntVar(&articlesNewest, "articles-newest", 0, "list the N newest published articles")
}

func main() {
	flag.Parse()
	if showHelp == true {
		fmt.Printf(`
 USAGE: epgo [OPTIONS] [EPRINT_URI|JAVASCRIPT_FILES]

 epgo wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
 a JSON view of the XML presentation as well as generates feeds and web pages.

 epgo can be configured with following environment variables

 + EPGO_API_URL (required) the URL to your E-Prints installation
 + EPGO_DBNAME   (required) the BoltDB name for exporting, site building, and content retrieval
 + EPGO_SITE_URL (optional) the URL to your public website (might be the same as E-Prints)
 + EPGO_HTDOCS   (optional) the htdocs root for site building
 + EPGO_TEMPLATES (optional) the template directory to use for site building

 If EPRINT_URI is provided then an individual EPrint is return as
 a JSON structure (e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are
 returned.

 OPTIONS

  -api	                       display EPrint REST API response
  -export int                  export N EPrints to local database,
                               if N is negative export all EPrints
  -build                       build pages and feeds from database
  -feed-size int               sets the number of items included in generated feeds

  -published-newest int        list the N newest published records
  -published-oldest int        list the N oldest published records
  -articles-newest int         list the N newest articles
  -articles-oldest int         list the N oldest articles

  -js                          run each JavaScript file, can be combined with -i
  -i                           interactive JavaScript REPL
  -p                           pretty print JSON output

  -h       display help info
  -l       show license information
  -v       display version info

 Version %s
`, epgo.Version)

		os.Exit(0)
	}

	if showVersion == true {
		fmt.Printf("Version %s\n", epgo.Version)
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Printf(`
 LICENSE

 Copyright (c) 2016, Caltech
 All rights not granted herein are expressly reserved by Caltech.

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are met:

 * Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

 * Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

 * Neither the name of epgo nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`)
		os.Exit(0)
	}

	// This will read in any settings from the environment
	api, err := epgo.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	args := flag.Args()
	if exportEPrints != 0 {
		if err := api.ExportEPrints(exportEPrints); err != nil {
			log.Fatalf("%s", err)
		}
		if buildSite == false {
			os.Exit(0)
		}
	}

	if buildSite == true {
		if err := api.BuildSite(feedSize); err != nil {
			log.Fatalf("%s", err)
		}
		os.Exit(0)
	}

	if runJavaScript == true || runInteractive == true {
		vm := otto.New()
		js := ostdlib.New(vm)
		// Add extensions
		js.AddExtensions()
		// Add EPrints API exentions
		vm = api.AddExtensions(js)

		if runJavaScript == true {
			for _, fname := range args {
				if strings.HasSuffix(fname, ".js") == true {
					src, err := ioutil.ReadFile(fname)
					if err != nil {
						log.Fatalf("%s, %s", fname, src)
					}
					// Load, Compile and run the JS file
					if script, err := vm.Compile(fname, src); err != nil {
						log.Fatalf("%s, %s", fname, err)
					} else {
						if _, err := vm.Run(script); err != nil {
							log.Fatalf("%s, %s", fname, err)
						}
					}
				}
			}
			if runInteractive == false {
				os.Exit(0)
			}
		}

		if runInteractive == true {
			// Add basic help
			js.AddHelp()
			// Add EPrint API help
			api.AddHelp(js)
			// build autocomplete list
			js.AddAutoComplete()
			// print welcome
			js.PrintDefaultWelcome()
			js.Repl()
			os.Exit(0)
		}
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
		data, err = api.GetPublishedRecords(0, publishedNewest, epgo.Descending)
	case publishedOldest > 0:
		data, err = api.GetPublishedRecords(0, publishedOldest, epgo.Ascending)
	case articlesNewest > 0:
		data, err = api.GetPublishedArticles(0, articlesNewest, epgo.Descending)
	case articlesOldest > 0:
		data, err = api.GetPublishedArticles(0, articlesOldest, epgo.Ascending)
	case restAPI == true:
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
		log.Fatalf("%s", err)
	}

	if prettyPrint == true {
		src, _ = json.MarshalIndent(data, "", "    ")
	} else {
		src, _ = json.Marshal(data)
	}
	fmt.Printf("%s", src)
}

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
	"log"
	"os"
	"strconv"

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

	exportEPrints   int
	feedSize        int
	publishedOldest int
	publishedNewest int
	articlesOldest  int
	articlesNewest  int
)

// addEPrintExtensionsAndHelp creates a *otto.Otto (JavaScript VM) with functions added to integrate
// epgo.
func addEPrintExtensionsAndHelp(api *epgo.EPrintsAPI, js *ostdlib.JavaScriptVM) *otto.Otto {
	vm := js.VM
	errorObject := func(obj *otto.Object, msg string) otto.Value {
		if obj == nil {
			obj, _ = vm.Object(`({})`)
		}
		log.Println(msg)
		obj.Set("status", "error")
		obj.Set("error", msg)
		return obj.Value()
	}

	// responseObject := func(data interface{}) otto.Value {
	// 	src, _ := json.Marshal(data)
	// 	obj, _ := vm.Object(fmt.Sprintf(`(%s)`, src))
	// 	return obj.Value()
	// }

	// EPrints REST API methods
	apiObj, _ := vm.Object(`api = {}`)
	js.SetHelp("api", "listEPrintsURI", []string{}, "returns and array of EPrint uris")
	apiObj.Set("listEPrintsURI", func(call otto.FunctionCall) otto.Value {
		uris, err := api.ListEPrintsURI()
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listEPrintsURI() failed %s, %s", call.CallerLocation(), err))
		}
		result, err := vm.ToValue(uris)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listEPrintsURI() failed %s, %s", call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "getEPrint", []string{"uri string"}, "given an EPrint uri (e.g. /eprint/2026) return the eprint as a JavaScript object")
	apiObj.Set("getEPrint", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("Missing uri for getEPrint() %s", call.CallerLocation()))
		}
		uri := call.Argument(0).String()
		record, err := api.GetEPrint(uri)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getEPrint(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(record)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getEPrint(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "exportEPrints", []string{"N int"}, "Export N items from eprints. Exports highest value id first. If N == -1 then export everything")
	apiObj.Set("exportEPrints", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("Missing export count for exportEPrints(n) %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		cnt, _ := strconv.Atoi(s)
		err := api.ExportEPrints(cnt)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("exportEprints(%d) failed %s, %s", cnt, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("exportEprints(%d) failed %s, %s", cnt, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "listURI", []string{}, "return a list of URI saved in the boltdb")
	apiObj.Set("listURI", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			return errorObject(nil, fmt.Sprintf("listURI(start, count) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		uris, err := api.ListURI(start, count)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listURI(%d, %d) failed %s, %s", start, count, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(uris)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listURI(%d, %d) failed %s, %s", start, count, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "get", []string{"uri string"}, "given a uri, return the results saved in the boltdb")
	apiObj.Set("get", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("get(uri) missing parameters %s", call.CallerLocation()))
		}
		uri := call.Argument(0).String()
		record, err := api.Get(uri)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("get(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(record)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("get(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "getPublishedRecords", []string{"start int", "count int", "direction int"}, "generate a list of published eprints starting at start for count by direction. Direction is 0 for ascending, 1 for descending")
	apiObj.Set("getPublishedRecords", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(start, count, direction) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		s = call.Argument(2).String()
		direction, _ := strconv.Atoi(s)
		records, err := api.GetPublishedRecords(start, count, direction)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "getPublishedArticles", []string{"start int", "count int", "direction int"}, "generate a list of published articles eprints starting at start for count by direction. Direction is 0 for ascending, 1 for descending")
	apiObj.Set("getPublishedArticles", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(start, count, direction) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		s = call.Argument(2).String()
		direction, _ := strconv.Atoi(s)
		records, err := api.GetPublishedArticles(start, count, direction)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "renderDocuments", []string{"docTitle string", "docDescription string", "basepath string", "records an array of eprints"},
		"render the eprint records list as an HTML file, HTML include, RSS and JSON documents")
	apiObj.Set("renderDocuments", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 4 {
			return errorObject(nil, fmt.Sprintf("renderDocuments(docTitle, docDescription, basepath, records) missing parameters %s", call.CallerLocation()))
		}
		docTitle := call.Argument(0).String()
		docDescription := call.Argument(1).String()
		basepath := call.Argument(2).String()
		records := []*epgo.Record{}
		err := ostdlib.ToStruct(call.Argument(3), records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) records error %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		if err := api.RenderDocuments(docTitle, docDescription, basepath, records); err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) failed %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) failed %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		return result
	})

	js.SetHelp("api", "buildSite", []string{"feedSize int"}, "build feeds recently-published and recent-articles. feedSize indicates the maximun number of items included")
	apiObj.Set("buildSite", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("buildSite(feed_size) missing parameter %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		feedSize, _ := strconv.Atoi(s)
		err := api.BuildSite(feedSize)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("buildSite(%d) failed %s, %s", feedSize, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("buildSite(%d) failed %s, %s", feedSize, call.CallerLocation(), err))
		}
		return result
	})
	return vm
}

func init() {
	publishedNewest = 0
	publishedOldest = 0
	feedSize = epgo.DefaultFeedSize

	flag.BoolVar(&showHelp, "h", false, "display help info")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
	flag.BoolVar(&showLicense, "l", false, "show license information")
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
 USAGE: epgo [OPTIONS] [EPRINT_URI]

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

	if runInteractive == true {
		vm := otto.New()
		js := ostdlib.New(vm)
		// Add basic help
		js.AddHelp()
		// Add extensions
		js.AddExtensions()
		// Add API specific extensions
		addEPrintExtensionsAndHelp(api, js)
		// build autocomplete list
		js.AddAutoComplete()
		// print welcome
		js.PrintDefaultWelcome()
		js.Repl()
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

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
// doi2eprintsxml.go is a Caltech Library centric command line utility
// to query CrossRef API and DataCite API for metadata and
// return the results as an EPrints XML file suitable for importing
// into EPrints.
//

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/crossrefapi"
	"github.com/caltechlibrary/dataciteapi"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/eprinttools/clsrules"
)

var (
	helpText = `---
title: "{app_name}"
author: "R. S. Doiel"
pubDate: 2023-01-11
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] DOI

# DESCRIPTION

{app_name} is a Caltech Library centric application that
takes one or more DOI, queries the CrossRef API
and if that fails the DataCite API and returns an
EPrints XML document suitable for import into
EPrints. The DOI can be in either their canonical
form or URL form (e.g. "10.1021/acsami.7b15651" or
"https://doi.org/10.1021/acsami.7b15651").

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-D
: attempt to download the digital object if object URL provided

-c
: only search CrossRef API for DOI records

-clsrules
: Apply current Caltech Library Specific Rules to EPrintXML output (default true)

-crossref
: only search CrossRef API for DOI records

-d
: only search DataCite API for DOI records

-datacite
: only search DataCite API for DOI records

-dot-initials
: Add period to initials in given name

-download
: attempt to download the digital object if object URL provided

-eprints-url string
: Sets the EPRints API URL

-i, -input
: (string) set input filename

-json
: output EPrint structure as JSON

-m 
: (string) set the mailto value for CrossRef API access (default "helpdesk@library.caltech.edu")

-mailto
: (string) set the mailto value for CrossRef API access (default "helpdesk@library.caltech.edu")

-normalize-publisher
: Use normalize publisher rule

-normalize-related-url
: Use normlize related url rule

-normlize-publication
: Use normalize publication rule

-o, -output
: (string) set output filename

-quiet
: set quiet output

-simple
: output EPrint structure as Simplified JSON

-trim-creators
: Use trim creators list rule

-trim-number
: Use trim number rule

-trim-series
: Use trim series rule

-trim-title
: Use trim title rule

-trim-volume
: Use trim volume rule

# EXAMPLES

Example generating an EPrintsXML for one DOI

~~~
	{app_name} "10.1021/acsami.7b15651" > article.xml
~~~

Example generating an EPrintsXML for two DOI

~~~
	{app_name} "10.1021/acsami.7b15651" "10.1093/mnras/stu2495" > articles.xml
~~~

Example processing a list of DOIs in a text file into
an XML document called "import-articles.xml".

~~~
	{app_name} -i doi-list.txt -o import-articles.xml
~~~

{app_name} {version}

`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	generateMarkdown bool
	generateManPage  bool
	inputFName       string
	outputFName      string
	quiet            bool

	// App specific options
	apiEPrintsURL                  string
	mailto                         string
	crossrefOnly                   bool
	dataciteOnly                   bool
	useCaltechLibrarySpecificRules bool
	asJSON                         bool
	asSimplified                   bool
	attemptDownload                bool
	trimTitleRule                  bool
	trimVolumeRule                 bool
	trimNumberRule                 bool
	pruneCreatorsRule              bool
	pruneSeriesRule                bool
	normalizeRelatedUrlRule        bool
	normalizePublisherRule         bool
	normalizePublicationRule       bool
	dotInitials                    bool
)

func fmtTxt(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display app version")
	flag.StringVar(&inputFName, "i", "", "set input filename")
	flag.StringVar(&inputFName, "input", "", "set input filename")
	flag.StringVar(&outputFName, "o", "", "set output filename")
	flag.StringVar(&outputFName, "output", "", "set output filename")
	flag.BoolVar(&quiet, "quiet", false, "set quiet output")

	// Application Options
	flag.StringVar(&apiEPrintsURL, "eprints-url", "", "Sets the EPRints API URL")
	flag.BoolVar(&crossrefOnly, "c", false, "only search CrossRef API for DOI records")
	flag.BoolVar(&crossrefOnly, "crossref", false, "only search CrossRef API for DOI records")
	flag.BoolVar(&dataciteOnly, "d", false, "only search DataCite API for DOI records")
	flag.BoolVar(&dataciteOnly, "datacite", false, "only search DataCite API for DOI records")
	flag.BoolVar(&useCaltechLibrarySpecificRules, "clsrules", true, "Apply current Caltech Library Specific Rules to EPrintXML output")
	flag.BoolVar(&trimTitleRule, "trim-title", false, "Use trim title rule")
	flag.BoolVar(&trimVolumeRule, "trim-volume", false, "Use trim volume rule")
	flag.BoolVar(&trimNumberRule, "trim-number", false, "Use trim number rule")
	flag.BoolVar(&pruneCreatorsRule, "trim-creators", false, "Use trim creators list rule")
	flag.BoolVar(&pruneSeriesRule, "trim-series", false, "Use trim series rule")
	flag.BoolVar(&normalizeRelatedUrlRule, "normalize-related-url", false, "Use normlize related url rule")
	flag.BoolVar(&normalizePublisherRule, "normalize-publisher", false, "Use normalize publisher rule")
	flag.BoolVar(&normalizePublicationRule, "normlize-publication", false, "Use normalize publication rule")
	flag.BoolVar(&dotInitials, "dot-initials", false, "Add period to initials in given name")
	flag.BoolVar(&asJSON, "json", false, "output EPrint structure as JSON")
	flag.BoolVar(&asSimplified, "simple", false, "output EPrint structure as Simplified JSON")
	flag.BoolVar(&attemptDownload, "D", false, "attempt to download the digital object if object URL provided")
	flag.BoolVar(&attemptDownload, "download", false, "attempt to download the digital object if object URL provided")

	//FIXME: Need to come up with a better way of setting this,
	// perhaps a config mode and save the setting in
	// $HOME/etc/${AppName}.json
	flag.StringVar(&mailto, "mailto", "helpdesk@library.caltech.edu", "set the mailto value for CrossRef API access")
	flag.StringVar(&mailto, "m", "helpdesk@library.caltech.edu", "set the mailto value for CrossRef API access")

	flag.Parse()
	args := flag.Args()

	// Setup I/O
	var err error

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if inputFName != "" && inputFName != "-" {
		in, err = os.Open(inputFName)
		if err != nil {
			fmt.Fprintln(eout, err)
			os.Exit(1)
		}
		defer in.Close()
	}

	if outputFName != "" && outputFName != "-" {
		out, err = os.Create(outputFName)
		if err != nil {
			fmt.Fprintln(eout, err)
			os.Exit(1)
		}
		defer out.Close()
	}

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, eprinttools.Version))
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintf(out, "%s\n", eprinttools.LicenseText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(out, "%s %s\n", eprinttools.Version)
		os.Exit(0)
	}


	if len(args) < 1 {
		fmt.Fprintln(eout, "Missing parameters see %s -help", appName)
		os.Exit(1)
	}

	if inputFName != "" {
		src, err := ioutil.ReadAll(in)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		for _, line := range strings.Split(fmt.Sprintf("%s", src), "\n") {
			arg := strings.TrimSpace(line)
			if len(arg) > 0 {
				args = append(args, arg)
			}
		}
	}

	// NOTE: OK we're ready to run our conversions
	eprintsList := new(eprinttools.EPrints)
	//FIXME: If crossrefapi returns 404 then we need to
	// query the dataciteapi
	apiCrossRef, err := crossrefapi.NewCrossRefClient(appName, mailto)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	apiDataCite, err := dataciteapi.NewDataCiteClient(appName, mailto)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}

	//NOTE: need to support processing one or more DOI
	for _, arg := range args {
		doi, objectURL := arg, ""
		if strings.Contains(arg, "|") {
			p := strings.SplitN(arg, "|", 2)
			doi, objectURL = p[0], p[1]
		}
		switch {
		case crossrefOnly:
			obj, err := apiCrossRef.Works(doi)
			if err != nil {
				fmt.Fprintf(eout, "ERROR (CrossRef API) %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiCrossRef.StatusCode == 200 {
				// NOTE: First we see if we can get a CrossRef record
				eprint, err := eprinttools.CrossRefWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(eout, "ERROR (CrossRef to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.Append(eprint)
				}
			} else {
				fmt.Fprintf(eout, "WARNING (CrossRef API): %q, %s\n", doi, apiCrossRef.Status)
			}
		case dataciteOnly:
			obj, err := apiDataCite.Works(doi)
			if err != nil {
				fmt.Fprintf(eout, "ERROR (DataCite API): %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiDataCite.StatusCode == 200 {
				eprint, err := eprinttools.DataCiteWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(eout, "ERROR (DataCite to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.Append(eprint)
				}
			} else {
				fmt.Fprintf(eout, "WARNING (DataCite API): %q, %s\n", doi, apiDataCite.Status)
			}
		default:
			// NOTE: just done for readability for flagging failed lookups
			isCrossRefDOI := false
			isDataCiteDOI := false

			obj, err := apiCrossRef.Works(doi)
			if err != nil {
				fmt.Fprintf(eout, "ERROR (CrossRef API): %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiCrossRef.StatusCode == 200 {
				isCrossRefDOI = true
				// NOTE: First we see if we can get a CrossRef record
				eprint, err := eprinttools.CrossRefWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(eout, "ERROR (CrossRef to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.Append(eprint)
				}
			}

			// NOTE: We try DataCite's API as a fallback when CrossRef fials...
			if isCrossRefDOI == false {
				obj, err := apiDataCite.Works(doi)
				if err != nil {
					fmt.Fprintf(eout, "ERROR (DataCite API): %q, %s\n", doi, err)
					os.Exit(1)
				}
				if apiDataCite.StatusCode == 200 {
					isDataCiteDOI = true
					eprint, err := eprinttools.DataCiteWorksToEPrint(obj)
					if err != nil {
						fmt.Fprintf(eout, "ERROR (DataCite to EPrintXML): skipping %q, %s\n", doi, err)
					} else {
						eprintsList.Append(eprint)
					}
				}
			}
			if attemptDownload && objectURL != "" {
				u, err := url.Parse(objectURL)
				if err != nil {
					fmt.Fprintf(eout, "Can't parse url %q, %s\n", objectURL, err)
					os.Exit(1)
				}
				response, err := http.Get(objectURL)
				if err != nil {
					fmt.Fprintf(eout, "Can't retrieve %q, %s\n", objectURL, err)
					os.Exit(1)
				}
				data, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Fprintf(eout, "Failed to read %q, %q\n", objectURL, err)
					os.Exit(1)
				}
				fName := path.Base(u.Path)
				if err := ioutil.WriteFile(fName, data, 0666); err != nil {
					fmt.Fprintf(eout, "Could not write %q from %q, %s\n", fName, objectURL, err)
					os.Exit(1)
				}
			}
			if isCrossRefDOI == false && isDataCiteDOI == false {
				fmt.Fprintf(eout, "WARNING: %s not found in CrossRef or DataCite API lookup\n", doi)
			}
		}
	}
	// NOTE: We can an option to apply all Caltech Library Special Rules
	// before marshaling our results or select individual rules.
	ruleSet := clsrules.ClearRuleSet()
	if useCaltechLibrarySpecificRules {
		ruleSet = clsrules.UseCLSRules()
	}
	if dotInitials {
		ruleSet["dot_initials"] = dotInitials
	}
	if trimTitleRule {
		ruleSet["trim_title"] = trimTitleRule
	}
	if trimVolumeRule {
		ruleSet["trim_volume"] = trimVolumeRule
	}
	if trimNumberRule {
		ruleSet["trim_number"] = trimNumberRule
	}
	if pruneCreatorsRule {
		ruleSet["prune_creators"] = pruneCreatorsRule
	}
	if pruneSeriesRule {
		ruleSet["prune_series"] = pruneSeriesRule
	}
	if normalizeRelatedUrlRule {
		ruleSet["normalize_related_url"] = normalizeRelatedUrlRule
	}
	if normalizePublisherRule {
		ruleSet["normalize_publisher"] = normalizePublisherRule
	}
	if normalizePublicationRule {
		ruleSet["normalize_publication"] = normalizePublicationRule
	}
	eprintsList, err = clsrules.Apply(eprintsList, ruleSet)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	if outputFName != `` {
		out, err = os.Create(outputFName)
		if err != nil {
			fmt.Fprintf(eout, `Can't write %q, %s`, outputFName, err)
			os.Exit(1)
		}
	}
	if asSimplified {
		fmt.Fprintln(out, "[")
		if eprintsList != nil && eprintsList.EPrint != nil {
			//for i, item := range eprintsList.EPrint {
			for i := 0; i < len(eprintsList.EPrint); i++ {
				item := eprintsList.EPrint[i]
				if i > 0 {
					fmt.Fprintf(out, ",\n")
				}
				rec, err := eprinttools.CrosswalkEPrintToRecord(item)
				if err != nil {
					fmt.Fprintf(eout, "%s\n", err)
					os.Exit(1)
				}
				fmt.Fprintf(out, "%s", rec.ToString())
			}
		}
		fmt.Fprintln(out, "\n]")
		os.Exit(0)
	}
	if asJSON {
		src, err := json.MarshalIndent(eprintsList, "", "   ")
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(out, "%s\n", src)
		os.Exit(0)
	}
	src, err := xml.MarshalIndent(eprintsList, "", "   ")
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(out, "%s\n", src)
}
